package util

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	mRand "math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	availableKeysHash = "rsa:available"
	inUseKeysHash     = "rsa:inuse"
	minKeyThreshold   = 5               // 密钥池最小阈值
	maxKeyThreshold   = 20              // 密钥池最大阈值
	checkInterval     = 5 * time.Second // 检查间隔
)

// KeyPair 返回公钥和过期时间
type KeyPair struct {
	PublicKey string `json:"public_key"` // Base64格式的公钥
}

// RSACrypto 管理RSA密钥池
type RSACrypto struct {
	Redis    *redis.Client
	stopChan chan struct{}
}

// NewRSACrypto 创建新的密钥池实例
func NewRSACrypto(redisClient *redis.Client) *RSACrypto {
	if redisClient == nil {
		panic("Redis client is nil")
	}
	pool := &RSACrypto{
		Redis:    redisClient,
		stopChan: make(chan struct{}),
	}
	go pool.startKeyPoolMonitor()
	return pool
}

// GenerateAndStoreKeys 批量生成并存储密钥对
func (r *RSACrypto) GenerateAndStoreKeys(count int) error {
	ctx := context.Background()

	for i := 0; i < count; i++ {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return err
		}

		publicKey := &privateKey.PublicKey

		privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
		privateKeyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		})

		publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
		publicKeyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		})

		publicKeyStr := base64.StdEncoding.EncodeToString(publicKeyPEM)
		privateKeyStr := base64.StdEncoding.EncodeToString(privateKeyPEM)

		// 使用HSET存储密钥对
		err = r.Redis.HSet(ctx, availableKeysHash, publicKeyStr, privateKeyStr).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO: 标记为已使用时添加过期时间，使用ZSET配合HSET
// GetRandomKeyPair 随机获取一个可用密钥对并标记为已使用
func (r *RSACrypto) GetRandomKeyPair() (KeyPair, error) {
	ctx := context.Background()
	const maxRetries = 3 // 改为常量声明

	for i := 0; i < maxRetries; i++ {
		var selectedKey string
		err := r.Redis.Watch(ctx, func(tx *redis.Tx) error {
			// 合并错误处理
			keys, err := tx.HKeys(ctx, availableKeysHash).Result()
			if err != nil || len(keys) == 0 {
				if err = r.generateKeysIfNeeded(tx, len(keys)); err != nil {
					return err
				}
				// 直接重新获取结果
				return tx.HKeys(ctx, availableKeysHash).Err()
			}

			// 使用更高效的随机选择方式
			selectedKey = keys[mRand.Intn(len(keys))] // 直接使用全局随机源
			// 补充获取私钥步骤
			privateKey, err := tx.HGet(ctx, availableKeysHash, selectedKey).Result()
			if err != nil {
				return err
			}

			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				// 新增：将私钥存入inUseKeysHash
				if err = pipe.HSet(ctx, inUseKeysHash, selectedKey, privateKey).Err(); err != nil {
					return err
				}
				// 保留原有的删除操作
				return pipe.HDel(ctx, availableKeysHash, selectedKey).Err()
			})
			return err
		}, availableKeysHash) // 明确watch的key

		switch {
		case err == nil:
			return KeyPair{
				PublicKey: selectedKey,
			}, nil
		case errors.Is(err, redis.TxFailedErr):
			time.Sleep(100 * time.Millisecond)
		default:
			return KeyPair{}, fmt.Errorf("failed to get key pair: %w", err)
		}
	}
	return KeyPair{}, errors.New("max retries reached")
}

// 新增辅助函数
func (r *RSACrypto) generateKeysIfNeeded(tx *redis.Tx, currentCount int) error {
	if currentCount < minKeyThreshold {
		if err := r.GenerateAndStoreKeys(minKeyThreshold - currentCount); err != nil {
			return fmt.Errorf("密钥生成失败: %w", err)
		}
		// 强制刷新事务中的key状态
		if _, err := tx.HKeys(context.Background(), availableKeysHash).Result(); err != nil {
			return fmt.Errorf("刷新密钥列表失败: %w", err)
		}
	}
	return nil
}

// DecryptWithKey 使用已标记的密钥对解密数据
func (r *RSACrypto) DecryptWithKey(publicKeyStr, encryptedData string) (string, error) {
	ctx := context.Background()

	// 从已使用哈希中获取私钥
	privateKeyStr, err := r.Redis.HGet(ctx, inUseKeysHash, publicKeyStr).Result()
	if err != nil {
		return "", errors.New("key not found in used keys")
	}

	// 删除已使用的密钥
	err = r.Redis.HDel(ctx, inUseKeysHash, publicKeyStr).Err()
	if err != nil {
		return "", fmt.Errorf("failed to delete used key: %v", err)
	}

	// 解码私钥
	privateKeyPEM, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return "", errors.New("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 解码加密数据
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	// 解密数据
	decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedBytes)
	if err != nil {
		return "", err
	}

	return string(decryptedBytes), nil
}

// startKeyPoolMonitor 启动密钥池监控
func (r *RSACrypto) startKeyPoolMonitor() {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.checkAndRefillKeys()
		case <-r.stopChan:
			return
		}
	}
}

// checkAndRefillKeys 检查并补充密钥
func (r *RSACrypto) checkAndRefillKeys() error {
	ctx := context.Background()

	// 获取当前可用密钥数量
	count, err := r.Redis.HLen(ctx, availableKeysHash).Result()
	if err != nil {
		return err
	}

	// 如果数量低于最小阈值，补充密钥
	if count < minKeyThreshold {
		generateCount := minKeyThreshold - int(count)
		return r.GenerateAndStoreKeys(generateCount)
	}

	// 如果数量超过最大阈值，不操作
	if count > maxKeyThreshold {
		return nil
	}

	return nil
}

// Stop 停止密钥池监控
func (r *RSACrypto) Stop() {
	close(r.stopChan)
}

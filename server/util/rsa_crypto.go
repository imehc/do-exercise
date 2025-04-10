package util

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// KeyPair 返回公钥和过期时间
type KeyPair struct {
	PublicKey  string  `json:"public_key"`  // Base64格式的公钥
	ExpireTime float64 `json:"expire_time"` // 过期时间（秒）
}

// RSACrypto 管理RSA加密、解密和密钥池功能
type RSACrypto struct {
	Redis         *redis.Client // Redis客户端，用于存储密钥对
	ExpireMinutes time.Duration // 密钥过期时间
	keyPool       chan *KeyPair // 缓存池，用于存储预生成的密钥对
	poolMutex     sync.Mutex    // 互斥锁，确保线程安全
	minPoolSize   int           // 缓存池最小密钥对数量
	maxRedisKeys  int           // Redis中允许的最大密钥对数量
}

// NewRSACrypto 创建一个新的RSACrypto实例，并初始化密钥池
func NewRSACrypto(redisClient *redis.Client) *RSACrypto {
	if redisClient == nil {
		panic("Redis client is nil")
	}
	rsaCrypto := &RSACrypto{
		Redis:         redisClient,
		ExpireMinutes: 5 * time.Minute,
		keyPool:       make(chan *KeyPair, 10), // 修改为 *KeyPair 类型
		minPoolSize:   5,
		maxRedisKeys:  5 << 2, // 默认最大Redis密钥对数量
	}
	go rsaCrypto.refillKeyPool()
	return rsaCrypto
}

// GenerateKeyPair 生成RSA密钥对并存储到Redis中
func (r *RSACrypto) GenerateKeyPair() (*KeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
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

	ctx := context.Background()
	expireTime := r.ExpireMinutes
	err = r.Redis.Set(ctx, "publKey:"+publicKeyStr, privateKeyStr, expireTime).Err()
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		PublicKey:  publicKeyStr,
		ExpireTime: expireTime.Seconds(),
	}, nil
}

// refillKeyPool 定期检查并补充密钥池中的密钥对
func (r *RSACrypto) refillKeyPool() {
	for {
		r.poolMutex.Lock()

		// 检查Redis中已有的密钥对数量
		ctx := context.Background()
		currentKeys, err := r.Redis.Keys(ctx, "publKey:*").Result()
		if err != nil {
			r.poolMutex.Unlock()
			time.Sleep(1 * time.Second)
			continue
		}

		// 如果Redis中的密钥对数量超过最大限制，则跳过生成
		if len(currentKeys) >= r.maxRedisKeys {
			r.poolMutex.Unlock()
			time.Sleep(1 * time.Second)
			continue
		}

		// 补充密钥池
		if len(r.keyPool) < r.minPoolSize {
			for i := 0; i < r.minPoolSize-len(r.keyPool); i++ {
				keyPair, err := r.GenerateKeyPair()
				if err != nil {
					continue
				}
				r.keyPool <- keyPair // 无需类型转换
			}
		}
		r.poolMutex.Unlock()
		time.Sleep(1 * time.Second)
	}
}

// GetKeyFromPool 从密钥池中获取一个密钥对，如果池为空则等待
func (r *RSACrypto) GetKeyFromPool() (*KeyPair, error) {
	select {
	case key := <-r.keyPool:
		return key, nil
	default:
		key := <-r.keyPool
		return key, nil
	}
}

// VerifyAndDecrypt 验证并使用Redis中的私钥解密数据
func (r *RSACrypto) VerifyAndDecrypt(publicKeyStr, encryptedData string) (string, error) {
	ctx := context.Background()

	defer func() {
		_ = r.Redis.Del(ctx, "publKey:"+publicKeyStr).Err()
	}()

	privateKeyStr, err := r.Redis.Get(ctx, "publKey:"+publicKeyStr).Result()
	if err != nil {
		return "", errors.New("未找到私钥")
	}

	privateKeyPEM, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return "", errors.New("解码PEM块失败")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}

	decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedBytes)
	if err != nil {
		return "", err
	}

	return string(decryptedBytes), nil
}

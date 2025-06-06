package util

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRSACrypto(t *testing.T) {
	// 初始化Redis客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	defer func() {
		ctx := context.Background()
		// 清理测试数据
		redisClient.Del(ctx, availableKeysHash)
		// 清理所有以 publicKey: 开头的key
		iter := redisClient.Scan(ctx, 0, publicKeyPrefix+"*", 0).Iterator()
		for iter.Next(ctx) {
			redisClient.Del(ctx, iter.Val())
		}
		redisClient.Close()
	}()

	rsaPool := NewRSACrypto(redisClient)

	t.Run("测试批量密钥生成", func(t *testing.T) {
		err := rsaPool.GenerateAndStoreKeys(5)
		assert.NoError(t, err)

		ctx := context.Background()
		count, err := redisClient.HLen(ctx, availableKeysHash).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})

	t.Run("测试随机获取密钥对", func(t *testing.T) {
		err := rsaPool.GenerateAndStoreKeys(3)
		assert.NoError(t, err)

		keyPair, err := rsaPool.GetRandomKeyPair()
		assert.NoError(t, err)
		assert.NotEmpty(t, keyPair.PublicKey)

		ctx := context.Background()
		// 验证密钥已存储为独立key
		key := publicKeyPrefix + keyPair.PublicKey
		exists, err := redisClient.Exists(ctx, key).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(1), exists)

		// 验证密钥已从available哈希移除
		isExist, err := redisClient.HExists(ctx, availableKeysHash, keyPair.PublicKey).Result()
		assert.NoError(t, err)
		assert.False(t, isExist)

		// 验证过期时间设置
		ttl, err := redisClient.TTL(ctx, key).Result()
		assert.NoError(t, err)
		assert.True(t, ttl > 0 && ttl <= keyExpiration)
	})

	t.Run("测试解密功能", func(t *testing.T) {
		err := rsaPool.GenerateAndStoreKeys(1)
		assert.NoError(t, err)

		keyPair, err := rsaPool.GetRandomKeyPair()
		assert.NoError(t, err)

		originalData := "test data 123"
		encryptedData, err := mockEncryptData(keyPair.PublicKey, originalData)
		assert.NoError(t, err)

		decryptedData, err := rsaPool.DecryptWithKey(keyPair.PublicKey, encryptedData, true)
		assert.NoError(t, err)
		assert.Equal(t, originalData, decryptedData)

		// 验证密钥在解密后被删除
		ctx := context.Background()
		key := publicKeyPrefix + keyPair.PublicKey
		exists, err := redisClient.Exists(ctx, key).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), exists)
	})

	t.Run("测试密钥过期", func(t *testing.T) {
		err := rsaPool.GenerateAndStoreKeys(1)
		assert.NoError(t, err)

		keyPair, err := rsaPool.GetRandomKeyPair()
		assert.NoError(t, err)

		ctx := context.Background()
		key := publicKeyPrefix + keyPair.PublicKey

		// 强制使密钥过期 (使用1毫秒而不是1纳秒)
		err = redisClient.PExpire(ctx, key, time.Millisecond).Err()
		assert.NoError(t, err)

		// 等待密钥过期 (稍微延长等待时间以确保过期)
		time.Sleep(5 * time.Millisecond)

		// 验证过期的密钥无法使用
		_, err = rsaPool.DecryptWithKey(keyPair.PublicKey, "some_data", true)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "operationTimeout")
	})
}

// mockEncryptData 保持不变
func mockEncryptData(publicKeyStr string, data string) (string, error) {
	// 解码base64的公钥
	publicKeyPEM, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		return "", err
	}

	// 解析PEM格式的公钥
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return "", errors.New("failed to decode PEM block")
	}

	// 解析公钥
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 使用公钥加密数据
	encryptedBytes, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(data))
	if err != nil {
		return "", err
	}

	// 将加密后的数据转换为base64
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

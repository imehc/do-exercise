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
		// 清理测试数据
		ctx := context.Background()
		redisClient.Del(ctx, availableKeysHash)
		redisClient.Del(ctx, inUseKeysHash)
		redisClient.Close()
	}()

	rsaPool := NewRSACrypto(redisClient)

	t.Run("测试批量密钥生成", func(t *testing.T) {
		err := rsaPool.GenerateAndStoreKeys(5)
		assert.NoError(t, err)

		// 验证密钥数量
		ctx := context.Background()
		count, err := redisClient.HLen(ctx, availableKeysHash).Result()
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})

	t.Run("测试随机获取密钥对", func(t *testing.T) {
		// 先生成一些密钥
		err := rsaPool.GenerateAndStoreKeys(3)
		assert.NoError(t, err)

		// 获取密钥对
		keyPair, err := rsaPool.GetRandomKeyPair()
		assert.NoError(t, err)
		assert.NotEmpty(t, keyPair.PublicKey)

		// 验证密钥已移动到inuse哈希
		ctx := context.Background()
		_, err = redisClient.HGet(ctx, inUseKeysHash, keyPair.PublicKey).Result()
		assert.NoError(t, err)

		// 验证密钥已从available哈希移除
		_, err = redisClient.HGet(ctx, availableKeysHash, keyPair.PublicKey).Result()
		assert.Error(t, err)
	})

	t.Run("测试解密功能", func(t *testing.T) {
		// 生成密钥对并获取
		err := rsaPool.GenerateAndStoreKeys(1)
		assert.NoError(t, err)

		keyPair, err := rsaPool.GetRandomKeyPair()
		assert.NoError(t, err)

		// 模拟加密数据
		originalData := "test data 123"
		encryptedData, err := mockEncryptData(keyPair.PublicKey, originalData)
		assert.NoError(t, err)

		// 解密数据
		decryptedData, err := rsaPool.DecryptWithKey(keyPair.PublicKey, encryptedData)
		assert.NoError(t, err)
		assert.Equal(t, originalData, decryptedData)
	})
}

// mockEncryptData 模拟加密过程
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

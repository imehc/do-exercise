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

	// 测试完成后关闭Redis连接
	defer func() {
		// 清理所有测试相关的key
		iter := redisClient.Scan(context.Background(), 0, "publKey:*", 0).Iterator()
		for iter.Next(context.Background()) {
			redisClient.Del(context.Background(), iter.Val())
		}
		redisClient.Close()
	}()

	// 初始化RSACrypto
	rsaCrypto := NewRSACrypto(redisClient)
	rsaCrypto.ExpireMinutes = 10 * time.Second // 直接设置10秒

	t.Run("测试密钥生成", func(t *testing.T) {
		result, err := rsaCrypto.GenerateKeyPair()
		assert.NoError(t, err)
		assert.NotEmpty(t, result["public_key"])
		assert.Equal(t, float64(10), result["expire_time"])

		// 测试完成后清理密钥
		defer func() {
			publicKey := result["public_key"].(string)
			redisClient.Del(context.Background(), "publKey:"+publicKey)
		}()
	})

	t.Run("测试加密解密流程", func(t *testing.T) {
		// 1. 生成密钥对
		result, err := rsaCrypto.GenerateKeyPair()
		assert.NoError(t, err)
		publicKey := result["public_key"].(string)

		// 2. 原始数据
		originalData := "test123456!@#$%^"

		// 3. 模拟加密过程
		encryptedData, err := mockJSEncryptEncryption(publicKey, originalData)
		assert.NoError(t, err)

		// 4. 解密数据
		decryptedData, err := rsaCrypto.VerifyAndDecrypt(publicKey, encryptedData)
		assert.NoError(t, err)
		assert.Equal(t, originalData, decryptedData)
	})

	t.Run("测试密钥过期", func(t *testing.T) {
		// 1. 生成密钥对
		result, err := rsaCrypto.GenerateKeyPair()
		assert.NoError(t, err)
		publicKey := result["public_key"].(string)

		// 2. 等待密钥过期
		time.Sleep(time.Second * 12) // 等待12秒

		// 3. 尝试使用过期的密钥
		_, err = rsaCrypto.VerifyAndDecrypt(publicKey, "some_encrypted_data")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "private key not found")
	})
}

// mockJSEncryptEncryption 模拟前端JSEncrypt的加密过程
func mockJSEncryptEncryption(publicKeyStr string, data string) (string, error) {
	// 1. 解码base64的公钥
	publicKeyPEM, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		return "", err
	}

	// 2. 解析PEM格式的公钥
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return "", errors.New("failed to decode PEM block")
	}

	// 3. 解析公钥
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	// 4. 使用公钥加密数据
	encryptedBytes, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(data))
	if err != nil {
		return "", err
	}

	// 5. 将加密后的数据转换为base64
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

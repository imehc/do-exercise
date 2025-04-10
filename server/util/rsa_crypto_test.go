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
		iter := redisClient.Scan(context.Background(), 0, "publKey:*", 0).Iterator()
		for iter.Next(context.Background()) {
			redisClient.Del(context.Background(), iter.Val())
		}
		redisClient.Close()
	}()

	rsaCrypto := NewRSACrypto(redisClient)
	rsaCrypto.ExpireMinutes = 20 * time.Second

	t.Run("测试批量密钥生成", func(t *testing.T) {
		numKeys := 5
		var publicKeys []string

		for i := 0; i < numKeys; i++ {
			result, err := rsaCrypto.GenerateKeyPair()
			assert.NoError(t, err)
			assert.NotEmpty(t, result.PublicKey)
			publicKeys = append(publicKeys, result.PublicKey)

			defer func(pubKey string) {
				redisClient.Del(context.Background(), "publKey:"+pubKey)
			}(result.PublicKey)
		}

		// 验证所有公钥都存储在Redis中
		for _, publicKey := range publicKeys {
			privateKeyStr, err := redisClient.Get(context.Background(), "publKey:"+publicKey).Result()
			assert.NoError(t, err)
			assert.NotEmpty(t, privateKeyStr)
		}
	})

	t.Run("测试并发加密解密", func(t *testing.T) {
		numWorkers := 10
		data := "test123456!@#$%^"
		results := make(chan error, numWorkers)

		// 并发执行加密和解密操作
		for i := 0; i < numWorkers; i++ {
			go func() {
				keyPair, err := rsaCrypto.GenerateKeyPair()
				if err != nil {
					results <- err
					return
				}

				encryptedData, err := mockJSEncryptEncryption(keyPair.PublicKey, data)
				if err != nil {
					results <- err
					return
				}

				decryptedData, err := rsaCrypto.VerifyAndDecrypt(keyPair.PublicKey, encryptedData)
				if err != nil {
					results <- err
					return
				}

				if decryptedData != data {
					results <- errors.New("解密数据不匹配")
					return
				}

				results <- nil
			}()
		}

		// 等待所有goroutine完成并检查结果
		for i := 0; i < numWorkers; i++ {
			err := <-results
			assert.NoError(t, err)
		}
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

package util

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RSACrypto struct {
	Redis         *redis.Client
	ExpireMinutes time.Duration
}

// NewRSACrypto 创建新的RSA加密实例
func NewRSACrypto(redisClient *redis.Client) *RSACrypto {
	if redisClient == nil {
		panic("Redis client is nil")
	}
	return &RSACrypto{
		Redis:         redisClient,
		ExpireMinutes: 5 * time.Minute, // 默认5分钟过期
	}
}

// GenerateKeyPair 生成RSA密钥对并存储到Redis
func (r *RSACrypto) GenerateKeyPair() (map[string]interface{}, error) {
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

	return map[string]interface{}{
		"public_key":  publicKeyStr,
		"expire_time": expireTime.Seconds(),
	}, nil
}

// VerifyAndDecrypt 验证并解密数据
func (r *RSACrypto) VerifyAndDecrypt(publicKeyStr string, encryptedData string) (string, error) {
	ctx := context.Background()

	// 确保在函数返回前删除密钥
	defer func() {
		_ = r.Redis.Del(ctx, "publKey:"+publicKeyStr).Err()
	}()

	privateKeyStr, err := r.Redis.Get(ctx, "publKey:"+publicKeyStr).Result()
	if err != nil {
		return "", errors.New("private key not found")
	}

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

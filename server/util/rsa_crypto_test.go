package util

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// 本用例需要真实 Redis。默认连本机，并允许用环境变量指向别处（与集成测试同一套变量名）。
// 用独立的 DB 15：本用例会 DEL availableKeysHash 与全部 publicKey:*，落在 DB 0 会清掉
// 开发环境正在用的密钥池。
const (
	rsaTestRedisDefaultAddr = "localhost:6379"
	rsaTestRedisDB          = 15
)

func rsaTestRedisAddr() string {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		return rsaTestRedisDefaultAddr
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}
	return host + ":" + port
}

func TestRSACrypto(t *testing.T) {
	// 初始化Redis客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     rsaTestRedisAddr(),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       rsaTestRedisDB,
	})

	// Redis 不可达时跳过而不是失败：否则干净环境里 `go test ./...` 默认是红的，
	// 真实回归会被这一条淹没。要强制跑（例如 CI）就设 RSA_TEST_REQUIRE_REDIS=1。
	pingCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := redisClient.Ping(pingCtx).Err(); err != nil {
		redisClient.Close()
		if os.Getenv("RSA_TEST_REQUIRE_REDIS") != "" {
			t.Fatalf("RSA_TEST_REQUIRE_REDIS 已设置但 Redis(%s) 不可用: %v", rsaTestRedisAddr(), err)
		}
		t.Skipf("跳过：Redis(%s) 不可用（%v）。起一个 Redis 或设置 REDIS_HOST/REDIS_PORT 后重跑", rsaTestRedisAddr(), err)
	}

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

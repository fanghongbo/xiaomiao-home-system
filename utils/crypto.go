package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// 默认迭代次数和密钥长度
	iterations = 10000
	keyLength  = 32 // AES-256
	saltLength = 16
)

// EncryptString 使用 AES-256-GCM 加密字符串
// key: 加密密钥（建议使用配置中的密钥）
// plaintext: 要加密的明文
// 返回: base64编码的加密字符串和错误
func EncryptString(key, plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	// 生成随机盐值
	salt := make([]byte, saltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// 使用 PBKDF2 派生密钥
	derivedKey := pbkdf2.Key([]byte(key), salt, iterations, keyLength, sha256.New)

	// 创建 AES 密码块
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// 创建 GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// 生成随机 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// 将 salt 和 ciphertext 组合并 base64 编码
	// 格式: base64(salt + ciphertext)
	combined := append(salt, ciphertext...)
	encoded := base64.StdEncoding.EncodeToString(combined)

	return encoded, nil
}

// DecryptString 解密字符串
// key: 解密密钥（必须与加密时使用的密钥相同）
// ciphertext: base64编码的加密字符串
// 返回: 解密后的明文和错误
func DecryptString(key, ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	// Base64 解码
	combined, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// 检查长度
	if len(combined) < saltLength {
		return "", fmt.Errorf("invalid ciphertext length")
	}

	// 提取 salt 和加密数据
	salt := combined[:saltLength]
	encryptedData := combined[saltLength:]

	// 使用 PBKDF2 派生密钥
	derivedKey := pbkdf2.Key([]byte(key), salt, iterations, keyLength, sha256.New)

	// 创建 AES 密码块
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// 创建 GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// 检查长度
	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// 提取 nonce 和密文
	nonce, ciphertextBytes := encryptedData[:nonceSize], encryptedData[nonceSize:]

	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

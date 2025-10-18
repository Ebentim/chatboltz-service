package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

type EncryptionKey struct {
	key []byte
}

func NewEncryptionKey(key []byte) *EncryptionKey {
	return &EncryptionKey{key: key}
}
func (c *EncryptionKey) GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *EncryptionKey) EncryptString(data []byte) (string, error) {
	blocked, err := aes.NewCipher(c.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(blocked)
	if err != nil {
		return "", err
	}
	nounce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nounce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nounce, nounce, data, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil

}

func (c *EncryptionKey) DecryptString(data string) ([]byte, error) {
	decodeData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	blocked, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(blocked)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := decodeData[:nonceSize], decodeData[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

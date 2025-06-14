package postgresql

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
)

type Cursor struct {
	secretKey []byte
}

func NewCursor(secretKey []byte) *Cursor {

	return &Cursor{secretKey: secretKey}
}

func (c *Cursor) Encrypt(id int) (string, error) {
	block, err := aes.NewCipher(c.secretKey)
	if err != nil {
		return "", err
	}

	plaintext := make([]byte, 4)
	binary.BigEndian.PutUint32(plaintext, uint32(id))

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func (c *Cursor) Decrypt(encodedCursor string) (int, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(encodedCursor)
	if err != nil {
		return 0, err
	}

	block, err := aes.NewCipher(c.secretKey)
	if err != nil {
		return 0, err
	}

	if len(ciphertext) < aes.BlockSize {
		return 0, fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return int(binary.BigEndian.Uint32(ciphertext)), nil
}

// Chase Weaver

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

// createHash(string)
// Creates an hashed string from a key
func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// encryptAES([]byte, string)
// Encrypts a 64bit string
func encryptAES(key []byte, text string) string {
	pt := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ct := make([]byte, aes.BlockSize+len(pt))
	iv := ct[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ct[aes.BlockSize:], pt)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ct)
}

// decryptAES([]byte, string)
// Decrpyts an encrypted 64bit string
func decryptAES(key []byte, cryptoText string) string {
	ct, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ct) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ct[:aes.BlockSize]
	ct = ct[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ct, ct)

	return fmt.Sprintf("%s", ct)
}

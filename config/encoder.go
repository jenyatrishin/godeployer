package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os/user"
)

type Encoder struct {
	Phrase []byte
	SecretKey []byte
}

func (e *Encoder) encode(sentence string) string {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	e.Phrase = []byte(sentence)
	e.SecretKey = []byte(currentUser.Name)

	c, err := aes.NewCipher(e.SecretKey)
	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		fmt.Println(err)
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}

	return string(gcm.Seal(nonce, nonce, e.Phrase, nil))
}

func (e *Encoder) decode(sentence string) string {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	e.SecretKey = []byte(currentUser.Name)

	c, err := aes.NewCipher(e.SecretKey)
	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}

	nonceSize := gcm.NonceSize()
	if len(sentence) < nonceSize {
		fmt.Println(err)
	}

	nonce, sentence := sentence[:nonceSize], sentence[nonceSize:]
	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(sentence), nil)
	if err != nil {
		fmt.Println(err)
	}

	return string(plaintext)
}
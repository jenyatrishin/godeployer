package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"gopkg.in/gookit/color.v1"
	"io"
	"io/ioutil"
	"os"
	mRand "math/rand"
	"../tools"
	"strings"
	"time"
)

const (
	FOLDERNAME   string = ".dep2go"
	KEY_FILENAME string = "key"
)

type Encoder struct {
	Phrase []byte
	SecretKey []byte
}

func (e *Encoder) Encode(sentence string) []byte {
	e.Phrase = []byte(sentence)
	e.SecretKey = []byte(e.ReadKeyFile())

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

	return gcm.Seal(nonce, nonce, e.Phrase, nil)
}

func (e *Encoder) Decode(sentence []byte) string {
	e.SecretKey = e.ReadKeyFile()

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

	nonce, sentenceE := sentence[:nonceSize], sentence[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, sentenceE, nil)
	if err != nil {
		fmt.Println("1")
		fmt.Println(err)
	}

	return string(plaintext)
}

func (e *Encoder) generateKey() string {
	mRand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 16
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[mRand.Intn(len(chars))])
	}
	return b.String()
}

func (e *Encoder) WriteKeyFile() {
	key := []byte(e.generateKey())

	writeErr := ioutil.WriteFile(keyFileName(), key, os.ModePerm)
	if writeErr != nil {
		tools.WriteLog("Init error: " + writeErr.Error())
		color.Red.Println(writeErr.Error())
		os.Exit(1)
	}
}

func (e *Encoder) ReadKeyFile() []byte {
	key, err := ioutil.ReadFile(keyFileName())
	if err != nil {
		tools.WriteLog("Can't read key file - " + err.Error())
		color.Red.Println("Can't read key file - " + err.Error())
		os.Exit(1)
	}

	return key
}

func keyFileName() string {
	return FOLDERNAME + "/" + KEY_FILENAME
}
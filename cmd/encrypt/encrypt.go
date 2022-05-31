package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: encrypt path")
		return
	}

	key := os.Getenv("ENCRYPTION_KEY")
	if key == "" {
		fmt.Println("Missing ENCRYPTION_KEY environment variable")

		return
	}

	path := os.Args[1]

	var content string
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	} else {
		iv, content = decryptFile(key, path)
	}

	ext := filepath.Ext(path)
	tmpfile, err := os.CreateTemp("", fmt.Sprintf("scrt.*%s", ext))
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	tmpfile.WriteString(content)
	tmpfile.Close()

	editor := getEnv("EDITOR", "nvim")
	editor_arguments := strings.Fields(os.Getenv("EDITOR_ARGUMENTS"))
	if strings.Contains(editor, "vi") && len(editor_arguments) == 0 {
		editor_arguments = []string{"-c", "set ft=json"}
	}

	editor_arguments = append(editor_arguments, tmpfile.Name())
	cmd := exec.Command(editor, editor_arguments[:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	contentB, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		panic(err)
	}

	encryptFile(key, string(contentB), path, iv)
}

func decryptFile(keyS, path string) ([]byte, string) {
	key, err := hex.DecodeString(keyS)
	if err != nil {
		panic(err)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(string(raw))
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)
	return iv, string(ciphertext)
}

func encryptFile(keyS, text, path string, iv []byte) {
	key, err := hex.DecodeString(keyS)
	if err != nil {
		panic(err)
	}

	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	copy(ciphertext[:aes.BlockSize], iv)

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	dir := filepath.Dir(path)

	os.MkdirAll(dir, 0755)

	b64 := base64.StdEncoding.EncodeToString(ciphertext)

	err = os.WriteFile(path, []byte(b64), 0600)
	if err != nil {
		panic(err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

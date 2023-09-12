package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type PIIEncryptor struct {
	Key []byte
}

func (c *PIIEncryptor) Encrypt(text []byte) (string, error) {
	ciphertext := make([]byte, len(text))
	for i := range text {
		ciphertext[i] = text[i] ^ c.Key[i%len(c.Key)]
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (c *PIIEncryptor) Decrypt(cipher string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(cipher)
	if err != nil {
		return cipher, err
	}

	plaintext := make([]byte, len(ciphertext))
	for i := range ciphertext {
		plaintext[i] = ciphertext[i] ^ c.Key[i%len(c.Key)]
	}

	return string(plaintext), nil
}

func main() {
	// Initialize the PIIEncryptor with a random string key.
	key := "Hello Amber 123!" // Replace this with your random string key.
	encdec := &PIIEncryptor{Key: []byte(key)}

	http.HandleFunc("/encrypt", func(w http.ResponseWriter, r *http.Request) {
		req := map[string]string{}

		data := json.NewDecoder(r.Body)
		data.Decode(&req)

		res, err := encdec.Encrypt([]byte(req["plaintext"]))
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte(res))
	})

	http.HandleFunc("/decrypt", func(w http.ResponseWriter, r *http.Request) {
		req := map[string]string{}

		data := json.NewDecoder(r.Body)
		data.Decode(&req)

		res, err := encdec.Decrypt(req["ciphertext"])
		if err != nil {
			log.Fatal(err)
		}

		w.Write([]byte(res))
	})

	fmt.Println("starting web server at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}

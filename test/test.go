package main

import (
	"fmt"
	"telegram_test_bot/BL/Encrypt"
)

func main() {
	Secret := "akrosbccho&1*~#^2^#s0^=)"
	StringToEncrypt := "6c857c8a19578f61d1d5da29dff56c28e7071c4017cf3a935b93899cad40c47e1601fd7b9c4a5a6a09ec8"
	// To encrypt the StringToEncrypt
	encText, err := Encrypt.Encrypt(StringToEncrypt, Secret)
	if err != nil {
		fmt.Println("error encrypting your classified text: ", err)
	}
	fmt.Println("Enc: ", encText)
	// To decrypt the original StringToEncrypt
	decText, err := Encrypt.Decrypt(encText, Secret)
	if err != nil {
		fmt.Println("error decrypting your encrypted text: ", err)
	}
	fmt.Println("Dec: ", decText)
}

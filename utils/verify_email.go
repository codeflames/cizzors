package utils

import (
	"fmt"
	"net/mail"
)

func EmailVerifier(email string) error {

	_, err := mail.ParseAddress(email)
	if err != nil {
		fmt.Println("Invalid email:", email)
		return err
	}

	fmt.Println("Valid email:", email)
	return nil
}

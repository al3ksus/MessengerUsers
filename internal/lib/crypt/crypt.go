package crypt

import "golang.org/x/crypto/bcrypt"

type Crypter struct{}

func (c *Crypter) GenerateFromPassword(pass []byte, cost int) ([]byte, error) {
	passHash, err := bcrypt.GenerateFromPassword(pass, cost)
	if err != nil {
		return nil, err
	}

	return passHash, nil
}

func (c *Crypter) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	if err := bcrypt.CompareHashAndPassword(hashedPassword, password); err != nil {
		return err
	}

	return nil
}

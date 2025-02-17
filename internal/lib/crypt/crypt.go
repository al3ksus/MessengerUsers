package crypt

import "golang.org/x/crypto/bcrypt"

//Crypter - структура реализует методы хэширования пароля и сравнения пароля с хэшем.
type Crypter struct{}

// GenerateFromPassword возвращает хэш указанного пароля с заданной стоимостью.
// Использует пакет bcrypt.
func (c *Crypter) GenerateFromPassword(pass []byte, cost int) ([]byte, error) {
	passHash, err := bcrypt.GenerateFromPassword(pass, cost)
	if err != nil {
		return nil, err
	}

	return passHash, nil
}

// CompareHashAndPassword сравнивает захэшированный пароль с исходным.
// Использует пакет bcrypt.
func (c *Crypter) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	if err := bcrypt.CompareHashAndPassword(hashedPassword, password); err != nil {
		return err
	}

	return nil
}

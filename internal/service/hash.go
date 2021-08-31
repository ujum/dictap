package service

import "golang.org/x/crypto/bcrypt"

type PassHashService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type BCryptPassHashService struct {
}

func (B *BCryptPassHashService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (B *BCryptPassHashService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

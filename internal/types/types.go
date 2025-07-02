package types

import (
	"errors"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

type HashString struct {
	str  string
	hash bool
}

func NewHashString(str string) *HashString {
	return &HashString{
		str:  str,
		hash: true,
	}
}

func (hs *HashString) Hash() (string, error) {
	if hs == nil {
		return "", nil
	}

	if hs.str == "" {
		return "", nil
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(hs.str), 14)
	return string(bytes), err
}

func (hs *HashString) Equal(pwd string) bool {
	if hs == nil {
		return false
	}

	if hs.hash {
		err := bcrypt.CompareHashAndPassword([]byte(hs.str), []byte(pwd))
		return err == nil
	}
	return hs.str == pwd
}

func (hs *HashString) String() string {
	if hs == nil {
		return ""
	}
	return hs.str
}

type Email string

func NewEmail(s string) (Email, error) {
	if s == "" {
		return "", errors.New("email is required")
	}
	addr, err := mail.ParseAddress(s)
	if err != nil {
		return "", errors.New(err.Error())
	}
	return Email(addr.Address), nil
}

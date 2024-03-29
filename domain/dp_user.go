package domain

import (
	"errors"
	"net/url"
	"strings"
)

// Account
type Account interface {
	Account() string
}

func NewAccount(v string) (Account, error) {
	if v == "" || strings.ToLower(v) == "root" || !reName.MatchString(v) {
		return nil, errors.New("invalid user name")
	}

	return dpAccount(v), nil
}

type dpAccount string

func (r dpAccount) Account() string {
	return string(r)
}

// Password
type Password interface {
	Password() string
}

func NewPassword(s string) (Password, error) {
	if n := len(s); n < 8 || n > 20 {
		return nil, errors.New("invalid password")
	}

	part := make([]bool, 4)

	for _, c := range s {
		if c >= 'a' && c <= 'z' {
			part[0] = true
		} else if c >= 'A' && c <= 'Z' {
			part[1] = true
		} else if c >= '0' && c <= '9' {
			part[2] = true
		} else {
			part[3] = true
		}
	}

	i := 0
	for _, b := range part {
		if b {
			i++
		}
	}
	if i < 3 {
		return nil, errors.New(
			"the password must includes three of lowercase, uppercase, digital and special character",
		)
	}

	return dpPassword(s), nil
}

type dpPassword string

func (r dpPassword) Password() string {
	return string(r)
}

// Nickname
type Nickname interface {
	Nickname() string
}

func NewNickname(v string) (Nickname, error) {
	if len(v) > config.User.MaxNicknameLength {
		return nil, errors.New("invalid nickname")
	}

	return dpNickname(v), nil
}

type dpNickname string

func (r dpNickname) Nickname() string {
	return string(r)
}

// Bio
type Bio interface {
	Bio() string
}

func NewBio(v string) (Bio, error) {
	if len(v) > config.User.MaxBioLength {
		return nil, errors.New("invalid bio")
	}

	return dpBio(v), nil
}

type dpBio string

func (r dpBio) Bio() string {
	return string(r)
}

// Email
type Email interface {
	Email() string
}

func NewEmail(v string) (Email, error) {
	if v == "" || !reEmail.MatchString(v) {
		return nil, errors.New("invalid email")
	}

	return dpEmail(v), nil
}

type dpEmail string

func (r dpEmail) Email() string {
	return string(r)
}

// AvatarId
type AvatarId interface {
	AvatarId() string
}

func NewAvatarId(v string) (AvatarId, error) {
	if v == "" {
		return nil, errors.New("invalid avatar")
	}

	if _, err := url.Parse(v); err != nil {
		return nil, errors.New("invalid avatar")
	}

	return dpAvatarId(v), nil
}

type dpAvatarId string

func (r dpAvatarId) AvatarId() string {
	return string(r)
}

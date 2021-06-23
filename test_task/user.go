package main

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"os"
	//    "strconv"
	//    "reflect"
)

type User struct {
	Email        string
	PasswordHash string
	Token        string
}

func PasswordHash(password string) string {
	salt := "Yeeh_zMVk3"
	hasher := sha512.New512_256()
	hasher.Write([]byte(password + salt))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func FindByEmail(usersdb *os.File, email string) (User, bool) {
	csvLines, err := csv.NewReader(usersdb).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	for _, line := range csvLines {
		if line[0] == email {
			u := User{
				Email:        line[0],
				PasswordHash: line[1],
				Token:        line[2],
			}
			return u, true
		}
	}
	return User{}, false
}

func FindByToken(usersdb *os.File, token string) (User, bool) {
	csvLines, err := csv.NewReader(usersdb).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	for _, line := range csvLines {
		if line[2] == token {
			u := User{
				Email:        line[0],
				PasswordHash: line[1],
				Token:        line[2],
			}
			return u, true
		}
	}
	return User{}, false
}

func UserCreate(email, pass string) error {
	if email == "" {
		return fmt.Errorf("Incorrect email")
	}

	if pass == "" {
		return fmt.Errorf("Password can't be empty")
	}
	passHash := PasswordHash(pass)

	usersDbFilename := "/home/admin/go/genesis_school/test_task/usersdb.csv"

	usersdb, err := os.OpenFile(
		usersDbFilename,
		os.O_APPEND|os.O_CREATE|os.O_RDWR,
		0644)
	if err != nil {
		return err
	}
	defer usersdb.Close()

	if _, has := FindByEmail(usersdb, email); has {
		return fmt.Errorf("Email already used")
	}

	u := User{
		Email:        email,
		PasswordHash: passHash,
		Token:        passHash,
	}
	_, err = usersdb.Write([]byte(u.Email + "," + u.PasswordHash + "," + u.Token + "\n"))
	if err != nil {
		return err
	}

	return nil
}

func UserLogin(email, pass string) (string, error) {
	if email == "" {
		return "", fmt.Errorf("Incorrect email")
	}

	passHash := PasswordHash(pass)

	usersDbFilename := "/home/admin/go/genesis_school/test_task/usersdb.csv"

	usersdb, _ := os.OpenFile(
		usersDbFilename,
		os.O_RDONLY,
		0644)
	defer usersdb.Close()

	u, has := FindByEmail(usersdb, email)
	if !has || u.PasswordHash != passHash {
		return "", fmt.Errorf("Incorrect login")
	}

	return u.Token, nil
}

func IsAvaiableToken(token string) bool {
	if token == "" {
		return false
	}

	usersDbFilename := "/home/admin/go/genesis_school/test_task/usersdb.csv"

	usersdb, _ := os.OpenFile(
		usersDbFilename,
		os.O_RDONLY,
		0644)
	defer usersdb.Close()

	_, has := FindByToken(usersdb, token)
	return has
}

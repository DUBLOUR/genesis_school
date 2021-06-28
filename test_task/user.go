package main

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"time"
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

func RandomString(length int) string {
	var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyz")

	rand.Seed(time.Now().UnixNano())
	s := make([]rune, length)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func FindByEmailOrToken(email, token string) (User, bool) {
	usersDb, err := os.OpenFile(dbFile, os.O_RDONLY, 0644)
	if err != nil {
		return User{}, false
	}

	defer usersDb.Close()

	csvLines, err := csv.NewReader(usersDb).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	for _, line := range csvLines {
		u := User{
			Email:        line[0],
			PasswordHash: line[1],
			Token:        line[2],
		}
		if u.Email == email || u.Token == token {
			return u, true
		}
	}
	return User{}, false
}

func FindByEmail(email string) (User, bool) {
	return FindByEmailOrToken(email, "")
}

func FindByToken(token string) (User, bool) {
	return FindByEmailOrToken("", token)
}

func AppendUser(u User) error {
	usersDb, err := os.OpenFile(dbFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer usersDb.Close()

	_, err = usersDb.Write([]byte(u.Email + "," + u.PasswordHash + "," + u.Token + "\n"))
	if err != nil {
		return err
	}
	return nil
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsEmailValid(email string) bool {
	return len(email) >= 3 && len(email) <= 255 && emailRegex.MatchString(email)
}

func UserRegister(email, pass string) error {
	if email == "" {
		return fmt.Errorf("Missed email")
	}

	if !IsEmailValid(email) {
		return fmt.Errorf("Incorrect email")
	}

	if pass == "" {
		return fmt.Errorf("Password can't be empty")
	}

	if _, has := FindByEmail(email); has {
		return fmt.Errorf("Email already used")
	}

	var u User
	u.Email = email
	u.PasswordHash = PasswordHash(pass)
	u.Token = RandomString(12)

	if err := AppendUser(u); err != nil {
		return err
	}

	return nil
}

func UserLogin(email, pass string) (string, error) {
	if email == "" {
		return "", fmt.Errorf("Incorrect email")
	}

	u, has := FindByEmail(email)
	if !has || u.PasswordHash != PasswordHash(pass) {
		return "", fmt.Errorf("Incorrect login")
	}

	return u.Token, nil
}

func IsAvaiableToken(token string) bool {
	if token == "" {
		return false
	}

	_, has := FindByToken(token)
	return has
}

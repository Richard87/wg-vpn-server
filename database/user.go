package database

import (
	crand "crypto/rand"
	"github.com/Richard87/wg-vpn-server/config"
	"github.com/alexedwards/argon2id"
	"gorm.io/gorm"
	"log"
	"math/big"
	"runtime"
	"strings"
)

type User struct {
	gorm.Model
	Username string
	Hash     string
	Role     string
}

func InitUsers() {
	usersCreated := false

	users := []User{}
	Connection.Find(&users)

	for _, user := range users {
		usersCreated = true

		for i, u := range config.Config.Users {
			parts := strings.Split(u, ":")
			if len(parts) < 2 {
				log.Fatalf("When creating a new user, it must be in format username:password! (%s used)", u)
			}

			if parts[0] != user.Username {
				continue
			}

			hash, err := argon2id.CreateHash(parts[1], &argon2id.Params{
				Memory:      65536,
				Iterations:  19,
				Parallelism: uint8(runtime.NumCPU()),
				SaltLength:  16,
				KeyLength:   16,
			})
			if err != nil {
				log.Fatalf("could not update hash for %s: %s", user.Username, err)
			}

			user.Hash = hash
			if len(parts) == 3 {
				user.Role = parts[2]
			}

			Connection.Save(user)
			removeIndexFromUsersList(i)
		}
	}

	for i, u := range config.Config.Users {
		parts := strings.Split(u, ":")
		if len(parts) < 2 {
			log.Fatalf("error in users database! Failed user: %d (%v)", i, parts)
		}

		hash, err := argon2id.CreateHash(parts[1], &argon2id.Params{
			Memory:      65536,
			Iterations:  19,
			Parallelism: uint8(runtime.NumCPU()),
			SaltLength:  16,
			KeyLength:   16,
		})
		if err != nil {
			log.Fatalf("could not update hash for %s: %s", u, err)
		}

		var user = &User{
			Username: parts[0],
			Hash:     hash,
			Role:     "admin",
		}
		if len(parts) == 3 {
			user.Role = parts[2]
		}

		Connection.Create(user)
		usersCreated = true
	}

	if !usersCreated {
		password, err := generatePassword(10)
		if err != nil {
			log.Fatalf("could not generate admin password: %s", err)
		}

		log.Println("Creating admin user with password: " + password)
		hash, err := argon2id.CreateHash(password, &argon2id.Params{
			Memory:      65536,
			Iterations:  19,
			Parallelism: uint8(runtime.NumCPU()),
			SaltLength:  16,
			KeyLength:   16,
		})
		if err != nil {
			log.Fatalf("could not generate admin password: %s", err)
		}

		newUser := User{
			Username: "admin",
			Hash:     hash,
			Role:     "admin",
		}

		Connection.Create(newUser)
	}
}

func removeIndexFromUsersList(i int) {
	config.Config.Users[i] = config.Config.Users[len(config.Config.Users)-1]
	config.Config.Users = config.Config.Users[:len(config.Config.Users)-1]
}

func generatePassword(length int) (string, error) {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅabcdefghijklmnopqrstuvwxyz0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		n, err := crand.Int(crand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		b.WriteRune(chars[int(n.Int64())])
	}
	s := b.String()
	return s, nil
}

package api

import (
	"fmt"
	"github.com/Richard87/wg-vpn-server/config"
	"github.com/Richard87/wg-vpn-server/database"
	"github.com/alexedwards/argon2id"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"
	"time"
)

func Authenticate(c *fiber.Ctx) error {
	var login = Login{}
	if err := c.BodyParser(login); err != nil {
		c.Status(http.StatusBadRequest)
		return nil
	}

	if login.Password == "" {
		c.Status(http.StatusUnauthorized)
		return nil
	}

	var user = database.User{}
	database.Connection.Find(&user, "username = ?", login.Username)
	if user.Username == "" {
		c.Status(http.StatusUnauthorized)
		return nil
	}

	match, err := argon2id.ComparePasswordAndHash(login.Password, user.Hash)
	if err != nil {
		return fmt.Errorf("could not check password hash!: %s", err)
	}

	if !match {
		c.Status(http.StatusUnauthorized)
		return nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": login.Username,
		"exp":      time.Now().Add(time.Minute * 115).UnixNano(),
	})

	tokenString, err := token.SignedString(config.Config.HttpsJwtSigningKey)
	if err != nil {
		return fmt.Errorf("could not sign jwt: %s", err)
	}

	parts := strings.Split(tokenString, ".")

	c.Cookie(&fiber.Cookie{
		Name:     "auth",
		Value:    parts[2],
		MaxAge:   3600 * 2,
		Expires:  time.Now().Add(time.Hour * 2),
		Secure:   true,
		HTTPOnly: true,
		SameSite: "lax",
	})

	return c.JSON(LoginResponse{Token: tokenString})
}

func NewAuthenticationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := c.Get("Authorization")
		signature := c.Cookies("auth")
		authPars := strings.Split(authorization, " ")
		jwtParts := strings.Split(authPars[1], ".")
		if signature == "" || len(jwtParts) != 2 {
			c.Status(http.StatusUnauthorized)
			return nil
		}

		tokenString := fmt.Sprintf("%s.%s", jwtParts[1], signature)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if token.Header["alg"] != "HS256" {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return config.Config.HttpsJwtSigningKey, nil
		})
		if err != nil || !token.Valid {
			c.Status(http.StatusForbidden)
			return nil
		}

		c.Locals("jwt", token)
		return c.Next()
	}
}

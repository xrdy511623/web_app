package models

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	ID          uint
	AuthorityId uint
	jwt.StandardClaims
}

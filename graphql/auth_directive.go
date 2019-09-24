package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/dgrijalva/jwt-go"
)

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://easy-isa.auth0.com/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}

var Auth_Directive = func(ctx context.Context, obj interface{}, next graphql.Resolver, role Role) (interface{}, error) {
	Itoken := ctx.Value("token")
	splitToken := strings.Split(Itoken.(string), "Bearer")
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(splitToken[1], claims, func(token *jwt.Token) (interface{}, error) {
		cert, nil := getPemCert(token)
		return []byte(cert), nil
	})
	if err == nil {
		fmt.Println("valid token", token, claims)
	} else {
		fmt.Println("invalid token", err)
	}
	return next(ctx)
}

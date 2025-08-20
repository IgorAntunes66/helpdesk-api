package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var chaveJWT = []byte(os.Getenv("SEGREDOJWT"))

// ClaimCustom continua o mesmo, pois não depende de nenhum modelo específico.
type ClaimCustom struct {
	UserID int64  `json:"userID"`
	Nome   string `json:"nome"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GerarToken agora recebe os dados do usuário diretamente.
// Isso quebra a dependência que tínhamos do pacote 'model' do 'users-service'.
func GerarToken(userID int64, nome, email string) (string, error) {
	claims := ClaimCustom{
		UserID: userID,
		Nome:   nome,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
			Issuer:    "help-desk-api",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(chaveJWT)
}

// ValidarToken pode ser melhorado para retornar as 'claims' em caso de sucesso.
func ValidarToken(tokenString string) (*ClaimCustom, error) {
	claims := &ClaimCustom{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return chaveJWT, nil
	})

	if err != nil {
		if err == jwt.ErrTokenExpired {
			return nil, fmt.Errorf("token expirado")
		}
		return nil, fmt.Errorf("token inválido: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	return claims, nil
}

func ExtairToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("autorização ausente no cabeçalho da requisição")
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
		return "", errors.New("header de autorização mal formatado")
	}

	tokenString := headerParts[1]

	return tokenString, nil
}

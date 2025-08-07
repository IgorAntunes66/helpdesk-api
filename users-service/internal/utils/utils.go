package utils

import (
	"errors"
	"fmt"
	"helpdesk/users-service/internal/model"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var chaveJWT = []byte(os.Getenv("SEGREDOJWT"))

type ClaimCustom struct {
	UserID int64  `json:"userID"`
	Nome   string `json:"nome"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GerarToken(user model.User) (string, error) {
	// Passo 1: Criando as claims
	claims := ClaimCustom{
		user.ID,
		user.Nome,
		user.Email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
			Issuer:    "help-desk-api",
		},
	}

	// Passo 2: Criando o token com as claims e o método de assinatura
	// O método de assinatura é parte do cabeçalho (Header). HS256 é um dos mais comuns
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Passo 4: Assinando o token com o nosso segredo
	// Isso cria a terceira parte do JWT (a assinatura)
	tokenAssinado, err := token.SignedString(chaveJWT)
	if err != nil {
		return "", err
	}

	return tokenAssinado, nil
}

func ValidarToken(tokenString string) error {
	// Vamos criar uma instancia vazia das nossas claims para popularmos
	claims := &ClaimCustom{}

	//Passo 1: Fazendo o parse e a validação do token
	// O ParseWithClaims vai:
	//1. Decodificar o token sem verificar a assinatura.
	//2.Chamar nossa função de callback para fornecer a chave de validação.
	//3. Validar a assinatura do token usando a chave que fornecemos.
	//4. Validar as claims padrão (Como a data de expiração).
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		//DICA DE SEGURANÇA: VERIFIQUE O MÉTODO DE ASSINATURA!
		// Isso previne um ataque onde um token malicioso usa "alg: none".
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		//Retorna a nossa chave secreta para a validação
		return chaveJWT, nil
	})

	//Passo 2: Tratando os erros de validação
	if err != nil {
		// O pacote jwt retorna erros especificos que podemos verificar.
		// Isso é otimo para dar feedback claro ao cliente (ex: Token Expirado).
		if err == jwt.ErrTokenExpired {
			return errors.New("token expirado")
		} else {
			return errors.New("erro ao validar token")
		}
	}

	//Passo 3: Verificando se o token é valido e extraindo as informações
	if token.Valid {
		return nil
	}

	return err
}

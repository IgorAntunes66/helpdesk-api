package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Definimos uma chave para usar no contexto e evitar colisões.
type contextKey string

const UserIDKey contextKey = "userID"

// AuthMiddleware é a "muralha magica".
// Ele recebe um 'handler' (a proxima sala do castelo) e retorna um novo 'handler' (o portao com o guarda).
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. O ponto de inspeção: Extrair o Selo (Token) do cabeçalho
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Cabeçalho de autorização não fornecido", http.StatusUnauthorized)
			return
		}

		// O cabeçalho geralmente vem como "Bearer <token>". Precisamos separar o token.
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
			http.Error(w, "Formato do cabeçalho de autorização é invalido", http.StatusUnauthorized)
			return
		}
		tokenString := headerParts[1]

		//2. A validação do Selo: Verificar a autenticidade do token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SEGREDOJWT")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Token invalido ou expirado", http.StatusUnauthorized)
			return
		}

		//3. O sussuro do guardiao: Extrair a indentidade e injetar no contexto
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Não foi possivel ler as reindivicações do token", http.StatusUnauthorized)
			return
		}

		//Extraimos o ID do usuario que colocamos no token durante o login.
		// O JWT armazena números como float64, então precisamos converter.
		userID := int64(claims["userID"].(float64))

		// Criamos um novo contexto que carrega o ID do usuario.
		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		// Criamos uma nova requisição com este novo contexto.
		reqWithContext := r.WithContext(ctx)

		// 4. Permissao concedida: chamar a proxima sala
		next.ServeHTTP(w, reqWithContext)
	})
}

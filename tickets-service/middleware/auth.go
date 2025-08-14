package middleware

import (
	"context"
	"helpdesk/tickets-service/auth"
	"net/http"
	"strings"
)

// Definimos uma chave para usar no contexto e evitar colisões.
type contextKey string

const UserIDKey contextKey = "userID"

// AuthMiddleware é a "muralha magica".
// Ele recebe um 'handler' (a proxima sala do castelo) e retorna um novo 'handler' (o portao com o guarda).
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ... (lógica para extrair o token do header)
		headerParts := strings.Split(r.Header.Get("Authorization"), " ")
		if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
			http.Error(w, "Formato do cabeçalho de autorização é invalido", http.StatusUnauthorized)
			return
		}
		tokenString := headerParts[1]

		// CHAMADA ATUALIZADA: Agora usamos nossa função de validação centralizada!
		claims, err := auth.ValidarToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Injetamos o ID do usuário (vindo das claims) no contexto da requisição.
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)

		// Passa a requisição com o novo contexto para o próximo handler.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

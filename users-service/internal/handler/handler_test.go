package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"helpdesk/users-service/internal/middleware"
	"helpdesk/users-service/internal/model"
	"helpdesk/users-service/internal/repository"
	"helpdesk/users-service/internal/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserHandler(t *testing.T) {
	//Arrange (Preparar)
	mockRepo := new(repository.MockUserRepository)

	//Criamos um usuario de exemplo que esperamos enviar na requisição
	userInput := model.User{
		Nome:     "John Doe",
		Senha:    "password123",
		TipoUser: "admin",
		Email:    "teste@gmail.com",
		Telefone: "123456789",
		CpfCnpj:  "12345678901",
	}

	//Configuramos o mock. Dizemos a ele:
	//"Eu espero que o método 'CreateUser' seja chamado com o 'userInput'.
	// Quando isso acontecer, você deve retornar o ID '1' e nenhum erro."
	mockRepo.On("CreateUser", userInput).Return(int64(1), nil)

	//Criamos nosso ApiServer usando o REPOSITORIO FALSO.
	apiServer := NewApiServer(mockRepo)

	//Convertemos nosso usuario de input para JSON, pois é assim que uma requisição HTTP funciona.
	body, _ := json.Marshal(userInput)
	req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	//Act (Ação)
	//Executamos o handler.
	apiServer.CreateUserHandler(rr, req)

	//Assert (Verificar)
	// 1. Verificamos se o status HTTP é 201 Created.
	assert.Equal(t, http.StatusCreated, rr.Code)

	// 2. Verificamos se o corpo da resposta contém o usuário com o ID que o mock retornou.
	var userResponse model.User
	err := json.NewDecoder(rr.Body).Decode(&userResponse)
	assert.NoError(t, err)                     // Não deve haver erro ao decodificar a resposta.
	assert.Equal(t, int64(1), userResponse.ID) //O ID deve ser 1.
	assert.Equal(t, userInput.Nome, userResponse.Nome)

	mockRepo.AssertExpectations(t)
}

func TestCreateUserHandler_RepositoryError(t *testing.T) {
	// Arrange (Preparar)
	mockRepo := new(repository.MockUserRepository)

	userInput := model.User{
		Nome:     "John Doe",
		Senha:    "password123",
		TipoUser: "admin",
		Email:    "teste@gmail.com",
		Telefone: "123456789",
		CpfCnpj:  "12345678901",
	}

	// A MUDANÇA CRUCIAL: A profecia da Falha
	//Agora, instruimos nosso dublê de uma forma diferente
	//"Eu espero que 'CreateUser' seja chamado com 'userInput'.
	//Quando isso acontecer, você deve retornar um ID zero E um NOVO ERRO."
	mockRepo.On("CreateUser", userInput).Return(int64(0), errors.New("erro de banco de dados"))
	apiServer := NewApiServer(mockRepo)

	body, _ := json.Marshal(userInput)
	req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	// Act (Agir)
	apiServer.CreateUserHandler(rr, req)

	// Assert (Verificar)
	//1. Verificamos se o handler agiu corretamente diante do erro.
	// Ele não deve retornar 201, mas sim um erro de servidor.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	//2. Garantimos que o dublê foi chamado exatamente como planejamos.
	mockRepo.AssertExpectations(t)
}

func TestListUserHandler(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)

	mockRepo.On("FindAllUsers").Return([]model.User{}, nil)
	apiServer := NewApiServer(mockRepo)

	req := httptest.NewRequest("GET", "/users/", nil)
	rr := httptest.NewRecorder()

	apiServer.ListUsersHandler(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	mockRepo.AssertExpectations(t)
}

func TestListUserHandler_RepositoryError(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)

	mockRepo.On("FindAllUsers").Return([]model.User{}, errors.New("erro de banco de dados"))
	apiServer := NewApiServer(mockRepo)

	req := httptest.NewRequest("GET", "/users", nil)
	rr := httptest.NewRecorder()

	apiServer.ListUsersHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	mockRepo.AssertExpectations(t)

}

func TestGetUserHandler(t *testing.T) {
	//Arrange (Preparar)
	mockRepo := new(repository.MockUserRepository)

	userInput := model.User{
		ID:       1,
		Nome:     "John Doe",
		Senha:    "password123",
		TipoUser: "admin",
		Email:    "teste@gmail.com",
		Telefone: "123456789",
		CpfCnpj:  "12345678901",
	}

	mockRepo.On("FindUserByID", int64(1)).Return(userInput, nil)
	apiServer := NewApiServer(mockRepo)

	req := httptest.NewRequest("GET", "/users/1", nil)
	rr := httptest.NewRecorder()

	//Injetando o contexto da rota
	//Criamos um contexto de rota do Chi e dizemos a ele que o parâmetro "id" tem o valor "1".
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "1")

	//Anexamos este contexto magico a nossa requisição.
	//Agora, quando o handler chamar 'chi.URLParam(r, "id")', ele encontrará o valor!
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	apiServer.GetUserHandler(rr, req)

	var userOutput model.User
	_ = json.NewDecoder(rr.Body).Decode(&userOutput)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, userInput.Nome, userOutput.Nome)
	assert.Equal(t, userInput.ID, userOutput.ID)

	mockRepo.AssertExpectations(t)
}

func TestGetUserHandler_RepositoryError(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)

	mockRepo.On("FindUserByID", int64(1)).Return(model.User{}, pgx.ErrNoRows)
	apiServer := NewApiServer(mockRepo)

	req := httptest.NewRequest("GET", "/users/1", nil)
	rr := httptest.NewRecorder()

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	apiServer.GetUserHandler(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestLoginUserHandler(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)

	senha := "senha123"

	mockUser := model.User{
		ID:    1,
		Email: "igorgantunes@hotmail.com",
		Senha: "senha123",
		Nome:  "Igor",
	}

	loginCredentials := model.LoginRequest{
		Email: mockUser.Email,
		Senha: senha,
	}

	mockRepo.On("FindUserByEmail", loginCredentials).Return(mockUser, nil)

	apiServer := NewApiServer(mockRepo)

	body, _ := json.Marshal(loginCredentials)
	req := httptest.NewRequest("POST", "/users/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	apiServer.LoginUserHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseBody map[string]string
	err := json.NewDecoder(rr.Body).Decode(&responseBody)
	assert.NoError(t, err, "O corpo da resposta deveria ser um JSON com a chave 'token'")

	tokenString, exists := responseBody["token"]
	assert.True(t, exists, "A resposta deveria conter um token")
	assert.NotEmpty(t, tokenString, "O token não pode estar vazio")

	claims := &utils.ClaimCustom{}

	// 2. Use ParseWithClaims para decodificar o token diretamente na sua struct.
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verificação de segurança crucial
		assert.IsType(t, &jwt.SigningMethodHMAC{}, token.Method, "Método de assinatura inesperado!")
		return []byte(os.Getenv("SEGREDOJWT")), nil
	})

	// 3. Verifique os resultados
	assert.NoError(t, err, "O token retornado deve ser válido")
	assert.True(t, token.Valid, "O token deve ser valido")

	// 4. Agora você pode acessar os campos da sua struct de forma segura e tipada!
	assert.Equal(t, mockUser.ID, claims.UserID, "O ID do usuario no token esta incorreto")
	assert.Equal(t, mockUser.Nome, claims.Nome, "O Nome do usuario no token esta incorreto")
	assert.Equal(t, mockUser.Email, claims.Email, "O Email do usuario no token esta incorreto")

	// Garante que a expectativa do mock foi atendida
	mockRepo.AssertExpectations(t)
}

func TestGetMeHandler_Success(t *testing.T) {
	mockRepo := new(repository.MockUserRepository)
	apiServer := NewApiServer(mockRepo)

	mockUser := model.User{
		ID:    1,
		Email: "igorgantunes@hotmail.com",
		Nome:  "Igor",
	}

	mockRepo.On("FindUserByID", mockUser.ID).Return(mockUser, nil)

	token, _ := utils.GerarToken(mockUser)

	req := httptest.NewRequest("GET", "/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/users/me", apiServer.GetMeHandler)
	})

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var userDB model.User
	_ = json.NewDecoder(rr.Body).Decode(&userDB)

	assert.Equal(t, mockUser.ID, userDB.ID)
	assert.Equal(t, mockUser.Nome, userDB.Nome)

	mockRepo.AssertExpectations(t)
}


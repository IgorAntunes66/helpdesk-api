package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"helpdesk/users-service/internal/model"
	"helpdesk/users-service/internal/repository"
	"net/http"
	"net/http/httptest"
	"testing"

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

	mockRepo.On("FindUserById", int64(1)).Return(userInput, nil)
	apiServer := NewApiServer(mockRepo)

	body, _ := json.Marshal(userInput)
	req := httptest.NewRequest("GET", "/users/1", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	apiServer.GetUserHandler(rr, req)

	var userOutput model.User
	_ = json.NewDecoder(rr.Body).Decode(&userOutput)

	assert.Equal(t, http.StatusFound, rr.Code)
	assert.Equal(t, userInput.Nome, userOutput.Nome)

	mockRepo.AssertExpectations(t)
}

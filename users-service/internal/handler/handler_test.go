package handler

import (
	"bytes"
	"encoding/json"
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

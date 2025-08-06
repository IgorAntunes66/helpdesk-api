package handler

import (
	"helpdesk/users-service/internal/model"
	"helpdesk/users-service/internal/repository"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestCreateUserHandler(t *testing.T) {
	t.Run("Deve criar um usuario com sucesso", func(t *testing.T) {
		mockUser := &model.User{ID: 1, Nome: "Joao Silva"}
		mockRepo := new(repository.MockUserRepository)

		mockRepo.On("CreateUser", mock.Anything, int64(1)).Return(mockUser, nil)

		handler := 
	})
}

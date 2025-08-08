package repository

import (
	"helpdesk/users-service/internal/model"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user model.User) (int64, error) {
	args := m.Called(user)
	return args.Get(0).(int64), args.Error(1)

}

func (m *MockUserRepository) FindAllUsers() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) FindUserByID(id int64) (model.User, error) {
	args := m.Called(id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) FindUserByEmail(email string) (model.User, error) {
	args := m.Called(email)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(id int64, user model.User) error {
	args := m.Called(id)
	return args.Error(1)
}

func (m *MockUserRepository) DeleteUser(id int64) error {
	args := m.Called(id)
	return args.Error(1)
}

func (m *MockUserRepository) GerarHashSenha(senha string) ([]byte, error) {
	args := m.Called(senha)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockUserRepository) VerificarSenha(hashSalvo string, senhaLogin string) error {
	args := m.Called(hashSalvo, senhaLogin)
	return args.Error(1)
}

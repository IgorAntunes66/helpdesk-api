package model

type User struct {
	ID       int64  `json:"id"`
	Nome     string `json:"nome"`
	Senha    string `json:"senha"`
	TipoUser string `json:"tipoUser"`
	Email    string `json:"email"`
	Telefone string `json:"telefone"`
	CpfCnpj  string `json:"cpf_cnpj"`
}

type UserRepository interface {
	CreateUser(user User) (int64, error)
	FindAllUsers() ([]User, error)
	FindUserByID(id int64) (User, error)
	UpdateUser(id int64, user User) error
	DeleteUser(id int64) error
}

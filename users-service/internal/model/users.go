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

package model

type User struct {
	ID       int    `json:"id"`
	Nome     string `json:"nome"`
	Senha    string `json:"senha"`
	Funcao   string `json:"funcao"`
	Telefone string `json:"telefone"`
	CpfCnpj  string `json:"cpf_cnpj"`
}

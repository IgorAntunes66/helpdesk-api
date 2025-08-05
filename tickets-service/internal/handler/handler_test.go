// O teste deve pertencer ao mesmo pacote que o código que está vigiando.
package handler

import (
	"net/http"          // Precisamos das constantes HTTP, como 'http.StatusOK'.
	"net/http/httptest" // A caixa de ferramentas mágica do Go para simular requisições e respostas HTTP.
	"testing"           // O pacote fundamental para a criação de qualquer teste.
)

// A função de teste deve começar com 'Test' e receber um ponteiro para 'testing.T'.
// 't' é o seu porta-voz para com o Oráculo; ele relata sucessos e falhas.
func TestHealthCheckHandler(t *testing.T) {
	// [O QUE FAZER]:
	// 1. Forjar uma Requisição Falsa:
	// Criamos uma requisição HTTP que simula um usuário acessando o endpoint.
	// Não há rede envolvida; é pura simulação.
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		// Se a própria criação da requisição falhar, é um erro fatal no teste.
		t.Fatalf("Não foi possível criar a requisição: %v", err)
	}

	// 2. Criar um "Gravador de Resposta":
	// 'httptest.NewRecorder()' é um tipo especial de 'ResponseWriter'
	// que, em vez de enviar uma resposta pela internet, "grava" o que
	// o handler tentou enviar (código de status, corpo, etc.).
	rr := httptest.NewRecorder()

	// 3. Invocar o Guardião (Handler):
	// Chamamos nosso 'HealthCheckHandler' diretamente, como uma função normal.
	// Passamos a ele nossa requisição falsa e nosso gravador.
	HealthCheckHandler(rr, req)

	// 4. Interrogar o Gravador:
	// Agora, verificamos o que o handler "gravou" em nosso 'rr'.
	// Esperamos que o código de status seja 'http.StatusOK' (200).
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("O handler retornou um status inesperado: esperado %v, recebido %v",
			http.StatusOK, status)
	}

	// [POR QUE FAZER]:
	// Este processo de 'Arrange, Act, Assert' (Preparar, Agir, Verificar) é a base
	// de todos os testes unitários. Ele nos permite isolar e verificar o comportamento
	// de uma única peça de nosso código (um 'handler') com total confiança,
	// garantindo que ela faça exatamente o que esperamos, sem depender do resto do sistema.
}

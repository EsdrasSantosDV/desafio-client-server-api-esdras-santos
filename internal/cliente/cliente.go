package cliente

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"client-server-api-esdras/internal/cotacao"
)

const requestTimeout = 300 * time.Millisecond

func Executar(endpoint, arquivo string) error {
	bid, err := consultarServidor(endpoint)
	if err != nil {
		return err
	}

	conteudo := fmt.Sprintf("Dólar: %s", bid)
	if err := os.WriteFile(arquivo, []byte(conteudo), 0644); err != nil {
		return err
	}

	log.Printf("cotacao salva em %s", arquivo)
	return nil
}

func consultarServidor(endpoint string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return "", fmt.Errorf("timeout de %s excedido: %w", requestTimeout, err)
		}
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("servidor retornou status %d", res.StatusCode)
	}

	var resposta cotacao.RespostaCliente
	if err := json.NewDecoder(res.Body).Decode(&resposta); err != nil {
		return "", err
	}

	return resposta.Bid, nil
}

package servidor

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"client-server-api-esdras/internal/cotacao"

	_ "modernc.org/sqlite"
)

const (
	apiURL          = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	apiTimeout      = 200 * time.Millisecond
	databaseTimeout = 10 * time.Millisecond
)

type Servidor struct {
	endereco string
	db       *sql.DB
}

func Novo(endereco, arquivoBanco string) (*Servidor, error) {
	db, err := sql.Open("sqlite", arquivoBanco)
	if err != nil {
		return nil, err
	}

	if err := criarTabela(db); err != nil {
		db.Close()
		return nil, err
	}

	return &Servidor{
		endereco: endereco,
		db:       db,
	}, nil
}

func (s *Servidor) Iniciar() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", s.buscarCotacao)

	log.Printf("servidor iniciado em http://localhost%s/cotacao", s.endereco)
	return http.ListenAndServe(s.endereco, mux)
}

func (s *Servidor) Fechar() error {
	return s.db.Close()
}

func criarTabela(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS cotacoes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT,
			codein TEXT,
			name TEXT,
			high TEXT,
			low TEXT,
			var_bid TEXT,
			pct_change TEXT,
			bid TEXT,
			ask TEXT,
			timestamp TEXT,
			create_date TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := db.Exec(query)
	return err
}

func (s *Servidor) buscarCotacao(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "metodo nao permitido", http.StatusMethodNotAllowed)
		return
	}

	dados, err := consultarAPI()
	if err != nil {
		log.Printf("erro ao buscar cotacao na API externa: %v", err)
		http.Error(w, "erro ao buscar cotacao", http.StatusGatewayTimeout)
		return
	}

	if err := s.salvar(dados); err != nil {
		log.Printf("erro ao salvar cotacao no banco: %v", err)
		http.Error(w, "erro ao salvar cotacao", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cotacao.RespostaCliente{Bid: dados.Bid}); err != nil {
		log.Printf("erro ao responder cliente: %v", err)
	}
}

func consultarAPI() (cotacao.Dados, error) {
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return cotacao.Dados{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return cotacao.Dados{}, fmt.Errorf("timeout de %s excedido: %w", apiTimeout, err)
		}
		return cotacao.Dados{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return cotacao.Dados{}, fmt.Errorf("API externa retornou status %d", res.StatusCode)
	}

	var resposta cotacao.RespostaAPI
	if err := json.NewDecoder(res.Body).Decode(&resposta); err != nil {
		return cotacao.Dados{}, err
	}

	return resposta.USDBRL, nil
}

func (s *Servidor) salvar(dados cotacao.Dados) error {
	ctx, cancel := context.WithTimeout(context.Background(), databaseTimeout)
	defer cancel()

	query := `
		INSERT INTO cotacoes (
			code, codein, name, high, low, var_bid, pct_change, bid, ask, timestamp, create_date
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	_, err := s.db.ExecContext(
		ctx,
		query,
		dados.Code,
		dados.CodeIn,
		dados.Name,
		dados.High,
		dados.Low,
		dados.VarBid,
		dados.PctChange,
		dados.Bid,
		dados.Ask,
		dados.Timestamp,
		dados.CreateDate,
	)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("timeout de %s excedido: %w", databaseTimeout, err)
		}
		return err
	}

	return nil
}

package main

import (
	"log"

	"client-server-api-esdras/internal/cliente"
)

func main() {
	if err := cliente.Executar("http://localhost:8080/cotacao", "cotacao.txt"); err != nil {
		log.Fatalf("erro ao executar cliente: %v", err)
	}
}

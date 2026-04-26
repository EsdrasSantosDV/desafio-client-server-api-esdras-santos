package main

import (
	"log"

	"client-server-api-esdras/internal/servidor"
)

func main() {
	app, err := servidor.Novo(":8080", "cotacoes.db")
	if err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
	defer app.Fechar()

	log.Fatal(app.Iniciar())
}

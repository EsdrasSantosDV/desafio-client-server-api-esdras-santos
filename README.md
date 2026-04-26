# Desafio Client Server API

Aluno: Esdras Santos

Este projeto implementa o desafio de Client Server API em Go. A aplicação possui um servidor HTTP que consulta a cotação do dólar, grava o resultado em SQLite e entrega o valor atual para um cliente. O cliente consulta o servidor e salva a cotação em um arquivo de texto.

## Requisitos atendidos

- Servidor HTTP rodando na porta `8080`.
- Endpoint `GET /cotacao`.
- Consulta da API externa `https://economia.awesomeapi.com.br/json/last/USD-BRL`.
- Timeout de `200ms` para chamada da API externa.
- Persistência da cotação em banco SQLite.
- Timeout de `10ms` para gravação no banco.
- Cliente com timeout de `300ms` para chamar o servidor.
- Geração do arquivo `cotacao.txt` no formato `Dólar: {valor}`.

## Estrutura

```text
.
├── .tool-versions
├── cmd
│   ├── client
│   │   └── main.go
│   └── server
│       └── main.go
├── internal
│   ├── cliente
│   │   └── cliente.go
│   ├── cotacao
│   │   └── cotacao.go
│   └── servidor
│       └── servidor.go
├── go.mod
├── go.sum
└── README.md
```

## Como executar

Entre na pasta do desafio:

```bash
cd "DESAFIO 1 - Client Server Api Esdras Santos"
```

Se ainda não tiver Go configurado no `asdf`, rode:

```bash
asdf plugin add golang
asdf install
```

Inicie o servidor:

```bash
go run ./cmd/server
```

Em outro terminal, execute o cliente:

```bash
go run ./cmd/client
```

Após executar o cliente, será criado o arquivo `cotacao.txt` com o valor atual do dólar.

## Arquivos gerados

- `cotacoes.db`: banco SQLite criado pelo servidor.
- `cotacao.txt`: arquivo criado pelo cliente com a cotação recebida.

## Observações

O servidor precisa estar rodando antes da execução do cliente. Caso algum timeout seja excedido, o erro será exibido no console da aplicação correspondente.

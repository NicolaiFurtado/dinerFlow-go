package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	connStr := "root:root@tcp(127.0.0.1:8889)/diner-flow"
	var err error
	DB, err = sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal("Erro ao conectar com o banco de dados:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Erro ao verificar a conexão com o banco de dados:", err)
	}

	fmt.Println("Conectado ao banco de dados com sucesso.")
}

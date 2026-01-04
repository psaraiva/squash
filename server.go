package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./bin/web"))
	log.Print("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", fs))
}

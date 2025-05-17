package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Обработчик корневого пути
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	addr := ":8080"
	log.Printf("Starting server on %s…", addr)
	// Слушаем и блокируем до ошибки
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

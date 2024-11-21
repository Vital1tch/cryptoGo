package main

import (
	"cryptoGo/gui"
	"log"
	"os"
)

func main() {
	gui.EnsureDirectories()

	// Временная проверка содержимого папки ./encrypted
	files, err := os.ReadDir("./encrypted")
	if err != nil {
		log.Fatalf("Не удалось прочитать папку ./encrypted: %v", err)
	}

	log.Println("Содержимое папки ./encrypted:")
	for _, file := range files {
		log.Printf(" - %s", file.Name())
	}

	gui.StartApp()
}

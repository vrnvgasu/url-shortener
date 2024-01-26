package main

import (
	"fmt"
	"url-shortener/internal/config"
)

func main() {
	// TODO init config: cleanenv - может читать из разных источников
	cnf := config.MustLoad()
	fmt.Println(cnf)

	// TODO init logger: slog - в ядре с версии 1.21

	// TODO init storage: sqlite

	// TODO init router: chi, "chi render"

	// TODO run server
}

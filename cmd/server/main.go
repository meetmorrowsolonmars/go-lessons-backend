package main

import (
	"os"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/app"
)

func main() {
	if err := app.RunServer(); err != nil {
		os.Exit(1)
	}
}

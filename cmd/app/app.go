package main

import (
	"log"

	"github.com/alibek2024/forum/internal/db"
)

func main() {
	_, err := db.InitSQLite("forum.db")
	if err != nil {
		log.Fatal(err)
	}

}

package main

import (
	"fmt"

	"github.com/bnema/gobatterywatcher/db"
	"github.com/bnema/gobatterywatcher/handlers"
	"github.com/bnema/gobatterywatcher/utils"
)

func main() {
	// Read power data from powertop
	data, err := utils.ReadPowerTop()
	if err != nil {
		panic(err)
	}

	// Create db and tables
	database, err := db.CreateDB()
	if err != nil {
		panic(err)
	}
	defer database.Close()

	err = db.CreateTables(database)
	if err != nil {
		panic(err)
	}

	// Process power data
	uniquePowerData := handlers.ProcessData(data)

	// Print unique device names
	for name := range uniquePowerData {
		fmt.Println(name)
	}

	// Insert power data into the database
	err = handlers.InsertData(database, uniquePowerData)
	if err != nil {
		panic(err)
	}
}

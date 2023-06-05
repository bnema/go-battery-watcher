package main

import (
	"log"
	"time"

	"github.com/bnema/gobatterywatcher/cli"
	"github.com/bnema/gobatterywatcher/db"
	"github.com/bnema/gobatterywatcher/handlers"
	"github.com/bnema/gobatterywatcher/utils"
)

func main() {
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

	data, err := handlers.ReadDataLive(database)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ReadDataLive returned %d records\n", len(data))

	// Start a goroutine to continuously refresh the power data
	go func() {
		for {
			// Read power data from powertop
			data, err := utils.ReadPowerTop()
			if err != nil {
				panic(err)
			}

			// Process power data
			uniquePowerData := handlers.ProcessData(data)

			// Insert power data into the database
			err = handlers.InsertData(database, uniquePowerData)
			if err != nil {
				panic(err)
			}

			// Sleep for 5 seconds
			time.Sleep(5 * time.Second)

		}
	}()

	// Start the CLI
	go cli.StartCLI(database)

	// Keep the main function running forever
	select {}
}

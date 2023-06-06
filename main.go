package main

import (
	"time"

	"github.com/bnema/gobatterywatcher/cli"
	"github.com/bnema/gobatterywatcher/db"
	"github.com/bnema/gobatterywatcher/handlers"
	"github.com/bnema/gobatterywatcher/types"
	"github.com/bnema/gobatterywatcher/utils"
	"github.com/distatus/battery"
)

var Battery types.BatteryInfo

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

			// Fetch battery discharge/charge info (value in mW)
			batteryStats, err := battery.GetAll()
			if err != nil {
				panic(err)
			}
			for _, battery := range batteryStats {

				// store the charge rate in the Battery struct
				Battery.Rate = battery.ChargeRate

			}

			// Sleep for 5 seconds
			time.Sleep(5 * time.Second)
		}
	}()

	// Start the CLI
	go cli.StartCLI(database, &Battery)

	// Keep the main function running forever
	select {}
}

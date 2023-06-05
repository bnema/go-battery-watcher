package cli

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	tb "github.com/nsf/termbox-go"

	"github.com/bnema/gobatterywatcher/handlers"
)

func StartCLI(db *sql.DB) {
	// Initialize termui
	if err := ui.Init(); err != nil {
		log.Fatalf("Failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// Create a new line graph widget
	lineGraph := widgets.NewPlot()
	lineGraph.Title = "Battery Usage"
	lineGraph.Data = make([][]float64, 1)
	// Create a new table widget
	table := widgets.NewTable()
	table.Title = "Top 10 Devices by Consumption"
	table.RowStyles[0] = ui.NewStyle(ui.ColorYellow)

	// Define sizes for the widgets
	lineGraph.SetRect(0, 0, 100, 15) // Set the size of the graph to 100x15
	table.SetRect(0, 16, 100, 40)    // Set the size of the table

	// Define a draw function that renders the line graph
	draw := func() {
		ui.Render(lineGraph, table)
	}

	// Initialize termbox
	err := tb.Init()
	if err != nil {
		panic(err)
	}
	defer tb.Close()

	// Create a channel to receive termbox events
	eventQueue := make(chan tb.Event)
	go func() {
		for {
			eventQueue <- tb.PollEvent()
		}
	}()

	// Create a ticker that ticks every 5 seconds
	ticker := time.NewTicker(3 * time.Second).C

	// Main loop
	for {
		select {
		case ev := <-eventQueue:
			// If the user presses Ctrl+C, exit the program
			if ev.Type == tb.EventKey && ev.Key == tb.KeyCtrlC {
				return
			}
		case <-ticker:
			// Read battery data from the database
			data, err := handlers.ReadDataLive(db)
			if err != nil {
				log.Fatal(err)
			}
			// Get top 10 devices by power consumption
			topDevices, err := handlers.GetTopDevices(db)
			if err != nil {
				log.Fatal(err)
			}

			// Prepare data for the table
			rows := make([][]string, len(topDevices)+1)
			// Header row
			rows[0] = []string{"Device Name", "Total Consumption"}

			// Data rows
			for i, device := range topDevices {
				rows[i+1] = []string{device.DeviceName, fmt.Sprintf("%.2f", device.PowerUsage)}
			}

			table.Rows = rows

			// Calculate the total power usage from the battery data
			var totalPower float64
			for _, d := range data {
				totalPower += d.PowerUsage
			}

			// Add the total power usage to the line graph data
			lineGraph.Data[0] = append(lineGraph.Data[0], totalPower, totalPower)

			// If the line graph data length is more than 100, then drop the oldest data point
			if len(lineGraph.Data[0]) > 100 {
				lineGraph.Data[0] = lineGraph.Data[0][1:]
			}

			// Render the widgets
			draw()
		}
	}
}

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
	table := make([]*widgets.Paragraph, 10)
	for i := range table {
		table[i] = widgets.NewParagraph()
		table[i].Border = false
	}

	// Define a draw function that renders the line graph
	draw := func() {
		w, h := ui.TerminalDimensions()
		if h >= 15 { // Check if there is enough space to display at least 5 paragraphs and the line graph
			lineGraph.SetRect(0, 0, w, h/4)
			ui.Render(lineGraph)
			for i := range table {
				table[i].SetRect(0, h/4+i*3, w, h/4+(i+1)*3)
				ui.Render(table[i])
			}
		} else { // Only display the paragraphs
			for i := range table {
				table[i].SetRect(0, h/2+i*3, w, h/2+(i+1)*3)
				ui.Render(table[i])
			}
		}
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
		case e := <-ui.PollEvents():
			if e.Type == ui.ResizeEvent {
				draw()
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
			for i, device := range topDevices {
				table[i].Text = fmt.Sprintf("%s: %.2fW", device.DeviceName, device.PowerUsage)
			}

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

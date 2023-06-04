package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Data struct {
	DeviceType string
	DeviceName string
	PowerUsage string
}

func main() {
	// Read power data from powertop
	data, err := readPowerTop()
	if err != nil {
		panic(err)
	}

	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", "powerdata.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create the power_usage table if it doesn't exist
	createTableQuery := `
        CREATE TABLE IF NOT EXISTS power_usage (
            device_name TEXT PRIMARY KEY,
            power_usage REAL,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `
	_, err = db.Exec(createTableQuery)
	if err != nil {
		panic(err)
	}

	// Combine power data for devices with multiple entries
	uniquePowerData := make(map[string]float64)
	for _, d := range data {
		// Remove PID and numbers from device name
		re := regexp.MustCompile(`\[\s*PID\s*\d+\]`)
		name := re.ReplaceAllString(d.DeviceName, "")
		name = strings.TrimSpace(name)

		// Convert power usage to watts
		powerWatts := convertToWatts(d.PowerUsage)

		// Add power usage to existing total or initialize total to 0
		uniquePowerData[name] += powerWatts

		// Print unique device names
		fmt.Println(name)
	}

	// Insert power data into the database
	insertDataQuery := `
        INSERT OR IGNORE INTO power_usage (device_name, power_usage)
        VALUES (?, ?)
    `
	stmt, err := db.Prepare(insertDataQuery)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	for name, power := range uniquePowerData {
		_, err := stmt.Exec(name, power)
		if err != nil {
			panic(err)
		}
	}
}

// Convert power usage to watts
func convertToWatts(power string) float64 {
	power = strings.TrimSpace(power)
	if strings.HasSuffix(power, "uW") {
		power = strings.TrimSuffix(power, "uW")
		parsed, err := strconv.ParseFloat(power, 64)
		if err != nil {
			return 0
		}
		return parsed / 1000000 // Convert micro-watts to watts
	} else if strings.HasSuffix(power, "mW") {
		power = strings.TrimSuffix(power, "mW")
		parsed, err := strconv.ParseFloat(power, 64)
		if err != nil {
			return 0
		}
		return parsed / 1000 // Convert milliwatts to watts
	} else if strings.HasSuffix(power, "W") {
		power = strings.TrimSuffix(power, "W")
		parsed, err := strconv.ParseFloat(power, 64)
		if err != nil {
			return 0
		}
		return parsed // Already in watts
	}
	return 0
}

// Read power data from powertop
func readPowerTop() ([]Data, error) {
	// Run the powertop command
	cmd := exec.Command("sudo", "powertop", "-C", "powertop.csv", "-t", "3")
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	// Open the CSV file
	f, err := os.Open("powertop.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Prepare the data
	data := make([]Data, 0)

	// Read the file line by line
	scanner := bufio.NewScanner(f)
	reading := false
	for scanner.Scan() {
		line := scanner.Text()

		// Start of sections to read
		if line == "Usage;Wakeups/s;GPU ops/s;Disk IO/s;GFX Wakeups/s;Category;Description;PW Estimate" || line == "Usage;Device Name;PW Estimate" {
			reading = true
			continue
		}

		// End of sections to read
		if line == "____________________________________________________________________" {
			reading = false
			continue
		}

		// Only read lines in sections to read
		if reading {
			fields := strings.Split(line, ";")

			// Make sure there are enough fields
			if len(fields) < 2 {
				continue
			}

			// Ignore if power usage is 0, 0mW, or empty
			powerUsage := fields[len(fields)-1]
			powerUsage = strings.TrimSpace(powerUsage)
			if powerUsage == "" || powerUsage == "0" || powerUsage == "0 mW" {
				continue
			}

			// Add the data
			deviceType := ""
			if len(fields) > 5 {
				deviceType = fields[5]
			}

			deviceName := fields[len(fields)-2]
			data = append(data, Data{
				DeviceType: deviceType,
				DeviceName: deviceName,
				PowerUsage: powerUsage,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

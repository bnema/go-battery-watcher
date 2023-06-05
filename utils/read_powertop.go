package utils

import (
	"bufio"
	"os"
	"os/exec"
	"strings"

	"github.com/bnema/gobatterywatcher/types"
)

// Read power data from powertop
func ReadPowerTop() ([]types.Data, error) {
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
	data := make([]types.Data, 0)

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
			data = append(data, types.Data{
				DeviceType: deviceType,
				DeviceName: deviceName,
				PowerUsage: ConvertToWatts(powerUsage),
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

// Normalize power data
func NormalizeData(data []types.Data, totalPower float64) []types.Data {
	// Calculate the sum of all power usages
	var sumPower float64
	for _, d := range data {
		sumPower += d.PowerUsage
	}

	// Create a new slice of data
	normalizedData := make([]types.Data, len(data))

	// Normalize each power usage
	for i, d := range data {
		normalizedData[i] = d
		normalizedData[i].PowerUsage = d.PowerUsage / sumPower * totalPower
	}

	return normalizedData
}

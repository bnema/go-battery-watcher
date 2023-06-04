package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Data struct {
	DeviceType string
	DeviceName string
	PowerUsage string
}

func main() {
	data, err := readPowerTop()
	if err != nil {
		panic(err)
	}

	for _, d := range data {
		fmt.Printf("Type: %s\nName: %s\nPower: %s\n\n", d.DeviceType, d.DeviceName, d.PowerUsage)
	}
}

func readPowerTop() ([]Data, error) {
	// Exécution de la commande powertop
	cmd := exec.Command("sudo", "powertop", "-C", "powertop.csv", "-t", "3")
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	// Ouvrir le fichier CSV
	f, err := os.Open("powertop.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Préparer les données
	data := make([]Data, 0)

	// Lire le fichier ligne par ligne
	scanner := bufio.NewScanner(f)
	reading := false
	for scanner.Scan() {
		line := scanner.Text()

		// Début des sections à lire
		if line == "Usage;Wakeups/s;GPU ops/s;Disk IO/s;GFX Wakeups/s;Category;Description;PW Estimate" || line == "Usage;Device Name;PW Estimate" {
			reading = true
			continue
		}

		// Fin des sections à lire
		if line == "____________________________________________________________________" {
			reading = false
			continue
		}

		// Lire seulement les lignes dans les sections à lire
		if reading {
			fields := strings.Split(line, ";")

			// S'assurer qu'il y a assez de champs
			if len(fields) < 2 {
				continue
			}

			// Ignorer si la consommation en watts est 0, 0mW ou vide
			powerUsage := fields[len(fields)-1]
			powerUsage = strings.TrimSpace(powerUsage)
			if powerUsage == "" || powerUsage == "0" || powerUsage == "0 mW" {
				continue
			}
			// Ajouter les données
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

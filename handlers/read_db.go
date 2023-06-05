package handlers

import (
	"database/sql"

	"github.com/bnema/gobatterywatcher/types"
)

// Function to read the data live in the table power_usage
func ReadDataLive(db *sql.DB) ([]types.Data, error) {
	rows, err := db.Query("SELECT device_name, power_usage FROM power_usage")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []types.Data

	for rows.Next() {
		var d types.Data
		if err := rows.Scan(&d.DeviceName, &d.PowerUsage); err != nil {
			return nil, err
		}
		data = append(data, d)
	}

	return data, nil
}

// Function to read the data history in the table power_history
// need to match device_id with device_name in power_usage

func ReadDataHistory(db *sql.DB) ([]types.Data, error) {
	query := `
	SELECT 
	    p.device_name, 
	    h.power_usage, 
	    h.timestamp 
	FROM 
	    power_history h
	INNER JOIN 
	    power_usage p
	ON 
	    h.device_id = p.id
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []types.Data
	for rows.Next() {
		var d types.Data
		if err := rows.Scan(&d.DeviceName, &d.PowerUsage, &d.Timestamp); err != nil {
			return nil, err
		}
		data = append(data, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

// GetTopDevices gets the top 10 devices by power consumption from the power_usage table.
func GetTopDevices(db *sql.DB) ([]types.Data, error) {
	rows, err := db.Query(`
		SELECT device_name, power_usage
		FROM power_usage
		ORDER BY power_usage DESC
		LIMIT 10
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []types.Data
	for rows.Next() {
		var device types.Data
		err := rows.Scan(&device.DeviceName, &device.PowerUsage)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

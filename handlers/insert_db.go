package handlers

import (
	"database/sql"
	"regexp"
	"strings"

	"github.com/bnema/gobatterywatcher/types"
)

func ProcessData(data []types.Data) map[string]float64 {
	uniquePowerData := make(map[string]float64)
	for _, d := range data {
		// Remove PID and numbers from device name
		re := regexp.MustCompile(`\[\s*PID\s*\d+\]`)
		name := re.ReplaceAllString(d.DeviceName, "")
		name = strings.TrimSpace(name)

		// Add power usage to existing total or initialize total to 0
		uniquePowerData[name] += d.PowerUsage
	}
	return uniquePowerData
}
func InsertData(db *sql.DB, uniquePowerData map[string]float64) error {
	insertDataQuery := `
	INSERT INTO power_usage (device_name, power_usage)
	VALUES (?, ?)
	ON CONFLICT(device_name) DO UPDATE SET power_usage = excluded.power_usage
	`
	stmt, err := db.Prepare(insertDataQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	insertHistoryQuery := `
	INSERT INTO power_history (device_id, power_usage)
	SELECT id, ?
	FROM power_usage
	WHERE device_name = ?
	`
	stmtHistory, err := db.Prepare(insertHistoryQuery)
	if err != nil {
		return err
	}
	defer stmtHistory.Close()

	for name, power := range uniquePowerData {
		_, err := stmt.Exec(name, power)
		if err != nil {
			return err
		}

		_, err = stmtHistory.Exec(power, name)
		if err != nil {
			return err
		}
	}

	return nil
}

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

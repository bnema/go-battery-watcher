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

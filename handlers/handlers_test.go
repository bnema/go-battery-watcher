package handlers

import (
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/bnema/gobatterywatcher/types"
	"github.com/stretchr/testify/assert"
)

func TestReadDataLive(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"device_name", "power_usage"}).
		AddRow("Device1", 10.5).
		AddRow("Device2", 15.5)

	mock.ExpectQuery("^SELECT device_name, power_usage FROM power_usage$").WillReturnRows(rows)

	result, err := ReadDataLive(db)
	assert.NoError(t, err)

	want := []types.Data{
		{DeviceName: "Device1", PowerUsage: 10.5},
		{DeviceName: "Device2", PowerUsage: 15.5},
	}

	assert.Equal(t, want, result)
}

func TestReadDataHistory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{"device_name", "power_usage", "timestamp"}).
		AddRow("Device1", 10.5, now).
		AddRow("Device2", 15.5, now)

	mock.ExpectQuery("^SELECT (.+) FROM power_history h INNER JOIN power_usage p ON h.device_id = p.id$").WillReturnRows(rows)

	result, err := ReadDataHistory(db)
	assert.NoError(t, err)

	want := []types.Data{
		{DeviceName: "Device1", PowerUsage: 10.5, Timestamp: now},
		{DeviceName: "Device2", PowerUsage: 15.5, Timestamp: now},
	}

	assert.Equal(t, want, result)
	// Ensure that timestamp is not empty
	for _, data := range result {
		assert.False(t, data.Timestamp.IsZero())
	}
}

package handlers

import (
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/bnema/gobatterywatcher/types"
	"github.com/stretchr/testify/assert"
)

func TestReadDataHistory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"device_name", "power_usage", "timestamp"}).
		AddRow("Device1", 10.5, time.Now()).
		AddRow("Device2", 15.5, time.Now())

	mock.ExpectQuery("^SELECT (.+) FROM power_history h INNER JOIN power_usage p ON h.device_id = p.id$").WillReturnRows(rows)

	result, err := ReadDataHistory(db)
	assert.NoError(t, err)

	want := []types.Data{
		// PowerUsage is a string (ex 10w)
		{DeviceName: "Device1", PowerUsage: 10.5},
		{DeviceName: "Device2", PowerUsage: 15.5},
	}

	assert.Equal(t, want, result)
	// ok
}

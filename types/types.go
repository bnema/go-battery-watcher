package types

import "time"

type Data struct {
	DeviceType string
	DeviceName string
	PowerUsage float64
	Timestamp  time.Time
}

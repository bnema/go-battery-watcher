package types

import (
	"sync"
	"time"
)

type Data struct {
	DeviceType string
	DeviceName string
	PowerUsage float64
	Timestamp  time.Time
}

type BatteryInfo struct {
	LastCharge float64
	LastTime   time.Time
	Rate       float64
	RateLock   sync.Mutex
}

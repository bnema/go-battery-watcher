package utils

import (
	"strconv"
	"strings"
)

// Convert power usage to watts
func ConvertToWatts(power string) float64 {
	power = strings.TrimSpace(power)
	power = strings.Replace(power, " ", "", -1)  // remove spaces
	power = strings.Replace(power, ",", ".", -1) // replace commas with dots
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

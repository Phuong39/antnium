package server

import (
	"time"
)

type ClientInfo struct {
	// From every packet
	ComputerId string    `json:"ComputerId"`
	FirstSeen  time.Time `json:"FirstSeen"`
	LastSeen   time.Time `json:"LastSeen"`
	LastIp     string    `json:"LastIp"`

	// From ping
	Hostname string   `json:"Hostname"`
	LocalIps []string `json:"LocalIps"`
}
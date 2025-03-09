package relay

import "time"

type RelayInfo struct {
	RelayMode         int
	BaseURL           string
	AlternateBaseURL  string
	APIKey            string
	APISecret         string
	APIToken          string
	IsStream          bool
	IsPlayground      bool
	StartTime         time.Time
	FirstResponseTime time.Time
	SetFirstResponse  bool
	RequestModelName  string
}

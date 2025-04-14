package structure

type ServerConfig struct {
	Port string `json:"port"`
	SessionActiveInterval int `json:"session_active_interval_hours"`
	RateLimitMaxToken int `json:"rate_limit_max_sessions"`
	RateLimitRefillDuration int `json:"rate_limit_refill_duration_secods"`
}
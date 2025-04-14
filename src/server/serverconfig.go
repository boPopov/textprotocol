package server

import(
	"fmt"
)

type Configer interface {
	Load(path string) error
	Print()
}

type ServerConfig struct {
	Port string `json:"port"`
	SessionActiveInterval int `json:"session_active_interval_hours"`
	RateLimitMaxSessions int `json:"rate_limit_max_sessions"`
	RateLimitRefillDuration int `json:"rate_limit_refill_duration_secods"`
	RateLimitMaxInputPerInterval int `json:"rate_limit_max_input_per_interval"`
	Configer
}

func (s *ServerConfig) Load(path) error {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err 
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&s)
	return err
}

func (s *ServerConfig) Print() {
	fmt.Println(fmt.Sprintf("Port: %s, SessionActiveInterval: %d,RateLimitMaxSessions: %d,RateLimitRefillDuration: %d, RateLimitMaxInputPerInterval: %d", s.Port, s.SessionActiveInterval, s.RateLimitMaxSessions, s.RateLimitRefillDuration, s.RateLimitMaxInputPerInterval)))
}

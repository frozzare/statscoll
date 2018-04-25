package stat

// Stat represents a stat value.
type Stat struct {
	ID        uint   `gorm:"primary_key" json:"-"`
	Count     int64  `json:"count,omitempty"`
	Metric    string `json:"metric,omitempty"`
	Project   string `json:"project,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

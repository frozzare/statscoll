package stat

// Stat represents a stat value.
type Stat struct {
	ID        uint   `gorm:"primary_key" json:"-"`
	Count     int64  `json:"count"`
	Metric    string `json:"metric"`
	Project   string `json:"project"`
	Timestamp int64  `json:"timestamp"`
}

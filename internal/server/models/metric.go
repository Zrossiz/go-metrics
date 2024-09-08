package models

import "time"

type Metric struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"column:name"`
	Type      string    `gorm:"column:metric_type"`
	Value     float64   `gorm:"column:value"`
	Delta     int64     `gorm:"column:delta"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

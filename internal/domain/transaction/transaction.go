package transaction

import (
    "time"
)

type Transaction struct {
    UserId      string    `gorm:"type:varchar(26);index;not null"`
    Id          string    `gorm:"type:varchar(26);primaryKey"`
    Type        Types     `gorm:"type:varchar(10);not null"`
    CategoryId  string    `gorm:"type:varchar(26);index"`
    Amount      float64   `gorm:"not null"`
    Description string    `gorm:"size:255"`
    Date        time.Time `gorm:"not null"`
    CreatedAt   time.Time `gorm:"not null"`
    UpdatedAt   time.Time `gorm:"not null"`
}

type Category struct {
    UserId string `gorm:"type:varchar(26);index;not null"`
    Id     string `gorm:"type:varchar(26);primaryKey"`
    Name   string `gorm:"size:100;not null"`
    Icon   string `gorm:"size:50"`
}

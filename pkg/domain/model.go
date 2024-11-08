package domain

import "time"

type User struct {
	ID        uint       `gorm:"primaryKey"`
	Username  string     `gorm:"unique;not null"`
	Email     string     `gorm:"unique;not null"`
	Password  string     `gorm:"not null"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	Accounts  []Account  `gorm:"foreignKey:UserID"`
	Orders    []Order    `gorm:"foreignKey:UserID"`
	Positions []Position `gorm:"foreignKey:UserID"`
	Trades    []Trade    `gorm:"foreignKey:UserID"`
}

// this is for adding multiple accounts for the user
type Account struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	Balance   float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// Pending ===== order
type Order struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	Symbol    string    `gorm:"not null"`
	Volume    float64   `gorm:"not null"`
	Type      string    `gorm:"not null"`
	Price     float64   `gorm:"not null"`
	Status    string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Position struct {
	ID            uint      `gorm:"primaryKey"`
	UserID        uint      `gorm:"not null;index"`
	Symbol        string    `gorm:"not null"`
	Volume        float64   `gorm:"not null"`
	EntryPrice    float64   `gorm:"not null"`
	UnrealizedPnl float64   `gorm:"not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

type Trade struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	Symbol    string    `gorm:"not null"`
	Volume    float64   `gorm:"not null"`
	Price     float64   `gorm:"not null"`
	Type      string    `gorm:"not null"`
	Timestamp time.Time `gorm:"not null"`
}

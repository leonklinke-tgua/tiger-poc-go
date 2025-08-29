package entities

import "time"

type PolicyType string

const (
	LG  PolicyType = "lg"  // Lease Guarantee
	SDR PolicyType = "sdr" // Security Deposit Replacement
)

type Policy struct {
	ID                string     `json:"id"`
	UserID            string     `json:"user_id"`
	User              *User      `json:"user"`
	PolicyType        PolicyType `json:"policy_type"`
	PolicyAmountCents int64      `json:"policy_amount_cents"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

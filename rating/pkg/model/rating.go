package model

// RecordID defines a record id. Together
// with RecordType identifies unique records across all types.
type RecordID string

// RecordType defines a record type. Together
// with RecordID identifies unique records across all types.
type RecordType string

// RecordTypeMovie is existing record types.
const (
	RecordTypeMovie = RecordType("movie")
)

// UserID is user identification number.
type UserID string

// RatingValue defines score set by user.
type RatingValue int

// Rating defines an individual rating set by user.
type Rating struct {
	RecordID   RecordID    `json:"recordId"`
	RecordType RecordType  `json:"recordType"`
	UserID     UserID      `json:"userId"`
	Value      RatingValue `json:"value"`
}

package constant

//go:generate go-enum --marshal --sql --values --names --file

// ENUM(pass, like)
type UserMatchType string

// List of internal constant for user match
const (
	MaxMatchPerDay = 10
)

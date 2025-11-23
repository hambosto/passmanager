package entity

import "time"

// EntryType represents the type of vault entry
type EntryType int

const (
	EntryTypeLogin EntryType = iota
	EntryTypeSecureNote
	EntryTypeCard
	EntryTypeIdentity
)

// String returns the string representation of the entry type
func (et EntryType) String() string {
	switch et {
	case EntryTypeLogin:
		return "Login"
	case EntryTypeSecureNote:
		return "Secure Note"
	case EntryTypeCard:
		return "Card"
	case EntryTypeIdentity:
		return "Identity"
	default:
		return "Unknown"
	}
}

// Entry represents a vault entry with credentials and metadata
type Entry struct {
	ID           string            `json:"id"`
	Type         EntryType         `json:"type"`
	Name         string            `json:"name"`
	Username     string            `json:"username,omitempty"`
	Password     string            `json:"password,omitempty"`
	URI          string            `json:"uri,omitempty"`
	Notes        string            `json:"notes,omitempty"`
	TOTPSecret   string            `json:"totp_secret,omitempty"`
	CustomFields map[string]string `json:"custom_fields,omitempty"`
	FolderID     string            `json:"folder_id,omitempty"`
	IsFavorite   bool              `json:"is_favorite"`
	Tags         []string          `json:"tags,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	AccessedAt   time.Time         `json:"accessed_at,omitempty"`
}

// Card represents credit card information
type Card struct {
	CardholderName string `json:"cardholder_name"`
	Number         string `json:"number"`
	Brand          string `json:"brand"`
	ExpMonth       string `json:"exp_month"`
	ExpYear        string `json:"exp_year"`
	CVV            string `json:"cvv"`
}

// Identity represents personal identity information
type Identity struct {
	Title      string `json:"title"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	Address1   string `json:"address1"`
	Address2   string `json:"address2"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	SSN        string `json:"ssn"`
	PassportNo string `json:"passport_number"`
}

// NewEntry creates a new entry with generated ID and timestamps
func NewEntry(entryType EntryType, name string) *Entry {
	now := time.Now()
	return &Entry{
		ID:           generateID(),
		Type:         entryType,
		Name:         name,
		CustomFields: make(map[string]string),
		Tags:         make([]string, 0),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// UpdateAccessTime updates the last accessed timestamp
func (e *Entry) UpdateAccessTime() {
	e.AccessedAt = time.Now()
}

// Update updates the entry and sets the updated timestamp
func (e *Entry) Update() {
	e.UpdatedAt = time.Now()
}

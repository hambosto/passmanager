package dto

// CreateEntryRequest represents a request to create a new entry
type CreateEntryRequest struct {
	Name       string
	Username   string
	Password   string
	URI        string
	TOTPSecret string
	Notes      string
	FolderID   string
	IsFavorite bool
}

// UpdateEntryRequest represents a request to update an entry
type UpdateEntryRequest struct {
	ID         string
	Name       string
	Username   string
	Password   string
	URI        string
	TOTPSecret string
	Notes      string
	FolderID   string
	IsFavorite bool
}

// DeleteEntryRequest represents a request to delete an entry
type DeleteEntryRequest struct {
	ID string
}

// SearchRequest represents a search request
type SearchRequest struct {
	Query    string
	FolderID string
	Type     string
}

// GeneratePasswordRequest represents a password generation request
type GeneratePasswordRequest struct {
	Length           int
	IncludeUppercase bool
	IncludeLowercase bool
	IncludeNumbers   bool
	IncludeSymbols   bool
	ExcludeAmbiguous bool
}

// GeneratePassphraseRequest represents a passphrase generation request
type GeneratePassphraseRequest struct {
	WordCount     int
	Separator     string
	Capitalize    bool
	IncludeNumber bool
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	OldPassword string
	NewPassword string
}

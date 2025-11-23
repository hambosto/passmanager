package entity

import "time"

// Folder represents an organizational folder for entries
type Folder struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ParentID  string    `json:"parent_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// NewFolder creates a new folder with generated ID and timestamp
func NewFolder(name string, parentID string) *Folder {
	return &Folder{
		ID:        generateID(),
		Name:      name,
		ParentID:  parentID,
		CreatedAt: time.Now(),
	}
}

// IsRoot returns true if the folder is a root folder (no parent)
func (f *Folder) IsRoot() bool {
	return f.ParentID == ""
}

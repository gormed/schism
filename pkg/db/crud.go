package db

// Identifyable defines the minimum CRUD resource
type Identifyable struct {
	Database *Sqlite `json:"-"`
	Id       *string `json:"id"`
}

// CRUD defines the interface to the database of a CRUD resource
type CRUD interface {
	Exists() (bool, error)

	Create() (*Identifyable, int, error)
	Read() (*Identifyable, int, error)
	Update() (*Identifyable, int, error)
	Delete() (*Identifyable, int, error)
}

package db

// SqlIdentifyable defines the minimum CRUD resource
type SqlIdentifyable struct {
	Database *Sqlite `json:"-"`
	Id       *string `json:"id"`
}

// CRUD defines the interface to the database of a CRUD resource
type CRUD interface {
	Exists() (bool, error)

	Create() (*SqlIdentifyable, int, error)
	Read() (*SqlIdentifyable, int, error)
	Update() (*SqlIdentifyable, int, error)
	Delete() (*SqlIdentifyable, int, error)
}

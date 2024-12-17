package service

// Environment represents an environment within an application.
// It holds details about the environment, such as its ID, associated application ID,
// level, name, and timestamps for when it was created, updated, and optionally deleted.
type Environment struct {
	// ID is the unique identifier of the environment.
	ID int64 `json:"id"`

	// ApplicationID is the identifier of the application to which this environment belongs.
	ApplicationID int64 `json:"applicationId"`

	// Level indicates the environment's level, which might be used to denote the hierarchy or order of environments.
	Level int `json:"level"`

	// Name is the name of the environment.
	Name string `json:"name"`
}

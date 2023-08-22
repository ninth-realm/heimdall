package store

import "fmt"

type NotFoundError struct {
	ResourceType string
	ResourceID   string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.ResourceType, e.ResourceID)
}

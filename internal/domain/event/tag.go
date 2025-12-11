package event

type Tag struct {
	TagID int64  `json:"tag_id"`
	Name  string `json:"name"`
}

type TagRepo interface {
	FindAll() ([]Tag, error)
	FindByName(name string) (*Tag, error)
	CreateIfNotExists(name string) (*Tag, error)
	SeedDefaults() error
}

var DefaultTagNames = []string{
	"Music",
	"Sports",
	"Food & Drink",
	"Networking",
	"Workshop",
	"Party",
	"Conference",
	"Meetup",
	"Art",
	"Charity",
	"Outdoor",
	"Gaming",
	"Tech",
	"Health & Wellness",
	"Education",
}

package repositories

import "errors"

type StorageType int64

const (
	InMemory StorageType = iota
	Mongo
)

var StorageTypeValues = []StorageType{
	InMemory,
	Mongo,
}

func (s StorageType) String() string {
	switch s {
	case InMemory:
		return "inmemory"
	case Mongo:
		return "mongo"
	default:
		return "unknown"
	}
}

func ToStorageType(s string) (StorageType, error) {
	switch s {
	case "inmemory":
		return InMemory, nil
	case "mongo":
		return Mongo, nil
	default:
		return -1, errors.New("unknown storage type")
	}
}

package repositories

import "errors"

type StorageType int64

const (
	InMemory StorageType = iota
)

var StorageTypeValues = []StorageType{
	InMemory,
}

func (s StorageType) String() string {
	switch s {
	case InMemory:
		return "inmemory"
	default:
		return "unknown"
	}
}

func ToStorageType(s string) (StorageType, error) {
	switch s {
	case "inmemory":
		return InMemory, nil
	default:
		return -1, errors.New("unknown storage type")
	}
}

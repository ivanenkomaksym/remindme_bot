package repositories

// UserRepositoryFactory creates repository instances based on storage type
type UserRepositoryFactory struct{}

// NewUserRepositoryFactory creates a new factory instance
func NewUserRepositoryFactory() *UserRepositoryFactory {
	return &UserRepositoryFactory{}
}

// CreateRepository creates a repository instance based on the storage type
func (f *UserRepositoryFactory) CreateRepository(storageType StorageType) UserRepository {
	switch storageType {
	case InMemory:
		return NewInMemoryUserRepository()
	default:
		// Default to in-memory if unknown type
		return NewInMemoryUserRepository()
	}
}

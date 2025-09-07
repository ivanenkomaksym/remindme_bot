package repositories

// ReminderRepositoryFactory creates repository instances based on storage type
type ReminderRepositoryFactory struct{}

// NewReminderRepositoryFactory creates a new factory instance
func NewReminderRepositoryFactory() *ReminderRepositoryFactory {
	return &ReminderRepositoryFactory{}
}

// CreateRepository creates a repository instance based on the storage type
func (f *ReminderRepositoryFactory) CreateRepository(storageType StorageType) ReminderRepository {
	switch storageType {
	case InMemory:
		return NewInMemoryReminderRepository()
	default:
		// Default to in-memory if unknown type
		return NewInMemoryReminderRepository()
	}
}

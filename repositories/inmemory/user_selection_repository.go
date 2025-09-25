package inmemory

import (
	"sync"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/repositories"
)

type InMemoryUserSelectionRepository struct {
	mu             sync.RWMutex
	userSelections map[int64]*entities.UserSelection
}

func NewInMemoryUserSelectionRepository() repositories.UserSelectionRepository {
	return &InMemoryUserSelectionRepository{userSelections: make(map[int64]*entities.UserSelection)}
}

func (r *InMemoryUserSelectionRepository) GetUserSelection(userID int64) (*entities.UserSelection, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	selection, exists := r.userSelections[userID]
	if !exists {
		return nil, nil
	}

	selectionCopy := *selection
	return &selectionCopy, nil
}

func (r *InMemoryUserSelectionRepository) SetUserSelection(userID int64, selection *entities.UserSelection) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	selectionCopy := *selection
	r.userSelections[userID] = &selectionCopy
	return nil
}

func (r *InMemoryUserSelectionRepository) UpdateUserSelection(userID int64, selection *entities.UserSelection) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	selectionCopy := *selection
	r.userSelections[userID] = &selectionCopy
	return nil
}

func (r *InMemoryUserSelectionRepository) ClearUserSelection(userID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.userSelections, userID)
	return nil
}

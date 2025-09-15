package usecases

import (
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/errors"
	"github.com/ivanenkomaksym/remindme_bot/domain/repositories"
)

// ReminderUseCase defines the interface for reminder business logic
type ReminderUseCase interface {
	CreateReminder(userID int64, selection *entities.UserSelection) (*entities.Reminder, error)
	GetUserReminders(userID int64) ([]entities.Reminder, error)
	GetAllReminders() ([]entities.Reminder, error)
	DeleteReminder(reminderID, userID int64) error
	UpdateReminder(reminder *entities.Reminder) error
	GetActiveReminders() ([]entities.Reminder, error)
}

type reminderUseCase struct {
	reminderRepo repositories.ReminderRepository
	userRepo     repositories.UserRepository
}

// NewReminderUseCase creates a new reminder use case
func NewReminderUseCase(reminderRepo repositories.ReminderRepository, userRepo repositories.UserRepository) ReminderUseCase {
	return &reminderUseCase{
		reminderRepo: reminderRepo,
		userRepo:     userRepo,
	}
}

func (r *reminderUseCase) CreateReminder(userID int64, selection *entities.UserSelection) (*entities.Reminder, error) {
	if userID <= 0 {
		return nil, errors.NewDomainError("INVALID_USER_ID", "User ID must be positive", nil)
	}
	if selection == nil {
		return nil, errors.NewDomainError("INVALID_SELECTION", "User selection cannot be nil", nil)
	}
	if selection.ReminderMessage == "" {
		return nil, errors.ErrEmptyMessage
	}
	if selection.SelectedTime == "" {
		return nil, errors.ErrInvalidTimeFormat
	}

	// Get user to ensure they exist
	user, err := r.userRepo.GetUser(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	// Create reminder based on recurrence type
	switch selection.RecurrenceType {
	case entities.Once:
		return r.createOnceReminder(user, selection)
	case entities.Daily:
		return r.createDailyReminder(user, selection)
	case entities.Weekly:
		return r.createWeeklyReminder(user, selection)
	case entities.Monthly:
		return r.createMonthlyReminder(user, selection)
	case entities.Interval, entities.Custom:
		// Treat as daily for now
		return r.createDailyReminder(user, selection)
	default:
		return nil, errors.ErrInvalidRecurrenceType
	}
}

func (r *reminderUseCase) createOnceReminder(user *entities.User, selection *entities.UserSelection) (*entities.Reminder, error) {
	reminder, err := r.reminderRepo.CreateOnceReminder(selection.SelectedTime, user, selection.ReminderMessage)
	if err != nil {
		return nil, err
	}

	return reminder, nil
}

func (r *reminderUseCase) createDailyReminder(user *entities.User, selection *entities.UserSelection) (*entities.Reminder, error) {
	reminder, err := r.reminderRepo.CreateDailyReminder(selection.SelectedTime, user, selection.ReminderMessage)
	if err != nil {
		return nil, err
	}

	return reminder, nil
}

func (r *reminderUseCase) createWeeklyReminder(user *entities.User, selection *entities.UserSelection) (*entities.Reminder, error) {
	// Convert week options to time.Weekday slice
	var daysOfWeek []time.Weekday
	for i, selected := range selection.WeekOptions {
		if selected {
			daysOfWeek = append(daysOfWeek, time.Weekday(i))
		}
	}

	if len(daysOfWeek) == 0 {
		return nil, errors.NewDomainError("NO_WEEKDAYS_SELECTED", "At least one weekday must be selected", nil)
	}

	reminder, err := r.reminderRepo.CreateWeeklyReminder(daysOfWeek, selection.SelectedTime, user, selection.ReminderMessage)
	if err != nil {
		return nil, err
	}

	return reminder, nil
}

func (r *reminderUseCase) createMonthlyReminder(user *entities.User, selection *entities.UserSelection) (*entities.Reminder, error) {
	// For monthly, use the 1st of each month for now
	daysOfMonth := []int{1}

	reminder, err := r.reminderRepo.CreateMonthlyReminder(daysOfMonth, selection.SelectedTime, user, selection.ReminderMessage)
	if err != nil {
		return nil, err
	}

	return reminder, nil
}

func (r *reminderUseCase) GetUserReminders(userID int64) ([]entities.Reminder, error) {
	if userID <= 0 {
		return nil, errors.NewDomainError("INVALID_USER_ID", "User ID must be positive", nil)
	}

	reminders, err := r.reminderRepo.GetRemindersByUser(userID)
	if err != nil {
		return nil, err
	}
	return reminders, nil
}

func (r *reminderUseCase) GetAllReminders() ([]entities.Reminder, error) {
	reminders, err := r.reminderRepo.GetReminders()
	if err != nil {
		return nil, err
	}
	return reminders, nil
}

func (r *reminderUseCase) DeleteReminder(reminderID, userID int64) error {
	if reminderID <= 0 {
		return errors.NewDomainError("INVALID_REMINDER_ID", "Reminder ID must be positive", nil)
	}
	if userID <= 0 {
		return errors.NewDomainError("INVALID_USER_ID", "User ID must be positive", nil)
	}

	// Check if reminder exists and belongs to user
	reminder, err := r.reminderRepo.GetReminder(reminderID)
	if err != nil {
		return err
	}
	if reminder == nil {
		return errors.ErrReminderNotFound
	}
	if reminder.UserID != userID {
		return errors.ErrUnauthorized
	}

	return r.reminderRepo.DeleteReminder(reminderID, userID)
}

func (r *reminderUseCase) UpdateReminder(reminder *entities.Reminder) error {
	if reminder == nil {
		return errors.NewDomainError("INVALID_REMINDER", "Reminder cannot be nil", nil)
	}
	if reminder.ID <= 0 {
		return errors.NewDomainError("INVALID_REMINDER_ID", "Reminder ID must be positive", nil)
	}

	return r.reminderRepo.UpdateReminder(reminder)
}

func (r *reminderUseCase) GetActiveReminders() ([]entities.Reminder, error) {
	reminders, err := r.reminderRepo.GetActiveReminders()
	if err != nil {
		return nil, err
	}
	return reminders, nil
}

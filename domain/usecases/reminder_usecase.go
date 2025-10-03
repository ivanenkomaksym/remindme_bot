package usecases

import (
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/errors"
	"github.com/ivanenkomaksym/remindme_bot/domain/repositories"
	"github.com/ivanenkomaksym/remindme_bot/scheduler"
)

// ReminderUseCase defines the interface for reminder business logic
type ReminderUseCase interface {
	CreateReminder(userID int64, selection *entities.UserSelection) (*entities.Reminder, error)
	GetUserReminders(userID int64) ([]entities.Reminder, error)
	GetReminder(userID, reminderID int64) (*entities.Reminder, error)
	GetAllReminders() ([]entities.Reminder, error)
	DeleteReminder(reminderID, userID int64) error
	UpdateReminder(userID, reminderID int64, reminder *entities.Reminder) (*entities.Reminder, error)
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

// helper: merge a selected date and time string into a time.Time in user's location
func buildDateTimeFromSelection(date time.Time, selectedTime string, loc *time.Location) (time.Time, error) {
	if loc == nil {
		loc = time.UTC
	}
	date = date.In(loc)
	hour, minute, ok := scheduler.ParseHourMinute(selectedTime)
	if !ok {
		return time.Time{}, errors.ErrInvalidTimeFormat
	}
	return time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, loc), nil
}

func (r *reminderUseCase) CreateReminder(userID int64, selection *entities.UserSelection) (*entities.Reminder, error) {
	// Get user to ensure they exist
	user, err := r.userRepo.GetUser(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	timeOfDay := time.Time{}

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
	} else {
		loc := user.GetLocation()
		timeOfDay, err = buildDateTimeFromSelection(selection.SelectedDate, selection.SelectedTime, loc)
		if err != nil {
			return nil, err
		}
	}

	// Create reminder based on recurrence type
	switch selection.RecurrenceType {
	case entities.Once:
		return r.createOnceReminder(user, selection, timeOfDay)
	case entities.Daily:
		return r.createDailyReminder(user, selection, timeOfDay)
	case entities.Weekly:
		return r.createWeeklyReminder(user, selection, timeOfDay)
	case entities.Monthly:
		return r.createMonthlyReminder(user, selection, timeOfDay)
	case entities.Interval:
		return r.createIntervalReminder(user, selection, timeOfDay)
	case entities.SpacedBasedRepetition:
		return r.createSpaceBasedRepetitionReminder(user, selection, timeOfDay)
	default:
		return nil, errors.ErrInvalidRecurrenceType
	}
}

func (r *reminderUseCase) createOnceReminder(user *entities.User, selection *entities.UserSelection, timeOfDay time.Time) (*entities.Reminder, error) {
	reminder, err := r.reminderRepo.CreateOnceReminder(timeOfDay, user, selection.ReminderMessage)
	if err != nil {
		return nil, err
	}

	return reminder, nil
}

func (r *reminderUseCase) createDailyReminder(user *entities.User, selection *entities.UserSelection, timeOfDay time.Time) (*entities.Reminder, error) {
	reminder, err := r.reminderRepo.CreateDailyReminder(timeOfDay, user, selection.ReminderMessage)
	if err != nil {
		return nil, err
	}

	return reminder, nil
}

func (r *reminderUseCase) createWeeklyReminder(user *entities.User, selection *entities.UserSelection, timeOfDay time.Time) (*entities.Reminder, error) {
	if len(selection.WeekOptions) == 0 {
		return nil, errors.NewDomainError("NO_WEEKDAYS_SELECTED", "At least one weekday must be selected", nil)
	}

	reminder, err := r.reminderRepo.CreateWeeklyReminder(selection.WeekOptions, timeOfDay, user, selection.ReminderMessage)
	if err != nil {
		return nil, err
	}

	return reminder, nil
}

func (r *reminderUseCase) createMonthlyReminder(user *entities.User, selection *entities.UserSelection, timeOfDay time.Time) (*entities.Reminder, error) {
	if len(selection.MonthOptions) == 0 {
		return nil, errors.NewDomainError("NO_DAYS_SELECTED", "At least one day of month must be selected", nil)
	}

	reminder, err := r.reminderRepo.CreateMonthlyReminder(selection.MonthOptions, timeOfDay, user, selection.ReminderMessage)
	if err != nil {
		return nil, err
	}

	return reminder, nil
}

func (r *reminderUseCase) createIntervalReminder(user *entities.User, selection *entities.UserSelection, timeOfDay time.Time) (*entities.Reminder, error) {
	if selection.IntervalDays <= 0 {
		return nil, errors.NewDomainError("INVALID_INTERVAL", "Interval must be a positive number of days", nil)
	}
	reminder, err := r.reminderRepo.CreateIntervalReminder(selection.IntervalDays, timeOfDay, user, selection.ReminderMessage)
	if err != nil {
		return nil, err
	}
	return reminder, nil
}

func (r *reminderUseCase) createSpaceBasedRepetitionReminder(user *entities.User, selection *entities.UserSelection, timeOfDay time.Time) (*entities.Reminder, error) {
	reminder, err := r.reminderRepo.CreateSpaceBasedRepetitionReminder(timeOfDay, user, selection.ReminderMessage)
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

func (r *reminderUseCase) GetReminder(userID, reminderID int64) (*entities.Reminder, error) {
	if reminderID <= 0 {
		return nil, errors.NewDomainError("INVALID_REMINDER_ID", "Reminder ID must be positive", nil)
	}
	if userID <= 0 {
		return nil, errors.NewDomainError("INVALID_USER_ID", "User ID must be positive", nil)
	}

	// Get reminder
	reminder, err := r.reminderRepo.GetReminder(reminderID)
	if err != nil {
		return nil, err
	}
	if reminder == nil {
		return nil, errors.ErrReminderNotFound
	}

	// Verify ownership
	if reminder.UserID != userID {
		return nil, errors.ErrUnauthorized
	}

	return reminder, nil
}

func (r *reminderUseCase) UpdateReminder(userID, reminderID int64, updatedFields *entities.Reminder) (*entities.Reminder, error) {
	if updatedFields == nil {
		return nil, errors.NewDomainError("INVALID_REMINDER", "Reminder cannot be nil", nil)
	}
	if reminderID <= 0 {
		return nil, errors.NewDomainError("INVALID_REMINDER_ID", "Reminder ID must be positive", nil)
	}
	if userID <= 0 {
		return nil, errors.NewDomainError("INVALID_USER_ID", "User ID must be positive", nil)
	}

	// Get existing reminder to verify ownership and use as base for update
	existingReminder, err := r.reminderRepo.GetReminder(reminderID)
	if err != nil {
		return nil, err
	}
	if existingReminder == nil {
		return nil, errors.ErrReminderNotFound
	}
	if existingReminder.UserID != userID {
		return nil, errors.ErrUnauthorized
	}

	// Update only fields that are provided
	if updatedFields.Message != "" {
		existingReminder.Message = updatedFields.Message
	}
	if updatedFields.NextTrigger != nil {
		existingReminder.NextTrigger = updatedFields.NextTrigger
	}
	if updatedFields.Recurrence != nil {
		existingReminder.Recurrence = updatedFields.Recurrence
	}

	// Only update IsActive if it's different from the existing value
	if updatedFields.IsActive != existingReminder.IsActive {
		existingReminder.IsActive = updatedFields.IsActive
	}

	// Keep the original ID, UserID and CreatedAt
	existingReminder.ID = reminderID
	existingReminder.UserID = userID
	existingReminder.CreatedAt = existingReminder.CreatedAt

	// Update the reminder
	err = r.reminderRepo.UpdateReminder(existingReminder)
	if err != nil {
		return nil, err
	}

	return existingReminder, nil
}

func (r *reminderUseCase) GetActiveReminders() ([]entities.Reminder, error) {
	reminders, err := r.reminderRepo.GetActiveReminders()
	if err != nil {
		return nil, err
	}
	return reminders, nil
}

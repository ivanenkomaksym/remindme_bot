package persistent

import (
	"context"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/repositories"
	"github.com/ivanenkomaksym/remindme_bot/scheduler"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoReminderRepository struct {
	client   *mongo.Client
	database string
	col      *mongo.Collection
}

func NewMongoReminderRepository(connectionString string, database string) (repositories.ReminderRepository, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	client, err := mongo.Connect(options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPI))
	if err != nil {
		return nil, err
	}
	db := client.Database(database)
	return &MongoReminderRepository{
		client:   client,
		database: database,
		col:      db.Collection("reminders"),
	}, nil
}

func (r *MongoReminderRepository) CreateOnceReminder(date time.Time, timeStr string, user *entities.User, message string) (*entities.Reminder, error) {
	recurrence := entities.OnceAt(date, timeStr, user.Location)
	rem := entities.NewReminder(0, user.ID, message, recurrence, recurrence.StartDate)
	return r.insertAndReturn(rem)
}

func (r *MongoReminderRepository) CreateDailyReminder(timeStr string, user *entities.User, message string) (*entities.Reminder, error) {
	now := time.Now()
	recurrence := entities.DailyAt(timeStr, user.Location)
	next := scheduler.NextDailyTrigger(now, timeStr, user.Location)
	rem := entities.NewReminder(0, user.ID, message, recurrence, &next)
	return r.insertAndReturn(rem)
}

func (r *MongoReminderRepository) CreateWeeklyReminder(daysOfWeek []time.Weekday, timeStr string, user *entities.User, message string) (*entities.Reminder, error) {
	now := time.Now()
	next := scheduler.NextWeeklyTrigger(now, daysOfWeek, timeStr, user.Location)
	rem := entities.NewReminder(0, user.ID, message, entities.CustomWeekly(daysOfWeek, timeStr, user.Location), &next)
	return r.insertAndReturn(rem)
}

func (r *MongoReminderRepository) CreateMonthlyReminder(daysOfMonth []int, timeStr string, user *entities.User, message string) (*entities.Reminder, error) {
	now := time.Now()
	next := scheduler.NextMonthlyTrigger(now, daysOfMonth, timeStr, user.Location)
	rem := entities.NewReminder(0, user.ID, message, entities.MonthlyOnDay(daysOfMonth, timeStr, user.Location), &next)
	return r.insertAndReturn(rem)
}

func (r *MongoReminderRepository) CreateIntervalReminder(intervalDays int, timeStr string, user *entities.User, message string) (*entities.Reminder, error) {
	now := time.Now()
	base := scheduler.NextDailyTrigger(now, timeStr, user.Location)
	next := base.Add(time.Duration(intervalDays-1) * 24 * time.Hour)
	rem := entities.NewReminder(0, user.ID, message, entities.IntervalEveryDays(intervalDays, timeStr, user.Location), &next)
	return r.insertAndReturn(rem)
}

func (r *MongoReminderRepository) insertAndReturn(rem *entities.Reminder) (*entities.Reminder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Let Mongo assign _id; store it in ID as an increment is not available; keep ID as 0 or derive from _id if needed in the future
	// If ID is zero, assign a time-based unique ID to fit int64
	if rem.ID == 0 {
		rem.ID = time.Now().UnixNano()
	}
	_, err := r.col.InsertOne(ctx, rem)
	if err != nil {
		return nil, err
	}
	return rem, nil
}

func (r *MongoReminderRepository) GetReminders() ([]entities.Reminder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := r.col.Find(ctx, struct{}{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var res []entities.Reminder
	for cur.Next(ctx) {
		var rm entities.Reminder
		if err := cur.Decode(&rm); err != nil {
			return nil, err
		}
		res = append(res, rm)
	}
	return res, cur.Err()
}

func (r *MongoReminderRepository) GetRemindersByUser(userID int64) ([]entities.Reminder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := r.col.Find(ctx, map[string]any{"userId": userID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var res []entities.Reminder
	for cur.Next(ctx) {
		var rm entities.Reminder
		if err := cur.Decode(&rm); err != nil {
			return nil, err
		}
		res = append(res, rm)
	}
	return res, cur.Err()
}

func (r *MongoReminderRepository) GetReminder(reminderID int64) (*entities.Reminder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var rm entities.Reminder
	err := r.col.FindOne(ctx, map[string]any{"id": reminderID}).Decode(&rm)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rm, nil
}

func (r *MongoReminderRepository) UpdateReminder(reminder *entities.Reminder) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.col.UpdateOne(ctx, map[string]any{"id": reminder.ID}, map[string]any{"$set": reminder})
	return err
}

func (r *MongoReminderRepository) DeleteReminder(reminderID int64, userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.col.DeleteOne(ctx, map[string]any{"id": reminderID, "userId": userID})
	return err
}

func (r *MongoReminderRepository) DeactivateReminder(reminderID int64, userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.col.UpdateOne(ctx, map[string]any{"id": reminderID, "userId": userID}, map[string]any{"$set": map[string]any{"isActive": false}})
	return err
}

func (r *MongoReminderRepository) GetActiveReminders() ([]entities.Reminder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := r.col.Find(ctx, map[string]any{"isActive": true})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var res []entities.Reminder
	for cur.Next(ctx) {
		var rm entities.Reminder
		if err := cur.Decode(&rm); err != nil {
			return nil, err
		}
		res = append(res, rm)
	}
	return res, cur.Err()
}

func (r *MongoReminderRepository) UpdateNextTrigger(reminderID int64, nextTrigger time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.col.UpdateOne(ctx, map[string]any{"id": reminderID}, map[string]any{"$set": map[string]any{"nextTrigger": nextTrigger}})
	return err
}

package persistent

import (
	"context"
	"fmt"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MongoPremiumUsageRepository implements PremiumUsageRepository using MongoDB
type MongoPremiumUsageRepository struct {
	collection *mongo.Collection
}

// NewMongoPremiumUsageRepository creates a new MongoDB premium usage repository
func NewMongoPremiumUsageRepository(connectionString string, databaseName string) (repositories.PremiumUsageRepository, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	client, err := mongo.Connect(options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	collection := client.Database(databaseName).Collection("premium_usage")

	return &MongoPremiumUsageRepository{
		collection: collection,
	}, nil
}

// GetUserUsage retrieves premium usage for a specific user
func (r *MongoPremiumUsageRepository) GetUserUsage(userID int64) (*entities.PremiumUsage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var usage entities.PremiumUsage
	err := r.collection.FindOne(ctx, bson.M{"userId": userID}).Decode(&usage)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("premium usage not found for user %d", userID)
		}
		return nil, fmt.Errorf("failed to get premium usage: %w", err)
	}

	return &usage, nil
}

// CreateUserUsage creates a new premium usage record for a user
func (r *MongoPremiumUsageRepository) CreateUserUsage(usage *entities.PremiumUsage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usage.CreatedAt = time.Now()
	usage.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, usage)
	if err != nil {
		return fmt.Errorf("failed to create premium usage: %w", err)
	}

	return nil
}

// UpdateUserUsage updates an existing premium usage record
func (r *MongoPremiumUsageRepository) UpdateUserUsage(usage *entities.PremiumUsage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usage.UpdatedAt = time.Now()

	filter := bson.M{"userId": usage.UserID}
	update := bson.M{"$set": usage}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update premium usage: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("premium usage not found for user %d", usage.UserID)
	}

	return nil
}

// GetOrCreateUserUsage gets existing usage or creates a new one
func (r *MongoPremiumUsageRepository) GetOrCreateUserUsage(userID int64) (*entities.PremiumUsage, error) {
	// First try to get existing usage
	if usage, err := r.GetUserUsage(userID); err == nil {
		return usage, nil
	}

	// Create new usage if not found
	newUsage := entities.NewPremiumUsage(userID)
	if err := r.CreateUserUsage(newUsage); err != nil {
		return nil, fmt.Errorf("failed to create premium usage for user %d: %w", userID, err)
	}

	return r.GetUserUsage(userID)
}

// GetAllUsage retrieves all premium usage records
func (r *MongoPremiumUsageRepository) GetAllUsage() ([]entities.PremiumUsage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all premium usage: %w", err)
	}
	defer cursor.Close(ctx)

	var usages []entities.PremiumUsage
	if err := cursor.All(ctx, &usages); err != nil {
		return nil, fmt.Errorf("failed to decode premium usage: %w", err)
	}

	return usages, nil
}

// DeleteUserUsage deletes premium usage record for a user
func (r *MongoPremiumUsageRepository) DeleteUserUsage(userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.M{"userId": userID})
	if err != nil {
		return fmt.Errorf("failed to delete premium usage: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("premium usage not found for user %d", userID)
	}

	return nil
}

// GetUsageByPremiumStatus gets usage records filtered by premium status
func (r *MongoPremiumUsageRepository) GetUsageByPremiumStatus(status entities.PremiumStatus) ([]entities.PremiumUsage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"premiumStatus": status})
	if err != nil {
		return nil, fmt.Errorf("failed to get premium usage by status: %w", err)
	}
	defer cursor.Close(ctx)

	var usages []entities.PremiumUsage
	if err := cursor.All(ctx, &usages); err != nil {
		return nil, fmt.Errorf("failed to decode premium usage: %w", err)
	}

	return usages, nil
}

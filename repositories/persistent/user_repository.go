package persistent

import (
	"context"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/errors"
	"github.com/ivanenkomaksym/remindme_bot/domain/repositories"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoUserRepository struct {
	client   *mongo.Client
	database string
	usersCol *mongo.Collection
}

func NewMongoUserRepository(connectionString string, database string) (repositories.UserRepository, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	client, err := mongo.Connect(options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPI))
	if err != nil {
		return nil, err
	}
	db := client.Database(database)
	return &MongoUserRepository{
		client:   client,
		database: database,
		usersCol: db.Collection("users"),
	}, nil
}

func (r *MongoUserRepository) GetUsers() ([]*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := r.usersCol.Find(ctx, struct{}{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var users []*entities.User
	for cur.Next(ctx) {
		var u entities.User
		if err := cur.Decode(&u); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, cur.Err()
}

func (r *MongoUserRepository) GetUser(userID int64) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var u entities.User
	err := r.usersCol.FindOne(ctx, map[string]any{"id": userID}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *MongoUserRepository) GetOrCreateUser(userID int64, userName, firstName, lastName, language string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	now := time.Now()

	var u entities.User
	err := r.usersCol.FindOne(ctx, map[string]any{"id": userID}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		u := entities.User{ID: userID, UserName: userName, FirstName: firstName, LastName: lastName, Language: language, CreatedAt: now}
		if _, err := r.usersCol.InsertOne(ctx, &u); err != nil {
			return nil, err
		}
		return &u, nil
	}

	return &u, nil
}

func (r *MongoUserRepository) UpdateUserLanguage(userID int64, language string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.usersCol.UpdateOne(ctx, map[string]any{"id": userID}, map[string]any{"$set": map[string]any{"language": language, "updatedAt": time.Now()}})
	return err
}

func (r *MongoUserRepository) UpdateUserInfo(userID int64, userName, firstName, lastName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.usersCol.UpdateOne(ctx, map[string]any{"id": userID}, map[string]any{"$set": map[string]any{"userName": userName, "firstName": firstName, "lastName": lastName, "updatedAt": time.Now()}})
	return err
}

func (r *MongoUserRepository) UpdateLocation(userID int64, location string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.usersCol.UpdateOne(ctx, map[string]any{"id": userID}, map[string]any{"$set": map[string]any{"location": location, "updatedAt": time.Now()}})
	return err
}

func (r *MongoUserRepository) CreateUser(userID int64, userName, firstName, lastName, language string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if user already exists
	var existingUser entities.User
	err := r.usersCol.FindOne(ctx, map[string]any{"id": userID}).Decode(&existingUser)
	if err != mongo.ErrNoDocuments {
		if err == nil {
			return nil, errors.ErrUserExists
		}
		return nil, err
	}

	user := entities.NewUser(userID, userName, firstName, lastName, language)

	if _, err := r.usersCol.InsertOne(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *MongoUserRepository) DeleteUser(userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.usersCol.DeleteOne(ctx, map[string]any{"id": userID})
	return err
}

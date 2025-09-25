package persistent

import (
	"context"
	"time"

	"github.com/ivanenkomaksym/remindme_bot/domain/entities"
	"github.com/ivanenkomaksym/remindme_bot/domain/repositories"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoUserRepository struct {
	client     *mongo.Client
	database   string
	usersCol   *mongo.Collection
	selectsCol *mongo.Collection
}

func NewMongoUserRepository(connectionString string, database string) (repositories.UserRepository, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	client, err := mongo.Connect(options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPI))
	if err != nil {
		return nil, err
	}
	db := client.Database(database)
	return &MongoUserRepository{
		client:     client,
		database:   database,
		usersCol:   db.Collection("users"),
		selectsCol: db.Collection("user_selections"),
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

func (r *MongoUserRepository) CreateOrUpdateUser(userID int64, userName, firstName, lastName, language string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	now := time.Now()
	u := entities.User{ID: userID, UserName: userName, FirstName: firstName, LastName: lastName, Language: language, UpdatedAt: now}
	// Try update
	res, err := r.usersCol.UpdateOne(ctx, map[string]any{"id": userID}, map[string]any{"$set": u})
	if err != nil {
		return nil, err
	}
	if res.MatchedCount == 0 {
		u.CreatedAt = now
		if _, err := r.usersCol.InsertOne(ctx, &u); err != nil {
			return nil, err
		}
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

func (r *MongoUserRepository) GetUserSelection(userID int64) (*entities.UserSelection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var s entities.UserSelection
	err := r.selectsCol.FindOne(ctx, map[string]any{"userId": userID}).Decode(&s)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *MongoUserRepository) SetUserSelection(userID int64, selection *entities.UserSelection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Attach user id to selection for storage
	// Ensure CreatedAt/UpdatedAt semantics are preserved within selection struct itself
	res, err := r.selectsCol.UpdateOne(ctx, map[string]any{"userId": userID}, map[string]any{"$set": selection})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		// Insert new selection document
		if _, err := r.selectsCol.InsertOne(ctx, map[string]any{"userId": userID, "selection": selection}); err != nil {
			return err
		}
	}
	return nil
}

func (r *MongoUserRepository) UpdateUserSelection(userID int64, selection *entities.UserSelection) error {
	return r.SetUserSelection(userID, selection)
}

func (r *MongoUserRepository) ClearUserSelection(userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.selectsCol.DeleteOne(ctx, map[string]any{"userId": userID})
	return err
}

func (r *MongoUserRepository) GetUserWithSelection(userID int64) (*entities.User, *entities.UserSelection, error) {
	u, err := r.GetUser(userID)
	if err != nil {
		return nil, nil, err
	}
	s, err := r.GetUserSelection(userID)
	if err != nil {
		return nil, nil, err
	}
	return u, s, nil
}

func (r *MongoUserRepository) CreateOrUpdateUserWithSelection(userID int64, userName, firstName, lastName, language string) (*entities.User, *entities.UserSelection, error) {
	u, err := r.CreateOrUpdateUser(userID, userName, firstName, lastName, language)
	if err != nil {
		return nil, nil, err
	}
	s, err := r.GetUserSelection(userID)
	if err != nil {
		return nil, nil, err
	}
	if s == nil {
		s = entities.NewUserSelection()
		if err := r.SetUserSelection(userID, s); err != nil {
			return nil, nil, err
		}
	}
	return u, s, nil
}

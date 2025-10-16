package store

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/ywl0806/yuno_kiroku/internal/api/services/user/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore struct {
	collection *mongo.Collection
}

func NewUserStore(db *mongo.Database) *UserStore {
	return &UserStore{
		collection: db.Collection("users"),
	}
}

func (s *UserStore) FindUsers(ctx context.Context) ([]models.User, error) {

	users, err := s.collection.Find(ctx, bson.D{})

	if err != nil {
		fmt.Println("error: ", err)
		return nil, err
	}

	var result []models.User
	fmt.Println("users: ", users)
	defer users.Close(ctx)
	if err := users.All(ctx, &result); err != nil {
		fmt.Println("error: ", err)
		return nil, err
	}
	fmt.Println("result: ", result)
	return result, nil
}

func (s *UserStore) FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, echo.NewHTTPError(404, "User not found")
		}
		return models.User{}, err
	}
	return user, nil
}

func (s *UserStore) CreateUser(ctx context.Context, user models.User) (models.User, error) {

	// Check if the user already exists
	existingUser := models.User{}
	if user.Email != nil {
		err := s.collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
		if err == nil {
			return models.User{}, echo.NewHTTPError(400, "User already exists")
		} else if err != mongo.ErrNoDocuments {
			return models.User{}, err
		}
	}

	// Insert the new user
	result, err := s.collection.InsertOne(ctx, user)

	if err != nil {
		return models.User{}, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return user, nil
}

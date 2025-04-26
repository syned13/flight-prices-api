package auth

import (
	"context"

	"github.com/syned13/flight-prices-api/internal/models"
	"github.com/syned13/flight-prices-api/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	usersCollection = "users"
)

type AuthRepository interface {
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(user models.User) error
}

type authRepository struct {
	mongo *mongo.Client
}

func NewAuthRepository(mongo *mongo.Client) AuthRepository {
	return &authRepository{
		mongo: mongo,
	}
}

func (r *authRepository) GetUserByUsername(username string) (*models.User, error) {
	collection := r.mongo.Database(config.GetConfig().Mongo().Database()).Collection(usersCollection)

	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	return &user, err
}

func (r *authRepository) CreateUser(user models.User) error {
	collection := r.mongo.Database(config.GetConfig().Mongo().Database()).Collection(usersCollection)
	_, err := collection.InsertOne(context.Background(), user)
	return err
}

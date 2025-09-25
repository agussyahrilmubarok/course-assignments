package store

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	USER_COLLECTION = "users"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password,omitempty" json:"-"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

//go:generate mockery --name=IUserMongoStore
type IUserMongoStore interface {
	Create(ctx context.Context, user *User) (*User, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type userMongoStore struct {
	collection *mongo.Collection
	log        zerolog.Logger
}

func NewUserMongoStore(db *mongo.Database, log zerolog.Logger) IUserMongoStore {
	return &userMongoStore{
		collection: db.Collection(USER_COLLECTION),
		log:        log,
	}
}

func (s *userMongoStore) Create(ctx context.Context, user *User) (*User, error) {
	now := time.Now()
	user.ID = primitive.NewObjectID()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to insert user")
		return nil, err
	}
	return user, nil
}

func (s *userMongoStore) FindByID(ctx context.Context, id primitive.ObjectID) (*User, error) {
	var user User
	err := s.collection.FindOne(ctx, primitive.M{"_id": id}).Decode(&user)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to find user by id")
		return nil, err
	}
	return &user, nil
}

func (s *userMongoStore) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := s.collection.FindOne(ctx, primitive.M{"email": email}).Decode(&user)
	if err != nil {
		s.log.Error().Err(err).Msg("failed to find user by email")
		return nil, err
	}
	return &user, nil
}

func (s *userMongoStore) Update(ctx context.Context, user *User) error {
	user.UpdatedAt = time.Now()
	_, err := s.collection.UpdateByID(ctx, user.ID, primitive.M{
		"$set": user,
	})
	if err != nil {
		s.log.Error().Err(err).Msg("failed to update user")
		return err
	}
	return nil
}

func (s *userMongoStore) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.collection.DeleteOne(ctx, primitive.M{"_id": id})
	if err != nil {
		s.log.Error().Err(err).Msg("failed to delete user")
		return err
	}
	return nil
}

package user

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IStore interface {
	CreateUser(ctx context.Context, user *User) error
	FindUserByID(ctx context.Context, userID string) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	ExistsUserEmailByIgnoreCase(ctx context.Context, email string) error
}

type store struct {
	db     *mongo.Database
	logger zerolog.Logger
}

func NewStore(db *mongo.Database, logger zerolog.Logger) IStore {
	return &store{
		db:     db,
		logger: logger,
	}
}

var (
	userColl = "users"
)

func (s *store) CreateUser(ctx context.Context, user *User) error {
	now := time.Now()

	user.CreatedAt = now
	user.UpdatedAt = now

	result, err := s.db.Collection(userColl).InsertOne(ctx, user)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create user in mongo db")
		return err
	}

	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		s.logger.Error().Err(err).Msg("failed to create user in mongo db")
		user.ID = id
	}

	s.logger.Info().
		Str("user_id", result.InsertedID.(primitive.ObjectID).Hex()).
		Msg("create user successfully")
	return nil
}

func (s *store) FindUserByID(ctx context.Context, userID string) (*User, error) {
	userIDObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("user_id", userID).
			Msg("failed to find user not found by id in mongo db")
		return nil, err
	}

	filter := bson.M{"_id": userIDObj}

	var user User
	err = s.db.Collection(userColl).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.logger.Error().
				Err(err).
				Str("user_id", userID).
				Msg("failed to find user not found by id in mongo db")
			return nil, errors.New("user not found")
		}

		s.logger.Error().
			Err(err).
			Str("user_id", userID).
			Msg("failed to find user by id in mongo db")
		return nil, err
	}

	s.logger.Info().
		Str("user_id", user.ID.Hex()).
		Msg("find user by id successfully")
	return &user, nil
}

func (s *store) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	filter := bson.M{"email": email}

	var user User
	err := s.db.Collection(userColl).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.logger.Error().
				Err(err).
				Str("user_email", email).
				Msg("failed to find user not found by email in mongo db")
			return nil, errors.New("user not found")
		}
		s.logger.Error().
			Err(err).
			Str("user_email", email).
			Msg("failed to find user by email in mongo db")
		return nil, err
	}

	s.logger.Info().
		Str("user_email", user.Email).
		Msg("find user by email successfully")
	return &user, nil
}

func (s *store) ExistsUserEmailByIgnoreCase(ctx context.Context, email string) error {
	filter := bson.M{
		"email": bson.M{
			"$regex": primitive.Regex{
				Pattern: "^" + email + "$",
				Options: "i",
			},
		},
	}

	var user struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	err := s.db.Collection(userColl).FindOne(ctx, filter).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}

		s.logger.Error().
			Err(err).
			Str("user_email", email).
			Msg("failed to check existence of email by case-insensitive in mongo db")
		return err
	}

	s.logger.Error().
		Str("user_email", email).
		Msg("email already exists in mongo db")
	return errors.New("email already exists (case-insensitive)")
}

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

//go:generate mockery --name=IStore
type IStore interface {
	FindAll(ctx context.Context) ([]User, error)
	FindByID(ctx context.Context, userID string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
	UpdateByID(ctx context.Context, userID string, user *User) error
	DeleteByID(ctx context.Context, userID string) error
}

type store struct {
	collection *mongo.Collection
	log        *zerolog.Logger
}

func NewStore(db *mongo.Database, log *zerolog.Logger) IStore {
	if db == nil {
		log.Fatal().Msg("Database connection is nil")
	}

	return &store{
		collection: db.Collection("users"),
		log:        log,
	}
}

func (s *store) FindAll(ctx context.Context) ([]User, error) {
	cur, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to fetch users from database")
		return nil, err
	}
	defer cur.Close(ctx)

	var users []User
	if err := cur.All(ctx, &users); err != nil {
		s.log.Error().Err(err).Msg("Failed to decode user list")
		return nil, err
	}

	s.log.Info().Int("count", len(users)).Msg("Successfully fetched all users")
	return users, nil
}

func (s *store) FindByID(ctx context.Context, userID string) (*User, error) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		s.log.Warn().Str("user_id", userID).Msg("Invalid user ID format")
		return nil, errors.New("invalid user ID format")
	}

	var user User
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.log.Warn().Str("user_id", userID).Msg("User not found")
			return nil, nil
		}
		s.log.Error().Err(err).Str("user_id", userID).Msg("failed to find user by ID")
		return nil, err
	}

	s.log.Info().Str("user_id", userID).Msg("User found by ID")
	return &user, nil
}

func (s *store) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.log.Warn().Str("email", email).Msg("User not found by email")
			return nil, nil
		}
		s.log.Error().Err(err).Str("email", email).Msg("Failed to find user by email")
		return nil, err
	}

	s.log.Info().Str("email", email).Msg("User found by email")
	return &user, nil
}

func (s *store) Create(ctx context.Context, user *User) error {
	now := time.Now()
	if user.ID == "" {
		objID := primitive.NewObjectID()
		user.ID = objID.Hex()
	}

	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			s.log.Error().Err(err).Msg("Email already exists")
			return err
		}

		s.log.Error().Err(err).Msg("Failed to create user")
		return err
	}

	s.log.Info().Str("user_id", user.ID).Msg("User created successfully")
	return nil
}

func (s *store) UpdateByID(ctx context.Context, userID string, user *User) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		s.log.Warn().Str("user_id", userID).Msg("Invalid user ID format during update")
		return errors.New("invalid user ID format")
	}

	user.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"email":      user.Email,
			"password":   user.Password,
			"updated_at": user.UpdatedAt,
		},
	}

	res, err := s.collection.UpdateByID(ctx, objID, update)
	if err != nil {
		s.log.Error().Err(err).Str("user_id", userID).Msg("Failed to update user")
		return err
	}

	if res.MatchedCount == 0 {
		s.log.Warn().Str("user_id", userID).Msg("No user found to update")
		return mongo.ErrNoDocuments
	}

	s.log.Info().Str("user_id", userID).Msg("User updated successfully")
	return nil
}

func (s *store) DeleteByID(ctx context.Context, userID string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		s.log.Warn().Str("user_id", userID).Msg("Invalid user ID format during delete")
		return errors.New("invalid user ID format")
	}

	res, err := s.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		s.log.Error().Err(err).Str("user_id", userID).Msg("Failed to delete user")
		return err
	}

	if res.DeletedCount == 0 {
		s.log.Warn().Str("user_id", userID).Msg("No user deleted (user not found)")
		return nil
	}

	s.log.Info().Str("user_id", userID).Msg("User deleted successfully")
	return nil
}

package store

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

//go:generate mockery --name=IUserStore
type IUserStore interface {
	Create(ctx context.Context, user *User) (*User, error)
	FindAll(ctx context.Context) ([]User, error)
	FindByID(ctx context.Context, userID string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	UpdateByID(ctx context.Context, userID string, updates *User) (*User, error)
	DeleteByID(ctx context.Context, userID string) error
	ExistsEmail(ctx context.Context, email string) (bool, error)
}

type userStore struct {
	collection *mongo.Collection
	log        *zap.Logger
}

func NewUserStore(db *mongo.Database, collectionName string, log *zap.Logger) IUserStore {
	return &userStore{
		collection: db.Collection(collectionName),
		log:        log,
	}
}

func (s *userStore) Create(ctx context.Context, user *User) (*User, error) {
	now := time.Now()
	if user.ID == "" {
		return nil, errors.New("user id is required")
	}
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			s.log.Warn("email already exists", zap.String("email", user.Email))
			return nil, errors.New("email already exists")
		}
		s.log.Error("failed to insert user", zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (s *userStore) FindAll(ctx context.Context) ([]User, error) {
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		s.log.Error("failed to find all users", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []User
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			s.log.Error("failed to decode user during find all", zap.Error(err))
			continue
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		s.log.Error("cursor error during find all", zap.Error(err))
		return nil, err
	}

	return users, nil
}

func (s *userStore) FindByID(ctx context.Context, userID string) (*User, error) {
	var user User
	err := s.collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		s.log.Error("failed to find user by id", zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (s *userStore) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		s.log.Error("failed to find user by email", zap.Error(err))
		return nil, err
	}
	return &user, nil
}

func (s *userStore) UpdateByID(ctx context.Context, userID string, updates *User) (*User, error) {
	updates.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name":       updates.Name,
			"email":      updates.Email,
			"password":   updates.Password,
			"role":       updates.Role,
			"updated_at": updates.UpdatedAt,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated User
	err := s.collection.FindOneAndUpdate(ctx, bson.M{"_id": userID}, update, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		s.log.Error("failed to update user", zap.Error(err))
		return nil, err
	}
	return &updated, nil
}

func (s *userStore) DeleteByID(ctx context.Context, userID string) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"_id": userID})
	if err != nil {
		s.log.Error("failed to delete user", zap.Error(err))
		return err
	}
	return nil
}

func (s *userStore) ExistsEmail(ctx context.Context, email string) (bool, error) {
	count, err := s.collection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		s.log.Error("failed to check if email exists", zap.Error(err))
		return false, err
	}
	return count > 0, nil
}

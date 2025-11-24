package account

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IRepository interface {
	FindByID(ctx context.Context, accountID primitive.ObjectID) (*Account, error)
	FindByEmail(ctx context.Context, accountEmail string) (*Account, error)
	CreateAccount(ctx context.Context, account *Account) (*Account, error)
	UpdateAccount(ctx context.Context, account *Account) (*Account, error)
	DeleteAccount(ctx context.Context, accountID primitive.ObjectID) error
}

type repository struct {
	db     *mongo.Database
	logger *logrus.Logger
}

func NewRepository(db *mongo.Database, logger *logrus.Logger) IRepository {
	logger.WithField("layer", "repository").Info("Account repository initialized")
	return &repository{
		db:     db,
		logger: logger,
	}
}

func (r *repository) FindByID(ctx context.Context, accountID primitive.ObjectID) (*Account, error) {
	var acc Account

	err := r.db.Collection("accounts").FindOne(ctx, bson.M{"_id": accountID}).Decode(&acc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.WithField("id", accountID.Hex()).Info("account not found")
			return nil, nil
		}
		r.logger.WithError(err).Error("failed to query account by id")
		return nil, err
	}

	r.logger.WithField("id", accountID.Hex()).Info("account found")
	return &acc, nil
}

func (r *repository) FindByEmail(ctx context.Context, accountEmail string) (*Account, error) {
	var acc Account
	err := r.db.Collection("accounts").FindOne(ctx, bson.M{"email": accountEmail}).Decode(&acc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.WithField("email", accountEmail).Info("account not found")
			return nil, nil
		}
		r.logger.WithError(err).Error("failed to query account by email")
		return nil, err
	}

	r.logger.WithField("email", accountEmail).Info("account found")
	return &acc, nil
}

func (r *repository) CreateAccount(ctx context.Context, account *Account) (*Account, error) {
	account.ID = primitive.NewObjectID()
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	_, err := r.db.Collection("accounts").InsertOne(ctx, account)
	if err != nil {
		r.logger.WithError(err).Error("failed to create account")
		return nil, err
	}

	r.logger.WithField("id", account.ID.Hex()).Info("account created")
	return account, nil
}

func (r *repository) UpdateAccount(ctx context.Context, account *Account) (*Account, error) {
	account.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name":       account.Name,
			"email":      account.Email,
			"password":   account.Password,
			"updated_at": account.UpdatedAt,
		},
	}

	_, err := r.db.Collection("accounts").UpdateByID(ctx, account.ID, update)
	if err != nil {
		r.logger.WithError(err).Error("failed to update account")
		return nil, err
	}

	r.logger.WithField("id", account.ID.Hex()).Info("account updated")
	return account, nil
}

func (r *repository) DeleteAccount(ctx context.Context, accountID primitive.ObjectID) error {
	_, err := r.db.Collection("accounts").DeleteOne(ctx, bson.M{"_id": accountID})
	if err != nil {
		r.logger.WithError(err).Error("failed to delete account")
		return err
	}

	r.logger.WithField("id", accountID.Hex()).Info("account deleted")
	return nil
}

package repos

import (
	"context"

	"example.com.backend/internal/domain"
	"example.com.backend/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IUserRepository interface {
	FindAll(ctx context.Context) ([]domain.User, error)
	FindAllByRole(ctx context.Context, role domain.UserRole) ([]domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
	DeleteByID(ctx context.Context, id string) error
	ExistsByEmailIgnoreCase(ctx context.Context, email string) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindAll(ctx context.Context) ([]domain.User, error) {
	log := logger.GetLoggerFromContext(ctx)

	var users []domain.User

	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		log.Error("failed fetching all users", zap.Error(err))
		return nil, err
	}

	log.Info("successfully fetched all users", zap.Int("count", len(users)))
	return users, nil
}

func (r *userRepository) FindAllByRole(ctx context.Context, role domain.UserRole) ([]domain.User, error) {
	log := logger.GetLoggerFromContext(ctx)

	var users []domain.User
	if err := r.db.WithContext(ctx).Where("role = ?", role).Find(&users).Error; err != nil {
		log.Error("failed fetching users by role", zap.String("user_role", string(role)), zap.Error(err))
		return nil, err
	}

	log.Info("successfully fetched users by role",
		zap.String("role", string(role)),
		zap.Int("count", len(users)),
	)

	return users, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	log := logger.GetLoggerFromContext(ctx)

	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		log.Error("failed fetching user by id", zap.String("user_id", id), zap.Error(err))
		return nil, err
	}

	log.Info("successfully fetched user by id", zap.String("user_id", id))
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	log := logger.GetLoggerFromContext(ctx)

	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		log.Error("failed fetching user by email", zap.String("user_email", email), zap.Error(err))
		return nil, err
	}

	log.Info("successfully fetched user by email", zap.String("user_email", email))
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	log := logger.GetLoggerFromContext(ctx)

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		log.Error("failed creating user", zap.String("user_email", user.Email), zap.Error(err))
		return nil, err
	}

	log.Info("successfully created user", zap.String("user_id", user.ID), zap.String("user_email", user.Email))
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	log := logger.GetLoggerFromContext(ctx)

	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		log.Error("failed updating user", zap.String("user_id", user.ID), zap.Error(err))
		return nil, err
	}

	log.Info("successfully updated user", zap.String("user_id", user.ID))
	return user, nil
}

func (r *userRepository) DeleteByID(ctx context.Context, id string) error {
	log := logger.GetLoggerFromContext(ctx)

	if err := r.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", id).Error; err != nil {
		log.Error("failed deleting user", zap.String("user_id", id), zap.Error(err))
		return err
	}

	log.Info("successfully deleted user", zap.String("user_id", id))
	return nil
}

func (r *userRepository) ExistsByEmailIgnoreCase(ctx context.Context, email string) (bool, error) {
	log := logger.GetLoggerFromContext(ctx)

	var count int64
	if err := r.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("LOWER(email) = LOWER(?)", email).
		Count(&count).Error; err != nil {
		log.Error("failed checking existing email", zap.String("email", email), zap.Error(err))
		return false, err
	}

	return count > 0, nil
}

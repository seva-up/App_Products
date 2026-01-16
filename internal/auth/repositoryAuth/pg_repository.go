package repositoryAuth

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seva-up/App_Products/internal/auth"
	"github.com/seva-up/App_Products/internal/models"
)

type authRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) auth.UserRepository {
	return &authRepository{db: db}
}

func (s *authRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	err := s.db.QueryRow(ctx, queryCreate, user.FirstName, user.LastName, user.Password, user.Email, user.Role).Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (s *authRepository) FindById(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(ctx, queryFindById, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Password, &user.Email, &user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to GetById user: %w", err)
	}
	return &user, nil
}

func (s *authRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	err := s.db.QueryRow(ctx, queryUpdate, user.FirstName, user.LastName, user.Password, user.Email, user.Role, user.ID).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Password, &user.Email, &user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user with if %d not found for update", user.ID)
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return user, nil
}

func (s *authRepository) Delete(ctx context.Context, id int) error {
	result, err := s.db.Exec(ctx, queryDelete, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("ERR NOT FOUND в репозитории бд удаление пользователя пользователя")
	}
	return nil
}

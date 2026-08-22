package service

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"go-waste-routes/internal/domain"
	"go-waste-routes/internal/repository/memory"
)

type UserService struct {
	store *memory.Store[domain.User]
}

func NewUserService(store *memory.Store[domain.User]) *UserService {
	return &UserService{store: store}
}

func (s *UserService) Create(ctx context.Context, user domain.User, rawPassword string, permissions []domain.Permission) (domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), 12)
	if err != nil {
		return domain.User{}, err
	}
	user.PasswordHash = string(hash)
	user.Status = domain.UserStatusEnabled
	user.Permissions = permissions
	return s.store.Create(ctx, user)
}

func (s *UserService) Update(ctx context.Context, id int64, user domain.User) (domain.User, error) {
	return s.store.Update(ctx, id, user)
}

func (s *UserService) GetByID(ctx context.Context, id int64) (domain.User, bool, error) {
	return s.store.Get(ctx, id)
}

func (s *UserService) List(ctx context.Context, page, pageSize int) ([]domain.User, int64, error) {
	return s.store.List(ctx, page, pageSize)
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	return s.store.Delete(ctx, id)
}

func (s *UserService) FindByUsername(ctx context.Context, username string) (domain.User, bool, error) {
	for _, user := range s.store.All() {
		if strings.EqualFold(user.Username, username) {
			return user, true, nil
		}
	}
	return domain.User{}, false, nil
}

func (s *UserService) VerifyPassword(user domain.User, rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(rawPassword))
}

func (s *UserService) ChangePassword(ctx context.Context, id int64, oldPassword, newPassword string) error {
	user, ok, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("user not found")
	}
	if err := s.VerifyPassword(user, oldPassword); err != nil {
		return err
	}
	return s.setPassword(ctx, id, user, newPassword)
}

func (s *UserService) ResetPassword(ctx context.Context, id int64, newPassword string) error {
	user, ok, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("user not found")
	}
	return s.setPassword(ctx, id, user, newPassword)
}

func (s *UserService) UpdateStatus(ctx context.Context, id int64, status string) error {
	user, ok, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("user not found")
	}
	user.Status = status
	_, err = s.store.Update(ctx, id, user)
	return err
}

func (s *UserService) setPassword(ctx context.Context, id int64, user domain.User, newPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	_, err = s.store.Update(ctx, id, user)
	return err
}

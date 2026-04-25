package users

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

// 🔹 REGISTER
func (s *Service) Register(email, password string) error {

	existingUser, _ := s.repo.FindByEmail(email)
	if existingUser != nil {
		return errors.New("email already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return err
	}

	user := &User{
		ID:       uuid.New(),
		Email:    email,
		Password: string(hashed),
		Role:     "client", // temporary (we’ll fix later)
	}

	return s.repo.Create(user)
}

// 🔹 LOGIN (this is what we were adding)
func (s *Service) Login(email, password string) (*User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

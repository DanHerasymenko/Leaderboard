package auth

import (
	"Leaderboard/internal/constants"
	"context"
	"errors"
	"fmt"
	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var (
	ErrUserAlreadyExists = fmt.Errorf("user with nickname %s already exists")
	ErrUserNotFound      = fmt.Errorf("user not found")
)

type User struct {
	ID           string    `json:"id" bson:"id"`
	Nickname     string    `json:"nickname" bson:"nickname"`
	PasswordHash string    `json:"password_hash" bson:"password_hash"`
	Role         string    `json:"role" bson:"role"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
}

type UserSearchParameters struct {
	ID       *string
	Nickname *string
}

func (s *Service) CreateUser(ctx context.Context, nickname, password string) (*User, error) {

	now := time.Now()

	user := &User{
		ID:           ulid.MustNew(ulid.Timestamp(now), ulid.DefaultEntropy()).String(),
		Nickname:     nickname,
		PasswordHash: password,
		Role:         constants.UserRole,
		CreatedAt:    now,
	}

	_, err := s.usersColl.InsertOne(ctx, user)
	if mongo.IsDuplicateKeyError(err) {
		return nil, ErrUserAlreadyExists
	} else if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return user, nil

}

func (s *Service) GetUserByParam(ctx context.Context, param *UserSearchParameters) (*User, error) {

	filter := bson.M{}

	switch {
	case param.ID != nil:
		filter = bson.M{"id": *param.ID}
	case param.Nickname != nil:
		filter = bson.M{"nickname": *param.Nickname}
	default:
		return nil, errors.New("ID or Nickname must be provided")
	}

	res := s.usersColl.FindOne(ctx, filter)
	if errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return nil, ErrUserNotFound
	} else if res.Err() != nil {
		return nil, fmt.Errorf("failed to find one: %w", res.Err())
	}

	u := &User{}

	err := res.Decode(u)
	if err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	return u, nil

}

func (s *Service) GetUserRole(ctx context.Context, userId *string) (string, error) {
	user, err := s.GetUserByParam(ctx, &UserSearchParameters{ID: userId})
	if err != nil {
		return "", err
	}
	return user.Role, nil

}

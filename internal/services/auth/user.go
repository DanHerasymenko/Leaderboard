package auth

import (
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
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
}

func (s *Service) CreateUser(ctx context.Context, nickname, password string) (*User, error) {

	now := time.Now()

	user := &User{
		ID:           ulid.MustNew(ulid.Timestamp(now), ulid.DefaultEntropy()).String(),
		Nickname:     nickname,
		PasswordHash: password,
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

func (s *Service) GetUserByName(ctx context.Context, nickname string) (*User, error) {

	filter := bson.M{"nickname": nickname}

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

package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/AdamWalker95/orders-api/model"
)

type RedisRepo struct {
	Client *redis.Client
}

func (r *RedisRepo) Insert(ctx context.Context, newUser model.LoginDetails) error {
	data, err := json.Marshal(newUser)
	if err != nil {
		return fmt.Errorf("failed to encode user details: %w", err)
	}

	txn := r.Client.TxPipeline()

	res := txn.SetNX(ctx, newUser.Email, string(data), 0)
	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to set: %w", err)
	}

	if err := txn.SAdd(ctx, "users", newUser.Email).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to add to orders set: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

var ErrNotExist = errors.New("user does not exist")

func (r *RedisRepo) FindByID(ctx context.Context, user string) (model.LoginDetails, error) {

	value, err := r.Client.Get(ctx, user).Result()
	if errors.Is(err, redis.Nil) {
		return model.LoginDetails{}, ErrNotExist
	} else if err != nil {
		return model.LoginDetails{}, fmt.Errorf("get order: %w", err)
	}

	var foundUser model.LoginDetails
	err = json.Unmarshal([]byte(value), &foundUser)
	if err != nil {
		return model.LoginDetails{}, fmt.Errorf("failed to decode order json: %w", err)
	}

	return foundUser, nil
}

func (r *RedisRepo) DeleteByID(ctx context.Context, user string) error {

	txn := r.Client.TxPipeline()

	err := txn.Del(ctx, user).Err()
	if errors.Is(err, redis.Nil) {
		txn.Discard()
		return ErrNotExist
	} else if err != nil {
		txn.Discard()
		return fmt.Errorf("get user: %w", err)
	}

	if err := txn.SRem(ctx, "users", user).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to remove from users: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

func (r *RedisRepo) Update(ctx context.Context, user model.LoginDetails) error {
	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user details: %w", err)
	}

	err = r.Client.SetXX(ctx, user.Email, string(data), 0).Err()
	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

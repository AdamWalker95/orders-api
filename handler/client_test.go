package handler_test

import (
	"testing"

	"github.com/AdamWalker95/orders-api/handler"
	"github.com/AdamWalker95/orders-api/repository/client"
	"github.com/stretchr/testify/assert"
)

func TestCPasswordValidation(t *testing.T) {

	userHandler := &handler.Client{
		Repo: &client.SqlRepo{
			Client: nil,
		},
	}

	t.Run("Attempt to use valid password", func(t *testing.T) {
		print("\n")
		t.Log("Attempt to use valid password")

		//Returns nothing as function only returns a valid if password doesn't pass validation
		assert.Equal(t, "", userHandler.ValidatePassword("password123"))
	})
	t.Run("Attempt to use too short password", func(t *testing.T) {
		print("\n")
		t.Log("Attempt to add user with a password that's too short")

		expectedError := "Error: Password is either less than 8 characters or doesn't contain a number"

		assert.Equal(t, expectedError, userHandler.ValidatePassword("pass1"))
	})
	t.Run("Attempt to use password with no numerical values", func(t *testing.T) {
		print("\n")
		t.Log("Attempt to use password with no numerical values")

		expectedError := "Error: Password is either less than 8 characters or doesn't contain a number"

		assert.Equal(t, expectedError, userHandler.ValidatePassword("password"))
	})
}

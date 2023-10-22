package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/AdamWalker95/orders-api/repository/client"
	"github.com/AdamWalker95/orders-api/repository/order"
)

type Login struct {
	OrdRedisRepo *order.RedisRepo
	OrdSqlRepo   *order.SqlRepo
	UsrSqlRepo   *client.SqlRepo
}

func (h *Login) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fmt.Println("failed: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.UsrSqlRepo.FindByID(body.Email)
	if errors.Is(err, client.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("Error occurred when trying to find user: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user.Password != body.Password {
		fmt.Println("User's password is incorrect")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.moveOrdersFromSqlToRedis(r.Context(), user.CustomerID)

	if err != nil {
		fmt.Println("Error with retrieving orders: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h *Login) moveOrdersFromSqlToRedis(ctx context.Context, customerID int) error {

	orders, err := h.OrdSqlRepo.FindByID(customerID)
	if err != nil {
		return err
	}

	err = h.OrdRedisRepo.InsertMulti(ctx, orders)

	if err != nil {
		return err
	}

	return nil
}

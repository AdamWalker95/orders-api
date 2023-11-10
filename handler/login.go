package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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

func (h *Login) Logout(w http.ResponseWriter, r *http.Request) {
	cursorStr := r.URL.Query().Get("cursor")
	if cursorStr == "" {
		cursorStr = "0"
	}

	const decimal = 10
	const bitSize = 64
	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	const size = 50
	err = h.OrdRedisRepo.RemoveMulti(r.Context(), order.FindAllPage{
		Offset: cursor,
		Size:   size,
	})
	if err != nil {
		criticalErr := fmt.Sprint("failed to Remove all orders from system when logging off: ", err)
		panic(criticalErr)
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

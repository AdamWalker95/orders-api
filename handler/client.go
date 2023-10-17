package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/AdamWalker95/orders-api/model"
	"github.com/AdamWalker95/orders-api/repository/client"
	"github.com/go-chi/chi/v5"
)

type Client struct {
	Repo *client.SqlRepo
}

func (h *Client) ValidatePassword(password string) string {
	if len(password) < 8 || regexp.MustCompile(`\d`).MatchString(password) == false {
		errMsg := "Error: Password is either less than 8 characters or doesn't contain a number"
		fmt.Println(errMsg)
		return errMsg
	}

	return ""
}
func (h *Client) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fmt.Println("failed: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userHTML := chi.URLParam(r, "id")

	// Checks to see if email is already in use
	_, emailInCurrentRecords := h.Repo.FindByID(userHTML)
	if emailInCurrentRecords == nil {
		fmt.Println("Record Already exists for email address")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.ValidatePassword(body.Password); err != "" {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("\n" + string(err) + "\n"))
		return
	}

	newUser := model.LoginDetails{
		Email:    body.Email,
		Password: body.Password,
	}

	err := h.Repo.Insert(r.Context(), newUser)
	if err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(newUser)
	if err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}

func (h *Client) GetByID(w http.ResponseWriter, r *http.Request) {
	userHTML := chi.URLParam(r, "id")

	user, err := h.Repo.FindByID(userHTML)
	if errors.Is(err, client.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// This is only used for updating password details for now
func (h *Client) UpdateByID(w http.ResponseWriter, r *http.Request) {
	var body struct {
		NewPassword string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userHTML := chi.URLParam(r, "id")

	user, err := h.Repo.FindByID(userHTML)
	if errors.Is(err, client.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := h.ValidatePassword(body.NewPassword); err != "" {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("\n" + string(err) + "\n"))
		return
	}

	user.Password = body.NewPassword

	err = h.Repo.Update(r.Context(), user)
	if err != nil {
		fmt.Println("failed to insert:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Client) DeleteByID(w http.ResponseWriter, r *http.Request) {
	userHTML := chi.URLParam(r, "id")

	err := h.Repo.DeleteByID(r.Context(), userHTML)
	if errors.Is(err, client.ErrNotExist) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		fmt.Println("failed to find user:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

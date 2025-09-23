package service

import (
	cmn "Codex-Backend/api/common"
	db "Codex-Backend/api/internal/database"
	"Codex-Backend/api/internal/domain"
	"context"
	"errors"
	"net/http"
)

func LoginUser(credentials domain.Credentials, ctx context.Context) (domain.User, error) {
	client, err := db.GetClient(ctx)
	if err != nil {
		return domain.User{}, err
	}

	user, err := client.GetUserByEmail(credentials.Email, ctx)
	if err != nil {
		return domain.User{}, err
	}

	err = cmn.VerifyPassword(user.Password, credentials.Password)
	if err != nil {
		return domain.User{}, &cmn.Error{Err: errors.New("Login Service Error - Invalid password"), Status: http.StatusUnauthorized}
	}

	return user, nil
}

func GetUserByID(id string, ctx context.Context) (domain.User, error) {
	client, err := db.GetClient(ctx)
	if err != nil {
		return domain.User{}, err
	}

	user, err := client.GetUserById(id, ctx)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func RegisterUser(newUser domain.NewUser, ctx context.Context) error {
	client, err := db.GetClient(ctx)
	if err != nil {
		return err
	}

	_, err = client.GetUserByEmail(newUser.Email, ctx)
	if e, ok := err.(*cmn.Error); ok {
		if e.StatusCode() != http.StatusNotFound {
			return &cmn.Error{Err: errors.New("Register Service Error - Getting User By Email: " + err.Error()), Status: http.StatusInternalServerError}
		}
	}

	hashedPassword, err := cmn.HashPassword(newUser.Password)
	if err != nil {
		return err
	}

	if err = client.CreateUser(domain.User{
		Username: newUser.Username,
		Password: string(hashedPassword),
		Email:    newUser.Email,
		Type:     "User",
	}, ctx); err != nil {
		return err
	}

	return nil
}

func LogoutUser(tokenString string) error {
	if tokenString == "" {
		return &cmn.Error{Err: errors.New("Logout Service Error - Token not found"), Status: http.StatusBadRequest}
	}

	return nil
}

func UpdateUser(updatedUser domain.User, ctx context.Context) error {
	client, err := db.GetClient(ctx)
	if err != nil {
		return err
	}

	if err = client.UpdateUser(updatedUser, ctx); err != nil {
		return err
	}

	return nil
}

func DeleteUser(id string, ctx context.Context) error {
	client, err := db.GetClient(ctx)
	if err != nil {
		return err
	}

	if err = client.DeleteUser(id, ctx); err != nil {
		return err
	}

	return nil
}

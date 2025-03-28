package users

import (
	"fmt"
	"go-image-processor/internal/db"
	"go-image-processor/internal/keys"
	"log/slog"

	"github.com/markbates/goth"
)

// CompleteAuthorization завершает процесс авторизации через OAuth.
// Здесь проверяется, записан ли пользователь в базу. Если не записан - записывается
func CompleteAuthorization(user *goth.User) (account *db.User, err error) {
	account, exists := FindByEmail(&user.Email)
	if !exists {
		acc, err := PersistAccount(user)
		if err != nil {
			return nil, err
		}

		return acc, nil
	}

	return account, nil
}

// PersistAccount сохраняет пользователя в базу. Принимает данные от goth.
func PersistAccount(u *goth.User) (*db.User, error) {
	keys, err := keys.NewKeysPair(u.Email)
	if err != nil {
		slog.Error("failed to generate keys pair. do user has email?", "error", err.Error())
	}

	user := &db.User{
		Name:          u.Name,
		Provider:      u.Provider,
		Email:         u.Email,
		AvatarURL:     &u.AvatarURL,
		NickName:      &u.NickName,
		PublicApiKey:  keys.PublicKey,
		PrivateApiKey: keys.PrivateKey,
	}

	dbr := db.Client.Create(user)
	if dbr.Error != nil {
		slog.Error("failed to create user record", "error", dbr.Error.Error())
		return nil, fmt.Errorf("failed to create user record: %s", dbr.Error.Error())
	}

	return user, nil
}

func FindByEmail(email *string) (user *db.User, exists bool) {
	accountResult := &db.User{}
	dbr := db.Client.Where("email = ?", email).First(&accountResult)
	if dbr.Error != nil {
		return nil, false
	}

	return accountResult, true
}

func FindById(userId string) (user *db.User, err error) {
	accountResult := &db.User{}
	dbr := db.Client.Where("id = ?", userId).First(&accountResult)
	if dbr.Error != nil {
		return nil, err
	}

	return accountResult, nil
}

package usecase

import "context"

func (u *UseCase) AddUser(ctx context.Context, userID, chatID int64, username string) error {
	const op = "usecase.addUser"

	err := u.postgres.AddUser(ctx, userID, chatID, username)
	if err != nil {
		u.log.Error("Failed to add user", "op", op, "error", err)
		return err
	}

	return nil
}

func (u *UseCase) RemoveUser(ctx context.Context, userID, chatID int64) error {
	const op = "usecase.removeUser"

	err := u.postgres.RemoveUser(ctx, userID, chatID)
	if err != nil {
		u.log.Error("Failed to remove user", "op", op, "error", err)
		return err
	}

	return nil
}

func (u *UseCase) CheckUserInChat(ctx context.Context, userName string, chatID int64) (bool, error) {
	const op = "usecase.CheckUserInChat"

	_, err := u.postgres.GetUserFromChat(ctx, chatID, userName)
	if err != nil {
		u.log.Error("error while checking user in chat", "op", op, "error", err)
		return false, err
	}

	return true, nil
}

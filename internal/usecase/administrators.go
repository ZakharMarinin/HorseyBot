package usecase

import "context"

func (u *UseCase) AddAdministrator(ctx context.Context, userID int64, username string) error {
	const op = "usecase.addAdministrator"

	err := u.postgres.AddAdministrator(ctx, userID, username)
	if err != nil {
		u.log.Error("Failed to add administrator", "op", op, "error", err)
		return err
	}

	return err
}

func (u *UseCase) RemoveAdministrator(ctx context.Context, userID int64) error {
	const op = "usecase.removeAdministrator"

	err := u.postgres.RemoveAdministrator(ctx, userID)
	if err != nil {
		u.log.Error("Failed to remove administrator", "op", op, "error", err)
		return err
	}

	return err
}

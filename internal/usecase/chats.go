package usecase

import (
	"context"
	domain "horsey/internal/domain/entity"
)

func (u *UseCase) AddNewChat(ctx context.Context, chatID int64, chatName string) error {
	const op = "usecase.AddNewChat"

	err := u.postgres.AddNewChat(ctx, chatID, chatName)
	if err != nil {
		u.log.Error("Could not add new chat", "op", op, "error", err)
		return err
	}

	return nil
}

func (u *UseCase) CheckChat(ctx context.Context, chatName string) (*domain.Chat, error) {
	const op = "usecase.GetChats"

	chat, err := u.postgres.GetChat(ctx, chatName)
	if err != nil {
		u.log.Error("Cannot get chat list: ", "op", op, "error", err)
		return nil, err
	}

	return chat, nil
}

func (u *UseCase) DeleteChat(ctx context.Context, chatID int64) error {
	const op = "usecase.DeleteChat"

	err := u.postgres.RemoveChat(ctx, chatID)
	if err != nil {
		u.log.Error("Could not delete chat", "op", op, "error", err)
		return err
	}

	return nil
}

func (u *UseCase) GetChats(ctx context.Context) ([]domain.Chat, error) {
	const op = "usecase.GetChats"

	chats, err := u.postgres.GetChats(ctx)
	if err != nil {
		u.log.Error("Cannot get chat list: ", "op", op, "error", err)
		return nil, err
	}

	return chats, nil
}

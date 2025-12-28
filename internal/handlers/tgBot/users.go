package tgBot

import (
	"context"
	domain "horsey/internal/domain/entity"

	"gopkg.in/telebot.v4"
)

func (b *TgBot) AddUser(ctx context.Context, c telebot.Context) error {
	err := b.useCase.AddUser(ctx, c.Sender().ID, c.Chat().ID, c.Sender().Username)
	if err != nil {
		return err
	}

	return nil
}

func (b *TgBot) UserJoined(ctx context.Context) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		newUsers := c.Message().UsersJoined

		for _, user := range newUsers {
			err := b.useCase.AddUser(ctx, user.ID, c.Chat().ID, user.Username)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func (b *TgBot) UserLeft(ctx context.Context) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		leftUser := c.Message().UserLeft

		if leftUser != nil {
			err := b.useCase.RemoveUser(ctx, leftUser.ID, c.Chat().ID)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func (b *TgBot) GetUserState(userID int64) *domain.TempUserState {
	if state, ok := tempUserState[userID]; ok {
		return state
	}

	newState := &domain.TempUserState{
		UserID: userID,
		State:  "",
	}
	tempUserState[userID] = newState
	return newState
}

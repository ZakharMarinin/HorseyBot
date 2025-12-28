package tgBot

import (
	"context"
	"fmt"
	domain "horsey/internal/domain/entity"
	"strconv"
	"strings"

	"gopkg.in/telebot.v4"
)

func (b *TgBot) AddSubscription(ctx context.Context) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		userState := b.GetUserState(c.Sender().ID)
		if userState.State == "" {
			return c.Send("Сперва нужно зарегистрироваться!")
		}
		menu := telebot.ReplyMarkup{ResizeKeyboard: true}

		imageBtn := menu.Data("Триггер слово", fmt.Sprintf("%d", domain.SendImage))
		pingBtn := menu.Data("Пинг", fmt.Sprintf("%d", domain.SendPing))

		menu.Inline(menu.Row(imageBtn, pingBtn))

		c.Send("Выдери действие, с которым ты бы хотел создать связь:", &menu)
		tempUserState[c.Sender().ID].State = domain.WaitingAction

		return nil
	}
}

func (b *TgBot) RemoveSubscription(ctx context.Context) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		userState := b.GetUserState(c.Sender().ID)

		if userState.State == "" {
			return c.Send("Сперва нужно зарегистрироваться!")
		}

		userState.Action = domain.DeleteSub
		userState.State = domain.WaitingChat

		b.SelectChat(ctx, c)
		return nil
	}
}

func (b *TgBot) GetSubs(ctx context.Context, c telebot.Context) error {
	subs, err := b.useCase.GetSubs(ctx, tempUserState[c.Sender().ID].ChatID)
	if err != nil {
		return c.Send("Не вижу связей в данном чате")
	}

	msg := "Теперь напиши мне номер связи для удаления:\n\n"

	for _, sub := range *subs {
		var actionName string
		switch sub.Feature {
		case domain.SendImage:
			actionName = "Триггер слово"
		case domain.SendPing:
			actionName = "Пинг"
		}
		msg += fmt.Sprintf("Номер связи: %d\nВ чате: %s\nТип: %s\nОжидаемый пользователь: %s\n\n", sub.ID, sub.ChatName, actionName, sub.Store.TrackedUser)
	}

	c.Send(msg)

	return nil
}

func (b *TgBot) GetSubsWithFilter(ctx context.Context, c telebot.Context, userFilter, userData string) error {
	subs, err := b.useCase.GetSubsWithFilter(ctx, tempUserState[c.Sender().ID].ChatID, userFilter, userData)
	if err != nil {
		return c.Send("Не вижу связей c подобным фильтром")
	}

	msg := ""

	for _, sub := range *subs {
		var actionName string
		switch sub.Feature {
		case domain.SendImage:
			actionName = "Триггер слово"
			msg += fmt.Sprintf("Номер связи: %d\nВ чате: %s\nТип: %s\nКлючевое слово: %s\nОжидаемый пользователь: %s\n\n", sub.ID, sub.ChatName, actionName, sub.Store.Keyword, sub.Store.TrackedUser)
		case domain.SendPing:
			actionName = "Пинг"
			msg += fmt.Sprintf("Номер связи: %d\nВ чате: %s\nТип: %s\nОжидаемый пользователь: %s\nПороговое значение: %d\n\n", sub.ID, sub.ChatName, actionName, sub.Store.TrackedUser, sub.Store.Threshold)
		}

	}

	return c.Send(msg)
}

func (b *TgBot) ShowSubsWithFilter(ctx context.Context) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		userState := b.GetUserState(c.Sender().ID)

		if userState.State == "" {
			return c.Send("Сперва нужно зарегистрироваться!")
		}

		userState.Action = domain.GetSubs
		userState.State = domain.WaitingChat

		b.SelectChat(ctx, c)
		return nil
	}
}

func (b *TgBot) HandleInlineButtonSelection(ctx context.Context) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		userState := b.GetUserState(c.Sender().ID)
		if userState.State == domain.WaitingAction {
			actionName := c.Data()
			defer c.Respond()

			userID := c.Sender().ID
			var err error

			actionNum, _ := strings.CutPrefix(actionName, "\f")

			tempUserState[userID].Action, err = strconv.Atoi(actionNum)
			if err != nil {
				return err
			}

			c.Edit("Отлично! Теперь выберите чат, в который нужно добавить связь: ")
			b.SelectChat(ctx, c)
		} else if userState.State == domain.WaitingFilter {
			filter := c.Data()
			defer c.Respond()

			filter, _ = strings.CutPrefix(filter, "\f")

			userState.Filter = filter
			userState.State = domain.WaitingFilterData

			if filter == "no-filter" {
				b.GetSubsWithFilter(ctx, c, userState.Filter, "")
			}

			c.Edit("Отлично! Теперь напиши, по какому значению нужно отфильтровать связи: ")
		}

		return nil
	}
}

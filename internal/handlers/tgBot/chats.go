package tgBot

import (
	"context"
	"fmt"
	domain "horsey/internal/domain/entity"

	"gopkg.in/telebot.v4"
)

func (b *TgBot) SelectChat(ctx context.Context, c telebot.Context) error {
	userID := c.Sender().ID
	foundInAny := false
	userState := b.GetUserState(c.Sender().ID)
	userState.State = domain.WaitingChat

	var availableChat []domain.Chat

	chats, err := b.useCase.GetChats(ctx)
	if err != nil {
		return c.Send("Похоже, я не состою ни в одном чате. Можете добавить меня в один, я выполню ваш запрос!")
	}

	for _, chat := range chats {
		member, err := b.Bot.ChatMemberOf(&telebot.Chat{ID: chat.ChatID}, &telebot.User{ID: userID})
		if err != nil {
			continue
		}

		if member.Role != telebot.Left && member.Role != telebot.Kicked {
			foundInAny = true
			availableChat = append(availableChat, chat)
		}
	}

	if foundInAny {
		if len(availableChat) == 1 {
			userState.ChatID = availableChat[0].ChatID
			userState.ChatName = availableChat[0].ChatName
			switch userState.Action {
			case domain.SendImage:
				userState.State = domain.WaitingImage
				c.Send("Мы состоим только в одном чате, так понимаю мы будем работать с ним.\nТеперь пришли мне изображение с которым хочешь создать связь.")
			case domain.SendPing:
				userState.State = domain.WaitingUser
				c.Send("Мы состоим только в одном чате, так понимаю мы будем работать с ним.\nТеперь пришли мне пользователя чата, с которым хочешь создать связь.")
			case domain.DeleteSub:
				userState.State = domain.WaitingDelete
				c.Send("Мы состоим только в одном чате, так понимаю, мы будем работать с ним.")
				b.GetSubs(ctx, c)
			case domain.GetSubs:
				menu := telebot.ReplyMarkup{ResizeKeyboard: true}

				typeBtn := menu.Data("Тип связи", "type-filter")
				userBtn := menu.Data("Пользователь", "user-filter")
				defaultBtn := menu.Data("Без фильтров", "no-filter")

				menu.Inline(menu.Row(typeBtn, userBtn), menu.Row(defaultBtn))

				userState.State = domain.WaitingFilter

				c.Send("Как бы ты хотел отфильтровать связи?", &menu)
			}
			return nil
		}
		message := "Отправь мне название чата, в котором мы будем работать:\n\n"
		for i, chat := range availableChat {
			message += fmt.Sprintf("%d: %s\n", i+1, chat.ChatName)
		}
		c.Send(message)
	}

	return nil
}

func (b *TgBot) ChatManagement(ctx context.Context) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		if c.ChatMember().NewChatMember.Role == telebot.Member {
			err := b.useCase.AddNewChat(ctx, c.Chat().ID, c.Chat().Title)
			if err != nil {
				return err
			}
		} else if c.ChatMember().NewChatMember.Role == telebot.Left {
			err := b.useCase.DeleteChat(ctx, c.Chat().ID)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func (b *TgBot) CheckUserInChat(ctx context.Context, c telebot.Context, userName string) (bool, error) {
	userState := b.GetUserState(c.Sender().ID)
	isIt, err := b.useCase.CheckUserInChat(ctx, userName, userState.ChatID)
	if err != nil {
		return isIt, c.Send("Не вижу такого пользователя в данной группе. Попробуй еще раз.")
	}

	return isIt, nil
}

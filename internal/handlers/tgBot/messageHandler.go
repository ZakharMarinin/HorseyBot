package tgBot

import (
	"context"
	domain "horsey/internal/domain/entity"
	"strconv"
	"time"

	"gopkg.in/telebot.v4"
)

func (b *TgBot) MessageHandler(ctx context.Context) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		if c.Chat().Type == telebot.ChatPrivate {
			userState := b.GetUserState(c.Sender().ID)
			switch userState.State {
			case "":
				c.Send("Сперва нужно зарегистрироваться!")
			case domain.WaitingChat:
				chatName := c.Message().Text

				chat, err := b.useCase.CheckChat(ctx, chatName)
				if err != nil {
					return c.Send("Кажется такого чата нет в списке, попробуйте еще раз.")
				}

				tempUserState[c.Sender().ID].ChatID = chat.ChatID
				tempUserState[c.Sender().ID].ChatName = chat.ChatName

				switch tempUserState[c.Sender().ID].Action {
				case domain.SendImage:
					tempUserState[c.Sender().ID].State = domain.WaitingImage
					c.Send("Теперь пришли мне изображение с которым хочешь создать связь.")
				case domain.SendPing:
					tempUserState[c.Sender().ID].State = domain.WaitingUser
					c.Send("Теперь пришли мне пользователя чата, с которым хочешь создать связь.")
				case domain.DeleteSub:
					tempUserState[c.Sender().ID].State = domain.WaitingDelete
					b.GetSubs(ctx, c)
				case domain.GetSubs:
					menu := telebot.ReplyMarkup{ResizeKeyboard: true}

					typeBtn := menu.Data("Тип связи", "type-filter")
					userBtn := menu.Data("Пользователь", "user-filter")
					defaultBtn := menu.Data("Без фильтров", "no-filter")

					menu.Inline(menu.Row(typeBtn, userBtn), menu.Row(defaultBtn))

					tempUserState[c.Sender().ID].State = domain.WaitingFilter

					c.Send("Как бы ты хотел отфильтровать связи?", &menu)
				}

				return nil
			case domain.WaitingImage:
				msg := c.Message()

				err := b.HandleMedia(c, msg)
				if err != nil {
					return err
				}

				return nil
			case domain.WaitingKeyword:
				msg := c.Message().Text

				tempUserState[c.Sender().ID].Store.Keyword = msg
				tempUserState[c.Sender().ID].State = domain.WaitingUser

				c.Send("Ключевое слово выбрано! Теперь вы можете выбрать пользователя, к которому привяжите данную связь.\n(Этот этап необязательный - напиши 'skip' чтобы пропустить)")

				return nil
			case domain.WaitingUser:
				userName := c.Message().Text

				if userState.Action != 2 {
					if userName == "skip" {
						tempUserState[c.Sender().ID].Store.TrackedUser = ""
						tempUserState[c.Sender().ID].State = domain.WaitingChance
						return c.Send("Супер! Осталось только указать с каким шансом будет выполняться данная связь.\nУкажите значение шанса 1 к X. Где X - ваше число.")
					}
				}

				isIt, err := b.CheckUserInChat(ctx, c, userName)
				if err != nil {
					return c.Send("Кажется, такого пользователя нет в группе. Попробуйте еще раз.")
				}

				if isIt {
					switch userState.Action {
					case 1:
						userState.State = domain.WaitingChance
						c.Send("Супер! Осталось только указать с каким шансом будет выполняться данная связь.\nУкажите значение шанса 1 к X. Где X - ваше число.")
					case 2:
						userState.State = domain.WaitingChance
						c.Send("Супер! Осталось только указать дни отсутствия человека\n(Укажите количество дней, при которых будет срабатывать пинг).")
					}
					userState.Store.TrackedUser = userName
				}

				return nil
			case domain.WaitingChance:
				msg := c.Message().Text

				num, err := strconv.Atoi(msg)
				if err != nil {
					return c.Send("Прошу, укажите целое, положительное число.")
				}

				if num+1 <= 0 {
					return c.Send("Прошу, укажите целое, положительное число.")
				}

				switch userState.Action {
				case 1:
					tempUserState[c.Sender().ID].Store.Chance = num + 1
				case 2:
					tempUserState[c.Sender().ID].Store.Threshold = num
					tempUserState[c.Sender().ID].Store.StartTime = time.Now().UTC().Format("2006-01-02 15:04:05")
				}

				tempUserState[c.Sender().ID].State = domain.WaitingCommand

				err = b.useCase.AddSub(ctx, tempUserState[c.Sender().ID])
				if err != nil {
					return c.Send("Возникла ошибка с добавлением связи.")
				}

				return c.Send("Готово! Связь создана!")
			case domain.WaitingDelete:
				msg := c.Message().Text

				num, err := strconv.Atoi(msg)
				if err != nil {
					return c.Send("Не вижу такого номера в списке. Попробуйте еще раз.")
				}

				b.useCase.RemoveSub(ctx, num)
				return c.Send("Готово! Связь была удалена.")
			case domain.WaitingFilterData:
				userData := c.Message().Text
				b.GetSubsWithFilter(ctx, c, tempUserState[c.Sender().ID].Filter, userData)
			}
		}

		if c.Chat().Type == telebot.ChatGroup {
			b.OgoMeter(ctx, c)
			b.SelectFeature(ctx, c)
			b.AddUser(ctx, c)
		}

		return nil
	}
}

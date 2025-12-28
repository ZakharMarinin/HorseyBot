package tgBot

import (
	"context"
	domain "horsey/internal/domain/entity"
	"io"
	"log/slog"

	"gopkg.in/telebot.v4"
)

type UseCase interface {
	AddAdministrator(ctx context.Context, userID int64, username string) error
	RemoveAdministrator(ctx context.Context, userID int64) error
	AddUser(ctx context.Context, userID, chatID int64, username string) error
	RemoveUser(ctx context.Context, userID, chatID int64) error
	AddNewChat(ctx context.Context, chatID int64, chatName string) error
	DeleteChat(ctx context.Context, chatID int64) error
	GetChats(ctx context.Context) ([]domain.Chat, error)
	CheckChat(ctx context.Context, chatName string) (*domain.Chat, error)
	AddSub(ctx context.Context, link *domain.TempUserState) error
	UpdateSub(ctx context.Context, link *domain.Subscription) error
	RemoveSub(ctx context.Context, subID int) error
	GetSubs(ctx context.Context, chatID int64) (*[]domain.Subscription, error)
	GetSubsWithFilter(ctx context.Context, chatID int64, userFilter, userData string) (*[]domain.Subscription, error)
	GetExpiredSubs(ctx context.Context) (*[]domain.Subscription, error)
	HandleMedia(file io.ReadCloser, mimeType string) (string, string, error)
	CheckUserInChat(ctx context.Context, userName string, chatID int64) (bool, error)
}

var (
	tempUserState = make(map[int64]*domain.TempUserState)
	ogoMeter      = make(map[int64]*domain.OgoMeter)
)

type TgBot struct {
	log     *slog.Logger
	Bot     *telebot.Bot
	useCase UseCase
}

func New(log *slog.Logger, bot *telebot.Bot, useCase UseCase) *TgBot {
	return &TgBot{
		log:     log,
		Bot:     bot,
		useCase: useCase,
	}
}

func (b *TgBot) Start(ctx context.Context) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		if c.Chat().Type != telebot.ChatPrivate {
			return nil
		}
		tempUserState[c.Sender().ID] = &domain.TempUserState{
			UserID: c.Sender().ID,
			ChatID: c.Chat().ID,
			Action: 0,
			Store: domain.Store{
				Threshold:   0,
				Chance:      0,
				Image:       "",
				ImageType:   "",
				Keyword:     "",
				TrackedUser: "",
			},
			State: domain.WaitingCommand,
		}

		menu := &telebot.ReplyMarkup{ResizeKeyboard: true}
		btnAddSub := menu.Text("Добавить связь")
		btnRemoveSub := menu.Text("Убрать связь")
		btnShowSubs := menu.Text("Показать связи")
		btnHelp := menu.Text("Помощь")

		menu.Reply(menu.Row(btnAddSub, btnRemoveSub), menu.Row(btnShowSubs, btnHelp))

		err := b.useCase.AddAdministrator(ctx, c.Sender().ID, c.Sender().Username)
		if err != nil {
			return c.Send("Вы уже зарегистрированы!", menu)
		}

		c.Send("Добро пожаловать!\nВы успешно зарегистрировались. Выберите команду:", menu)
		b.HelpMessage(c)

		return nil
	}
}

func (b *TgBot) HelpMessage(c telebot.Context) error {
	userState := b.GetUserState(c.Sender().ID)
	if userState.State == "" {
		return c.Send("Сперва нужно зарегистрироваться!")
	}
	if c.Chat().Type != telebot.ChatPrivate {
		return nil
	}
	var parse telebot.ParseMode = "MarkDown"

	c.Send("Список действующий команд: "+
		"\n\n*Создать связь* - привязывает медиафайлы с ключевыми словами для их вызова в чате."+
		"\n\n*Удалить связь* - удаляет созданные связи."+
		"\n\n*Все связи* - отображает набор связей созданные вами.", parse)

	return nil
}

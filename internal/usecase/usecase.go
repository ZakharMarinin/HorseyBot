package usecase

import (
	"context"
	domain "horsey/internal/domain/entity"
	"log/slog"
)

type Postgres interface {
	AddAdministrator(ctx context.Context, userID int64, username string) error
	RemoveAdministrator(ctx context.Context, userID int64) error
	AddUser(ctx context.Context, userID, chatID int64, username string) error
	RemoveUser(ctx context.Context, userID, chatID int64) error
	UpdateUser(ctx context.Context, userID, chatID int64) error
	GetUserFromChat(ctx context.Context, chatID int64, userName string) (*domain.User, error)
	UpdateSubscription(ctx context.Context, link *domain.Subscription) error
	GetSubscriptions(ctx context.Context, chatID int64) (*[]domain.Subscription, error)
	GetSubsWithFilter(ctx context.Context, chatID int64, userFilter, userData string) (*[]domain.Subscription, error)
	GetExpiredSubscriptions(ctx context.Context) (*[]domain.Subscription, error)
	AddSubscription(ctx context.Context, link *domain.TempUserState) error
	RemoveSubscription(ctx context.Context, subID int) error
	AddNewChat(ctx context.Context, chatID int64, chatName string) error
	RemoveChat(ctx context.Context, chatID int64) error
	GetChats(ctx context.Context) ([]domain.Chat, error)
	GetChat(ctx context.Context, chatName string) (*domain.Chat, error)
}

type UseCase struct {
	log      *slog.Logger
	postgres Postgres
}

func NewUseCase(log *slog.Logger, postgres Postgres) *UseCase {
	return &UseCase{
		log:      log,
		postgres: postgres,
	}
}

package usecase

import (
	"context"
	domain "horsey/internal/domain/entity"
	"strings"
)

func (u *UseCase) AddSub(ctx context.Context, link *domain.TempUserState) error {
	const op = "usecase.AddLink"

	if link.Store.TrackedUser != "" {
		link.Store.TrackedUser, _ = strings.CutPrefix(link.Store.TrackedUser, "@")
	}

	err := u.postgres.AddSubscription(ctx, link)
	if err != nil {
		u.log.Error("Cannot add new subscription", "op", op, "error", err)
		return err
	}

	return nil
}

func (u *UseCase) UpdateSub(ctx context.Context, link *domain.Subscription) error {
	const op = "usecase.UpdateSub"

	err := u.postgres.UpdateSubscription(ctx, link)
	if err != nil {
		u.log.Error("Cannot update subscription", "op", op, "error", err)
		return err
	}

	return nil
}

func (u *UseCase) GetSubs(ctx context.Context, chatID int64) (*[]domain.Subscription, error) {
	const op = "usecase.GetSubByChatID"

	subs, err := u.postgres.GetSubscriptions(ctx, chatID)
	if err != nil {
		u.log.Error("Cannot get subscriptions", "op", op, "error", err)
		return nil, err
	}

	return subs, nil
}

func (u *UseCase) GetSubsWithFilter(ctx context.Context, chatID int64, userFilter, userData string) (*[]domain.Subscription, error) {
	const op = "usecase.GetSubsWithFilter"

	switch userFilter {
	case "user-filter":
		userData, _ = strings.CutPrefix(userData, "@")
	case "type-filter":
		userData = strings.ToLower(userData)
		userData = strings.TrimSpace(userData)

		if userData == "триггер слово" {
			userData = "1"
		} else if userData == "пинг" {
			userData = "2"
		}
	}

	subs, err := u.postgres.GetSubsWithFilter(ctx, chatID, userFilter, userData)
	if err != nil {
		u.log.Error("Cannot get subscriptions", "op", op, "error", err)
		return nil, err
	}

	return subs, nil
}

func (u *UseCase) GetExpiredSubs(ctx context.Context) (*[]domain.Subscription, error) {
	const op = "usecase.GetExpiredSubs"

	subs, err := u.postgres.GetExpiredSubscriptions(ctx)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

func (u *UseCase) RemoveSub(ctx context.Context, subID int) error {
	const op = "usecase.RemoveSub"

	err := u.postgres.RemoveSubscription(ctx, subID)
	if err != nil {
		u.log.Error("Cannot remove subscription", "op", op, "error", err)
		return err
	}

	return nil
}

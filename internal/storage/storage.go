package storage

import (
	"context"
	"embed"
	"fmt"
	domain "horsey/internal/domain/entity"
	"strconv"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type PostgresStorage struct {
	DB *pgxpool.Pool
}

func New(ctx context.Context, storagePath string) (*PostgresStorage, error) {
	const op = "storage.postgresql.SQL.NEW"

	if err := runMigrations(storagePath); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	poolConfig, err := pgxpool.ParseConfig(storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &PostgresStorage{DB: db}, nil
}

func (p *PostgresStorage) AddAdministrator(ctx context.Context, userID int64, username string) error {
	const op = "storage.addAdministrator"

	query, args, err := sq.
		Insert("administrators").
		Columns("user_id", "username").
		Values(userID, username).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = p.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *PostgresStorage) RemoveAdministrator(ctx context.Context, userID int64) error {
	const op = "storage.removeAdministrator"

	query, args, err := sq.
		Delete("administrators").
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = p.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *PostgresStorage) AddUser(ctx context.Context, userID, chatID int64, username string) error {
	const op = "storage.addUser"

	query := `
        INSERT INTO users (user_id, username, chat_ids) 
        VALUES ($1, $2, ARRAY[$3]::bigint[])
        ON CONFLICT (user_id) DO UPDATE SET 
            username = EXCLUDED.username,
            chat_ids = ARRAY(SELECT DISTINCT e FROM unnest(array_append(users.chat_ids, $3)) AS e)
    `
	_, err := p.DB.Exec(ctx, query, userID, username, chatID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *PostgresStorage) UpdateUser(ctx context.Context, userID, chatID int64) error {
	const op = "storage.updateUser"

	query, args, err := sq.
		Update("users").
		Set("chat_ids", sq.Expr("array_remove(chat_ids, ?)", chatID)).
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = p.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *PostgresStorage) RemoveUser(ctx context.Context, userID, chatID int64) error {
	const op = "storage.removeUser"

	query, args, err := sq.
		Delete("users").
		Where(sq.Eq{"user_id": userID, "chat_ids": chatID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = p.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *PostgresStorage) AddSubscription(ctx context.Context, link *domain.TempUserState) error {
	const op = "storage.addSubscription"

	query, args, err := sq.
		Insert("subscriptions").
		Columns("chat_id", "chat_name", "user_id", "feature", "store", "created_at", "last_run_at").
		Values(link.ChatID, link.ChatName, link.UserID, link.Action, link.Store, time.Now(), time.Now()).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = p.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *PostgresStorage) UpdateSubscription(ctx context.Context, link *domain.Subscription) error {
	const op = "storage.updateSubscription"

	query, args, err := sq.
		Update("subscriptions").
		Set("store", link.Store).
		Where(sq.Eq{"chat_id": link.ChatID, "feature": link.Feature}).
		Where("store->>'tracked_user' = ?", link.Store.TrackedUser).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = p.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *PostgresStorage) RemoveSubscription(ctx context.Context, subID int) error {
	const op = "storage.removeSubscription"
	query, args, err := sq.
		Delete("subscriptions").
		Where(sq.Eq{"id": subID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = p.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *PostgresStorage) GetSubscriptions(ctx context.Context, chatID int64) (*[]domain.Subscription, error) {
	const op = "storage.getSubscriptions"

	query, args, err := sq.
		Select("id", "chat_id", "chat_name", "user_id", "feature", "store").
		From("subscriptions").
		Where(sq.Eq{"chat_id": chatID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := p.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var Subscriptions []domain.Subscription

	for rows.Next() {
		var s domain.Subscription
		if err := rows.Scan(&s.ID, &s.ChatID, &s.ChatName, &s.UserID, &s.Feature, &s.Store); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		Subscriptions = append(Subscriptions, s)
	}

	if Subscriptions == nil {
		return nil, fmt.Errorf("%s: no subscriptions found", op)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Subscriptions, nil
}

func (p *PostgresStorage) GetSubsWithFilter(ctx context.Context, chatID int64, userFilter, userData string) (*[]domain.Subscription, error) {
	const op = "storage.getSubsWithFilter"

	var (
		query string
		args  []interface{}
		err   error
	)

	switch userFilter {
	case "no-filter":
		query, args, err = sq.
			Select("id", "chat_id", "chat_name", "user_id", "feature", "store").
			From("subscriptions").
			Where(sq.Eq{"chat_id": chatID}).
			PlaceholderFormat(sq.Dollar).
			ToSql()
	case "user-filter":
		query, args, err = sq.
			Select("id", "chat_id", "chat_name", "user_id", "feature", "store").
			From("subscriptions").
			Where(sq.Eq{"chat_id": chatID}).
			Where("store->>'tracked_user' = ?", userData).
			PlaceholderFormat(sq.Dollar).
			ToSql()
	case "type-filter":
		actionFilter, err := strconv.Atoi(userData)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		query, args, err = sq.
			Select("id", "chat_id", "chat_name", "user_id", "feature", "store").
			From("subscriptions").
			Where(sq.Eq{"chat_id": chatID, "feature": actionFilter}).
			PlaceholderFormat(sq.Dollar).
			ToSql()
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := p.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var Subscriptions []domain.Subscription

	for rows.Next() {
		var s domain.Subscription
		if err := rows.Scan(&s.ID, &s.ChatID, &s.ChatName, &s.UserID, &s.Feature, &s.Store); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		Subscriptions = append(Subscriptions, s)
	}

	if Subscriptions == nil {
		return nil, fmt.Errorf("%s: no subscriptions found", op)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Subscriptions, nil
}

func (p *PostgresStorage) GetExpiredSubscriptions(ctx context.Context) (*[]domain.Subscription, error) {
	const op = "storage.getExpiredSubscriptions"

	query, args, err := sq.
		Select("id", "chat_id", "user_id", "feature", "store").
		From("subscriptions").
		Where("((store->>'start_time')::timestamp + (store->>'threshold' || ' minutes')::interval) <= LOCALTIMESTAMP").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := p.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()
	var Subscriptions []domain.Subscription
	for rows.Next() {
		var s domain.Subscription
		err := rows.Scan(&s.ID, &s.ChatID, &s.UserID, &s.Feature, &s.Store)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		Subscriptions = append(Subscriptions, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Subscriptions, nil
}

func (p *PostgresStorage) AddNewChat(ctx context.Context, chatID int64, chatName string) error {
	const op = "storage.addNewChat"

	query, args, err := sq.
		Insert("chats").
		Columns("chat_id", "chat_name").
		Values(chatID, chatName).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = p.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *PostgresStorage) RemoveChat(ctx context.Context, chatID int64) error {
	const op = "storage.removeChat"
	query, args, err := sq.
		Delete("chats").
		Where(sq.Eq{"chat_id": chatID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = p.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *PostgresStorage) GetUserFromChat(ctx context.Context, chatID int64, userName string) (*domain.User, error) {
	const op = "storage.getUserFromChat"

	var user domain.User

	userName = strings.TrimPrefix(userName, "@")

	query, args, err := sq.
		Select("user_id, username, chat_ids").
		From("users").
		Where(sq.Eq{"username": userName}).
		Where("chat_ids && ARRAY[?]::bigint[]", chatID).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	fmt.Println(chatID)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = p.DB.QueryRow(ctx, query, args...).Scan(&user.UserID, &user.UserName, &user.ChatID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("%s: chat not found", op)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (p *PostgresStorage) GetChat(ctx context.Context, chatName string) (*domain.Chat, error) {
	const op = "storage.getChat"

	query, args, err := sq.
		Select("chat_id", "chat_name").
		From("chats").
		Where(sq.Eq{"chat_name": chatName}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var chat domain.Chat

	err = p.DB.QueryRow(ctx, query, args...).Scan(&chat.ChatID, &chat.ChatName)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, fmt.Errorf("%s: chat not found", op)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &chat, nil
}

func (p *PostgresStorage) GetChats(ctx context.Context) ([]domain.Chat, error) {
	const op = "storage.getChats"

	query, args, err := sq.
		Select("chat_id", "chat_name").
		From("chats").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	row, err := p.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer row.Close()

	var chats []domain.Chat
	for row.Next() {
		var chatID int64
		var chatName string
		if err := row.Scan(&chatID, &chatName); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		chats = append(chats, domain.Chat{
			ChatID:   chatID,
			ChatName: chatName,
		})
	}

	if len(chats) == 0 {
		return nil, fmt.Errorf("%s: no chats", op)
	}

	return chats, nil
}

func (p *PostgresStorage) GetSubscription(ctx context.Context, ID int64, choice string) ([]*domain.Subscription, error) {
	const op = "storage.getUserSubscription"

	query, args, err := sq.
		Select("*").
		From("subscriptions").
		Where(sq.Eq{choice: ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := p.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var subs []*domain.Subscription

	for rows.Next() {
		var sub *domain.Subscription
		if err := rows.Scan(&ID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return subs, nil
}

func runMigrations(dbURL string) error {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return err
	}

	db := stdlib.OpenDB(*config.ConnConfig)
	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}

func (p *PostgresStorage) Close() error {
	p.DB.Close()
	return nil
}

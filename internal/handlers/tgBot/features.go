package tgBot

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	domain "horsey/internal/domain/entity"
	"math/rand/v2"
	"strings"
	"time"

	"gopkg.in/telebot.v4"
)

func (b *TgBot) SelectFeature(ctx context.Context, c telebot.Context) error {
	subs, err := b.useCase.GetSubs(ctx, c.Chat().ID)
	if err != nil {
		return err
	}

	userMsg := c.Message().Text
	userName := c.Sender().Username

	for _, sub := range *subs {
		if sub.Feature == domain.SendImage {
			err = sendImage(c, sub, userMsg, userName)
			if err != nil {
				return err
			}
		}
		if sub.Feature == domain.SendPing {
			err = b.UserTimer(ctx, c, sub)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func sendImage(c telebot.Context, sub domain.Subscription, userMsg, userName string) error {
	if strings.Contains(sub.Store.Keyword, userMsg) {
		chance := rand.IntN(sub.Store.Chance + 1)
		if chance == 1 {
			rawB64 := sub.Store.Image
			if i := strings.Index(rawB64, ","); i != -1 {
				rawB64 = rawB64[i+1:]
			}

			decoded, err := base64.StdEncoding.DecodeString(rawB64)
			if err != nil {
				return fmt.Errorf("base64 decoding failed: %w", err)
			}
			if sub.Store.TrackedUser != "" {
				if userName == sub.Store.TrackedUser {
					selectMedia(c, sub.Store, decoded)
				}
			} else {
				selectMedia(c, sub.Store, decoded)
			}
		}
	}

	return nil
}

func selectMedia(c telebot.Context, store domain.Store, decoded []byte) error {
	switch store.ImageType {
	case "photo":
		photo := &telebot.Photo{
			File: telebot.FromReader(bytes.NewReader(decoded)),
		}

		return c.Reply(photo)
	case "video":
		video := &telebot.Video{
			File: telebot.FromReader(bytes.NewReader(decoded)),
		}

		video.FileName = "video.mp4"

		return c.Reply(video)
	case "audio":
		audio := &telebot.Voice{
			File: telebot.FromReader(bytes.NewReader(decoded)),
		}

		return c.Reply(audio)
	case "voice":
		audio := &telebot.Voice{
			File: telebot.FromReader(bytes.NewReader(decoded)),
		}

		return c.Reply(audio)
	case "animation":
		anim := &telebot.Animation{
			File: telebot.FromReader(bytes.NewReader(decoded)),
		}

		anim.FileName = "animation.gif"

		return c.Reply(anim)
	}
	return nil
}

func (b *TgBot) SendPing(sub domain.Subscription) error {
	slova := []string{"очнись", "я вызываю тебя", "приди", "вернись...", "там новая фурри новелла вышла", "твое время пришло. Восстань!", "даже конченный идиот зарабатывает 1000$ в месяц.\nЧитать далее...", "ինչպես ես?", "тунг тунг тунг"}
	line := ""
	if sub.Store.LastMessage != "" {
		line += fmt.Sprintf("Последнее сообщение @%s было %s\n\nСообщение: %s", sub.Store.TrackedUser, sub.Store.StartTime, sub.Store.LastMessage)
	} else {
		line += fmt.Sprintf("@%s, %s", sub.Store.TrackedUser, slova[rand.IntN(len(slova))])
	}

	b.Bot.Send(&telebot.Chat{ID: sub.ChatID}, line)

	return nil
}

func (b *TgBot) BackgroundTimer(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	b.log.Info("BackgroundTimer worker started")

	for {
		select {
		case <-ctx.Done():
			b.log.Info("BackgroundTimer worker stopped")
			return
		case <-ticker.C:
			b.ProcessTimer(ctx)
		}
	}
}

func (b *TgBot) ProcessTimer(ctx context.Context) {
	expiredSubs, err := b.useCase.GetExpiredSubs(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		b.log.Error("ProcessTimer failed", "err", err)
		return
	}

	for _, sub := range *expiredSubs {
		b.SendPing(sub)
		sub.Store.StartTime = time.Now().UTC().Format("2006-01-02 15:04:05")
		b.useCase.UpdateSub(ctx, &sub)
	}
}

func (b *TgBot) UserTimer(ctx context.Context, c telebot.Context, sub domain.Subscription) error {
	if c.Sender().Username == sub.Store.TrackedUser {
		msgText := c.Message().Text
		msgTime := c.Message().Time
		sub.Store.StartTime = msgTime().UTC().Format("2006-01-02 15:04:05")
		sub.Store.LastMessage = msgText

		err := b.useCase.UpdateSub(ctx, &sub)
		if err != nil {
			b.log.Error("Cannot update sub store: ", err)
			return err
		}
	}

	return nil
}

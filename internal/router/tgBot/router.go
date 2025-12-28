package router

import (
	"context"
	"horsey/internal/handlers/tgBot"

	"gopkg.in/telebot.v4"
)

func Router(ctx context.Context, b *tgBot.TgBot) {
	b.Bot.Handle("/start", b.Start(ctx))
	b.Bot.Handle("Помощь", b.HelpMessage)
	b.Bot.Handle("Добавить связь", b.AddSubscription(ctx))
	b.Bot.Handle("Убрать связь", b.RemoveSubscription(ctx))
	b.Bot.Handle("Показать связи", b.ShowSubsWithFilter(ctx))
	b.Bot.Handle(telebot.OnMyChatMember, b.ChatManagement(ctx))
	b.Bot.Handle(telebot.OnText, b.MessageHandler(ctx))
	b.Bot.Handle(telebot.OnPhoto, b.MessageHandler(ctx))
	b.Bot.Handle(telebot.OnAudio, b.MessageHandler(ctx))
	b.Bot.Handle(telebot.OnVideo, b.MessageHandler(ctx))
	b.Bot.Handle(telebot.OnAnimation, b.MessageHandler(ctx))
	b.Bot.Handle(telebot.OnVoice, b.MessageHandler(ctx))
	b.Bot.Handle(telebot.OnUserLeft, b.UserLeft(ctx))
	b.Bot.Handle(telebot.OnUserJoined, b.UserJoined(ctx))
	b.Bot.Handle(telebot.OnCallback, b.HandleInlineButtonSelection(ctx))
}

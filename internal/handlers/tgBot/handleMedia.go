package tgBot

import (
	domain "horsey/internal/domain/entity"

	"gopkg.in/telebot.v4"
)

func (b *TgBot) HandleMedia(c telebot.Context, msg *telebot.Message) error {
	var file *telebot.File
	var mimeType string

	switch {
	case msg.Photo != nil:
		photo := msg.Photo
		file = &photo.File
		mimeType = photo.MediaType()
	case msg.Audio != nil:
		audio := msg.Audio
		file = &audio.File
		mimeType = audio.MediaType()
	case msg.Voice != nil:
		voice := msg.Voice
		file = &voice.File
		mimeType = voice.MediaType()
	case msg.Video != nil:
		video := msg.Video
		file = &video.File
		mimeType = video.MediaType()
	case msg.Animation != nil:
		anim := msg.Animation
		file = &anim.File
		mimeType = anim.MediaType()
	default:
		c.Send("Пожалуйста, отправьте изображение, видео, гифку или аудио файл.")
		return nil
	}

	var maxSize int64 = 20971520

	if file.FileSize > maxSize {
		return c.Send("Файл слишком большой, я не могу обрабатывать файлы размером больше 20МБ.")
	}

	reader, err := b.Bot.File(file)
	if err != nil {
		b.log.Error("Error while opening file: ", err)
		return err
	}
	defer reader.Close()

	dataurl, trueMimeType, err := b.useCase.HandleMedia(reader, mimeType)

	tempUserState[c.Sender().ID].Store.Image = dataurl
	tempUserState[c.Sender().ID].Store.ImageType = trueMimeType
	tempUserState[c.Sender().ID].State = domain.WaitingKeyword

	c.Send("Теперь отправь мне ключевое слово(-а), по которому будет создана связь.\n(Если хочешь использовать несколько слов, пиши их через пробел)")

	return nil
}

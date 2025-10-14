package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot    *tgbotapi.BotAPI
	chatID string
}

func NewTelegram(botToken, chatID string) (*TelegramBot, error) {

	if botToken == "" || chatID == "" {
		return nil, fmt.Errorf("token or chat ID is not set")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	msg := tgbotapi.NewMessageToChannel(chatID, "Bot starting...")
	_, err = bot.Send(msg)
	if err != nil {
		return nil, fmt.Errorf("not able to send starting message: %s", err)
	}
	return &TelegramBot{bot: bot, chatID: chatID}, nil
}

func (t *TelegramBot) Send(ctx context.Context, message string) error {
	msg := tgbotapi.NewMessageToChannel(t.chatID, message)
	_, err := t.bot.Send(msg)
	return err
}

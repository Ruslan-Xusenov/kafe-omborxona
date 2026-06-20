package service

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramService struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramService(token string) (*TelegramService, error) {
	if token == "" {
		return &TelegramService{}, nil // Bot disabled
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize telegram bot: %v", err)
	}

	log.Printf("🤖 Telegram Bot ulangan: %s", bot.Self.UserName)

	svc := &TelegramService{bot: bot}
	go svc.startPolling()

	return svc, nil
}

// Global Chat ID for admin notifications
var adminChatID int64

func (s *TelegramService) startPolling() {
	if s.bot == nil {
		return
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := s.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// When admin sends /start, save their chat ID
		if update.Message.Command() == "start" {
			adminChatID = update.Message.Chat.ID
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "✅ Salom! Men Kafe Omborxona tizimining botiman. Barcha muhim xabarlarni va xulosalarni shu yerga yuboraman.")
			s.bot.Send(msg)
			log.Printf("Telegram Admin Chat ID o'rnatildi: %d", adminChatID)
		}
	}
}

func (s *TelegramService) SendAdminMessage(text string) {
	if s.bot == nil || adminChatID == 0 {
		return
	}

	msg := tgbotapi.NewMessage(adminChatID, text)
	msg.ParseMode = "HTML"
	_, err := s.bot.Send(msg)
	if err != nil {
		log.Printf("Telegram xabar yuborishda xatolik: %v", err)
	}
}

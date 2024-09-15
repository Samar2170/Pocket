package main

import (
	"bufio"
	"context"
	"os"
	"pocket/internal"
	"pocket/internal/models"
	"pocket/pkg/auditlog"
	"pocket/pkg/db"
	"pocket/pkg/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

var (
	bot *tgbotapi.BotAPI
	// logger tgbotapi.BotLogger
	logger zerolog.Logger
)

var userIds = []int64{
	6983528406,
	6474112057,
	7543595397,
}

func init() {
	tgbotLogFile := &lumberjack.Logger{
		Filename:   "logs/tgbot.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   false,
	}
	logger = zerolog.New(tgbotLogFile).With().Timestamp().Logger()
	logger.Printf("Telegram server started")

}

func RunTelegramServer() {
	var err error
	bot, err = tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		auditlog.AuditLogger.Error().Str("error", err.Error()).Msg("Error creating bot")
		panic(err)
	}
	tgbotapi.SetLogger(&logger)
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	updates := bot.GetUpdatesChan(u)
	go recieveUpdates(ctx, updates)

	auditlog.AuditLogger.Println("Telegram server started")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()
}

func recieveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			handleUpdate(update)
		}
	}
}

func handleUpdate(update tgbotapi.Update) {
	auditlog.AuditLogger.Info().Msgf("Update: %+v", update)
	var err error
	message := update.Message
	user := message.From
	text := message.Text
	if user == nil {
		return
	}
	chatID := message.Chat.ID
	if !utils.CheckArray[int64](userIds, user.ID) {
		bot.Send(tgbotapi.NewMessage(chatID, "I don't know you!"))
		return
	}
	auditlog.AuditLogger.Info().Msgf("%s wrote %s", user.UserName, text)
	var msgString string
	if message.Photo != nil {
		files := message.Photo
		fileId, err := downloadAndSaveFileFromTg(files[len(files)-1].FileID)
		if err != nil {
			msgString = err.Error()
			auditlog.Errorlogger.Error().Str("error", err.Error()).Msg("Error saving file")
		}
		if message.Caption != "" {
			err = internal.SaveFileTags(fileId, message.Caption)
			if err != nil {
				msgString = err.Error()
				auditlog.Errorlogger.Error().Str("error", err.Error()).Msg("Error saving file caption")
			}
		}
	} else if message.Document != nil {
		fileId, err := downloadAndSaveFileFromTg(message.Document.FileID)
		if err != nil {
			msgString = err.Error()
			auditlog.Errorlogger.Error().Str("error", err.Error()).Msg("Error saving file")
		}
		if message.Caption != "" {
			err = internal.SaveFileTags(fileId, message.Caption)
			if err != nil {
				msgString = err.Error()
				auditlog.Errorlogger.Error().Str("error", err.Error()).Msg("Error saving file caption")
			}
		}
	} else if text != "" {
		err = db.DB.Save(&models.Note{NoteContent: text}).Error
		if err != nil {
			msgString = err.Error()
			auditlog.Errorlogger.Error().Str("error", err.Error()).Msg("Error saving note")
		}
	}
	if msgString != "" {
		bot.Send(tgbotapi.NewMessage(chatID, msgString))
	}
}

func downloadAndSaveFileFromTg(fileId string) (string, error) {
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileId})
	if err != nil {
		return "", err
	}
	fileUrl := file.Link(BotToken)
	return internal.SaveFileTelegram(fileUrl)
}

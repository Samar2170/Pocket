package main

import (
	"bufio"
	"context"
	"log"
	"os"
	"pocket/internal"
	"pocket/internal/models"
	"pocket/pkg/db"
	"pocket/pkg/utils"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var userIds = []int64{
	6983528406,
	6474112057,
}

var (
	// Menu texts
	firstMenu  = "<b>Menu 1</b>\n\nActions"
	secondMenu = "<b>Menu 2</b>\n\nSelect Account"

	// Button texts
	createContentButton      = "Create content"
	createImageContentButton = "Create image content"
	backButton               = "Back"
	accounts                 = []string{
		"sillybutcher1",
	}

	// Store bot screaming status
	takingInput = false
	bot         *tgbotapi.BotAPI

	// Keyboard layout for the first menu. One button, one row
	firstMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(createContentButton, createContentButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(createImageContentButton, createImageContentButton),
		),
	)

	// Keyboard layout for the second menu. Two buttons, one per row
)

var accountMenuMarkup tgbotapi.InlineKeyboardMarkup

func init() {
	for _, account := range accounts {
		accountMenuMarkup.InlineKeyboard = append(accountMenuMarkup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(account, account),
		))
	}

	accountMenuMarkup.InlineKeyboard = append(accountMenuMarkup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(backButton, backButton),
	))
}

func RunTelegramServer() {
	var err error
	bot, err = tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		panic(err)
	}
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	updates := bot.GetUpdatesChan(u)
	go recieveUpdates(ctx, updates)

	log.Println("Telegram server started")
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
	var accountUsername string
	var account models.Account
	var err error
	if update.CallbackQuery != nil {
		if utils.CheckArray(accounts, accountUsername) {
			account, err = models.GetAccountByUsername(accountUsername)
			if err != nil {
				log.Println(err)
			}
		}

	}

	switch {
	case update.CallbackQuery != nil:
		handleButton(update.CallbackQuery)
		return
	case update.Message != nil:
		handleMessage(update.Message, account)
		return
	}
}

func handleMessage(message *tgbotapi.Message, account models.Account) {
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
	log.Printf("%s wrote %s", user.UserName, text)
	var err error
	if strings.HasPrefix(text, "/") {
		err = handleCommand(message.Chat.ID, text)
	} else if takingInput && len(text) > 0 {
		var msg tgbotapi.MessageConfig
		err = db.DB.Create(&models.TextContent{Text: text, Account: account}).Error
		if err != nil {
			msg = tgbotapi.NewMessage(chatID, "Something went wrong")
		} else {
			msg = tgbotapi.NewMessage(message.Chat.ID, "input received")
		}
		msg.Entities = message.Entities
		_, err = bot.Send(msg)
		takingInput = false
	} else if takingInput && message.Photo != nil {
		files := message.Photo
		var msgString string
		imageUrl, err := downloadAndSaveFileFromTg(files[len(files)-1].FileID)
		if err != nil {
			msgString = "Something went wrong while downloading image"
			log.Println(err)
		}
		err = db.DB.Create(&models.ImageContent{ImageURL: imageUrl, Account: account}).Error
		if err != nil {
			msgString = "Something went wrong while saving image to db"
			log.Println(err)
		}
		if msgString == "" {
			msgString = "Image saved successfully"
		}
		msg := tgbotapi.NewMessage(message.Chat.ID, msgString)
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
		takingInput = false
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "I don't know what you want")
		_, err = bot.Send(msg)
	}
	if err != nil {
		log.Println(err)
	}
	err = sendMenu(chatID)
	if err != nil {
		log.Println(err)
	}
}

func handleCommand(chatId int64, command string) error {
	var err error
	switch command {
	case "/start":
		err = sendMenu(chatId)
		break
	case "/menu":
		err = sendMenu(chatId)
		break
	}
	return err
}
func handleButton(query *tgbotapi.CallbackQuery) {
	var text string
	var err error
	markup := tgbotapi.NewInlineKeyboardMarkup()
	message := query.Message

	switch {
	case query.Data == createContentButton || query.Data == createImageContentButton:
		text = secondMenu
		markup = accountMenuMarkup
	case utils.CheckArray(accounts, query.Data):
		takingInput = true
		msg := tgbotapi.NewMessage(query.Message.Chat.ID, "waiting for input")
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	case query.Data == backButton:
		text = firstMenu
		markup = firstMenuMarkup
	}

	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	bot.Send(callbackCfg)

	// Replace menu text and keyboard
	msg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	bot.Send(msg)
}

func sendMenu(chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, firstMenu)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = firstMenuMarkup
	_, err := bot.Send(msg)
	return err
}
func downloadAndSaveFileFromTg(fileId string) (string, error) {
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileId})
	if err != nil {
		return "", err
	}
	fileUrl := file.Link(BotToken)
	return internal.SaveImageTelegram(fileUrl)

}

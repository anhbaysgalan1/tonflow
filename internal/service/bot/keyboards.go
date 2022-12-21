package bot

import (
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"park-wallet/internal/service/bot/template"
)

var (
	mainKeyboard = tgBotAPI.NewReplyKeyboard(
		tgBotAPI.NewKeyboardButtonRow(
			tgBotAPI.NewKeyboardButton(template.ReceiveButton),
			tgBotAPI.NewKeyboardButton(template.SendButton),
		), tgBotAPI.NewKeyboardButtonRow(
			tgBotAPI.NewKeyboardButton(template.BalanceButton),
		))
	confirmKeyboard = tgBotAPI.NewReplyKeyboard(
		tgBotAPI.NewKeyboardButtonRow(
			tgBotAPI.NewKeyboardButton(template.ConfirmButton),
			tgBotAPI.NewKeyboardButton(template.CancelButton),
		))
	cancelKeyboard = tgBotAPI.NewReplyKeyboard(
		tgBotAPI.NewKeyboardButtonRow(
			tgBotAPI.NewKeyboardButton(template.CancelButton),
		))

	depositInline = tgBotAPI.NewInlineKeyboardMarkup(
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonURL(
				"Tonkeeper",
				"https://app.tonkeeper.com/transfer/kQDgEEX3G0xKggopmwKrLowIR_QMxq-zgRqA9jF6JSi5DRl_"),
			tgBotAPI.NewInlineKeyboardButtonURL(
				"Tonhub",
				"https://tonhub.com/transfer/kQDgEEX3G0xKggopmwKrLowIR_QMxq-zgRqA9jF6JSi5DRl_"),
		),
	)
)

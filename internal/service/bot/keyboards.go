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

	addressWithQR = tgBotAPI.NewInlineKeyboardMarkup(
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonSwitch(
				"Tonkeeper",
				"https://app.tonkeeper.com/transfer/kQDgEEX3G0xKggopmwKrLowIR_QMxq-zgRqA9jF6JSi5DRl_?amount=1000000000&text=Comment"),
		),
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonURL(
				"Tonhub",
				"https://tonhub.com/transfer/kQDgEEX3G0xKggopmwKrLowIR_QMxq-zgRqA9jF6JSi5DRl_?amount=1000000000&text=Comment"),
		),
	)
)

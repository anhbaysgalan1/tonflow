package bot

import telegramBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var (
	mainKeyboard = telegramBotAPI.NewReplyKeyboard(
		telegramBotAPI.NewKeyboardButtonRow(
			telegramBotAPI.NewKeyboardButton("📥 Receive"),
			telegramBotAPI.NewKeyboardButton("📤 Send"),
		), telegramBotAPI.NewKeyboardButtonRow(
			telegramBotAPI.NewKeyboardButton("💎 Balance"),
		))

	addressWithQR = telegramBotAPI.NewInlineKeyboardMarkup(
		telegramBotAPI.NewInlineKeyboardRow(
			telegramBotAPI.NewInlineKeyboardButtonSwitch(
				"Tonkeeper",
				"https://app.tonkeeper.com/transfer/kQDgEEX3G0xKggopmwKrLowIR_QMxq-zgRqA9jF6JSi5DRl_?amount=1000000000&text=Comment"),
		),
		telegramBotAPI.NewInlineKeyboardRow(
			telegramBotAPI.NewInlineKeyboardButtonURL(
				"Tonhub",
				"https://tonhub.com/transfer/kQDgEEX3G0xKggopmwKrLowIR_QMxq-zgRqA9jF6JSi5DRl_?amount=1000000000&text=Comment"),
		),
	)
)

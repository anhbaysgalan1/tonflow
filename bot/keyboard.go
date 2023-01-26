package bot

import (
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	inlineMainKeyboard = tgBotAPI.NewInlineKeyboardMarkup(
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData(
				ReceiveButton,
				"receive"),
			tgBotAPI.NewInlineKeyboardButtonData(
				SendButton,
				"send"),
		),
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData(
				BalanceButton,
				"balance"),
		))

	inlineBalanceKeyboard = tgBotAPI.NewInlineKeyboardMarkup(
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData(
				ReceiveButton,
				"receive"),
			tgBotAPI.NewInlineKeyboardButtonData(
				SendButton,
				"send"),
		),
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData(
				UpdateBalanceButton,
				"update balance"),
		))

	inlineConfirmKeyboard = tgBotAPI.NewInlineKeyboardMarkup(
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData(
				ConfirmButton,
				"confirm"),
			tgBotAPI.NewInlineKeyboardButtonData(
				CancelButton,
				"cancel"),
		),
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData(
				AddCommentButton,
				"add comment"),
		),
	)
)

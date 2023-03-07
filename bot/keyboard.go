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

	inlineReceiveSendKeyboard = tgBotAPI.NewInlineKeyboardMarkup(
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData(
				ReceiveButton,
				"receive"),
			tgBotAPI.NewInlineKeyboardButtonData(
				SendButton,
				"send"),
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

	inlineSendAllKeyboard = tgBotAPI.NewInlineKeyboardMarkup(
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData(
				SendAllButton,
				"send all"),
		),
	)

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

	inlineConfirmWithCommentKeyboard = tgBotAPI.NewInlineKeyboardMarkup(
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData(
				ConfirmButton,
				"confirm"),
			tgBotAPI.NewInlineKeyboardButtonData(
				CancelButton,
				"cancel"),
		),
	)
)

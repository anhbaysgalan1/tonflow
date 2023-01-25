package bot

import (
	"encoding/json"
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	mainKeyboard = tgBotAPI.NewReplyKeyboard(
		tgBotAPI.NewKeyboardButtonRow(
			tgBotAPI.NewKeyboardButton(ReceiveButton),
			tgBotAPI.NewKeyboardButton(SendButton),
		), tgBotAPI.NewKeyboardButtonRow(
			tgBotAPI.NewKeyboardButton(BalanceButton),
		))
	mainInlineKeyboard = tgBotAPI.NewInlineKeyboardMarkup(
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

	mainInlineKeyboardCheckBalance = tgBotAPI.NewInlineKeyboardMarkup(
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

	confirmKeyboard = tgBotAPI.NewReplyKeyboard(
		tgBotAPI.NewKeyboardButtonRow(
			tgBotAPI.NewKeyboardButton(ConfirmButton),
			tgBotAPI.NewKeyboardButton(CancelButton),
		))
	confirmInlineKeyboard = tgBotAPI.NewInlineKeyboardMarkup(
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

	cancelKeyboard = tgBotAPI.NewReplyKeyboard(
		tgBotAPI.NewKeyboardButtonRow(
			tgBotAPI.NewKeyboardButton(CancelButton),
		))
	cancelInlineKeyboard = tgBotAPI.NewInlineKeyboardMarkup(
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData(
				CancelButton,
				"cancel"),
		))

	depositInline = tgBotAPI.NewInlineKeyboardMarkup(
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonURL(
				"Tonkeeper",
				"https://app.tonkeeper.com/transfer/EQBTTcHwhP7v-U59QZ-QO9c_BAvum0VfvtSdI7wBT9zKbFOy"),
			tgBotAPI.NewInlineKeyboardButtonURL(
				"Tonhub",
				"https://tonhub.com/transfer/EQBTTcHwhP7v-U59QZ-QO9c_BAvum0VfvtSdI7wBT9zKbFOy"),
		),
	)
)

type buttonData struct {
	Method string `json:"method,omitempty"`
	Wallet string `json:"wallet,omitempty"`
}

func newButtonData(method, wallet string) (string, error) {
	data := &buttonData{
		Method: method,
		Wallet: wallet,
	}
	result, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func newInlineMainButtons(data string) tgBotAPI.InlineKeyboardMarkup {

	return tgBotAPI.NewInlineKeyboardMarkup(
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData(
				ReceiveButton,
				"receive"),
			tgBotAPI.NewInlineKeyboardButtonData(
				SendButton,
				data),
		),
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData(
				BalanceButton,
				"balance"),
		))
}

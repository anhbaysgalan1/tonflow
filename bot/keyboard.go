package bot

import (
	"encoding/json"
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tonflow/bot/template"
)

var (
	mainKeyboard = tgBotAPI.NewReplyKeyboard(
		tgBotAPI.NewKeyboardButtonRow(
			tgBotAPI.NewKeyboardButton(template.ReceiveButton),
			tgBotAPI.NewKeyboardButton(template.SendButton),
		), tgBotAPI.NewKeyboardButtonRow(
			tgBotAPI.NewKeyboardButton(template.BalanceButton),
		))
	mainInlineKeyboard = tgBotAPI.NewInlineKeyboardMarkup(
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData(
				template.ReceiveButton,
				"receive"),
			tgBotAPI.NewInlineKeyboardButtonData(
				template.SendButton,
				"send"),
		),
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData(
				template.BalanceButton,
				"balance"),
		))

	mainInlineKeyboardCheckBalance = tgBotAPI.NewInlineKeyboardMarkup(
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData(
				template.ReceiveButton,
				"receive"),
			tgBotAPI.NewInlineKeyboardButtonData(
				template.SendButton,
				"send"),
		),
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData(
				template.UpdateBalanceButton,
				"update balance"),
		))

	confirmKeyboard = tgBotAPI.NewReplyKeyboard(
		tgBotAPI.NewKeyboardButtonRow(
			tgBotAPI.NewKeyboardButton(template.ConfirmButton),
			tgBotAPI.NewKeyboardButton(template.CancelButton),
		))
	confirmInlineKeyboard = tgBotAPI.NewInlineKeyboardMarkup(
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData(
				template.ConfirmButton,
				"confirm"),
			tgBotAPI.NewInlineKeyboardButtonData(
				template.CancelButton,
				"cancel"),
		),
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData(
				template.AddCommentButton,
				"add comment"),
		),
	)

	cancelKeyboard = tgBotAPI.NewReplyKeyboard(
		tgBotAPI.NewKeyboardButtonRow(
			tgBotAPI.NewKeyboardButton(template.CancelButton),
		))
	cancelInlineKeyboard = tgBotAPI.NewInlineKeyboardMarkup(
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData(
				template.CancelButton,
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
				template.ReceiveButton,
				"receive"),
			tgBotAPI.NewInlineKeyboardButtonData(
				template.SendButton,
				data),
		),
		tgBotAPI.NewInlineKeyboardRow(
			tgBotAPI.NewInlineKeyboardButtonData(
				template.BalanceButton,
				"balance"),
		))
}

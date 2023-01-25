package bot

import (
	"context"
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"tonflow/bot/model"
	"tonflow/pkg"
)

func (bot *Bot) handleUpdate(ctx context.Context, update tgBotAPI.Update) {
	switch {
	case update.CallbackQuery != nil:
		bot.handleCallbackQuery(ctx, update)
	case update.Message != nil:
		bot.handleMessage(ctx, update)
	}
}

func (bot *Bot) handleCallbackQuery(ctx context.Context, update tgBotAPI.Update) {
	user, _, err := bot.getTonflowUser(ctx, update.SentFrom())
	if err != nil {
		bot.sendErr(err, "getTonflowUser")
		// возможно, в случае ошибок, нужно отвечать сообщением, что что-то пошло не так
		return
	}

	switch update.CallbackData() {
	case "receive":
		bot.inlineReceiveCoins(user)
	case "send":
		bot.inlineSendCoins(ctx, update, user)
	case "balance":
		bot.inlineBalance(update, user)
	case "update balance":
		bot.inlineUpdateBalance(update, user)
	case "cancel":
		bot.inlineCancel(ctx, update, user)
		bot.inlineBalance(update, user)
	}
}

func (bot *Bot) handleMessage(ctx context.Context, update tgBotAPI.Update) {
	user, isExisted, err := bot.getTonflowUser(ctx, update.SentFrom())
	if err != nil {
		bot.sendErr(err, "getTonflowUser")
		// возможно, в случае ошибок, нужно отвечать сообщением, что что-то пошло не так
		return
	}

	log.Debugf("getTonflowUser():\n%v\nisExisted: %v", pkg.AnyPrint(user), pkg.AnyPrint(isExisted))

	switch {
	case update.Message.IsCommand():
		switch update.Message.Command() {
		case "start":
			bot.cmdStart(update, user, isExisted)
		}
	default:
		if user.StageData == nil {
			log.Warnf("nil user.StageData of user %v", pkg.AnyPrint(user.ID))
		}
		if user.StageData.Stage == model.AddressWait && len(update.Message.Photo) != 0 {
			bot.parseQR(ctx, update, user)
		}
		if user.StageData.Stage == model.AmountWait && update.Message.Text != "" {
			bot.validateAmount(ctx, update, user)
		}
	}
}

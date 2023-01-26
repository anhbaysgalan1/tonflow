package bot

import (
	"context"
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"tonflow/model"
	"tonflow/pkg"
)

func (bot *Bot) handleUpdate(ctx context.Context, update tgBotAPI.Update) {
	user, isExisted, err := bot.getTonflowUser(ctx, update.SentFrom())
	if err != nil {
		log.Error(err)
		return
	}

	log.Debugf("getTonflowUser():\n%v\nisExisted: %v", pkg.AnyPrint(user), pkg.AnyPrint(isExisted))

	switch {
	case update.Message != nil:
		bot.handleMessage(ctx, update, user, isExisted)
	case update.CallbackQuery != nil:
		bot.handleCallback(ctx, update, user)
	}
}

func (bot *Bot) handleMessage(ctx context.Context, update tgBotAPI.Update, user *model.User, isExisted bool) {
	switch {
	case update.Message.IsCommand():
		switch update.Message.Command() {
		case "start":
			bot.cmdStart(update, user, isExisted)
		}
	default:
		if user.StageData.Stage == model.AddressWait && len(update.Message.Photo) != 0 {
			bot.parseSendingQR(ctx, update, user)
		}
		if user.StageData.Stage == model.AmountWait && update.Message.Text != "" {
			bot.validateSendingAmount(ctx, update, user)
		}
		// дефолтное сообщение не удовлетворяеющее ни однму стейджу
		// ...
	}
}

func (bot *Bot) handleCallback(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	switch update.CallbackData() {
	case "receive":
		bot.inlineReceiveCoins(update, user)
	case "send":
		bot.inlineSendCoins(ctx, update, user)
	case "balance":
		bot.inlineBalance(update, user)
	case "update balance":
		bot.inlineBalanceUpdate(update, user)
	case "cancel":
		bot.inlineCancel(ctx, update, user)
		bot.inlineBalance(update, user)
	}
}

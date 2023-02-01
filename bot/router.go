package bot

import (
	"context"
	"fmt"
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"tonflow/model"
	"tonflow/pkg"
)

func (bot *Bot) handleUpdate(ctx context.Context, update tgBotAPI.Update) {
	user := &model.User{}
	isExisted := false
	var err error
	if update.Message != nil || update.CallbackQuery != nil {
		user, isExisted, err = bot.getTonflowUser(ctx, update.SentFrom())
		if err != nil {
			log.Error(err)
			return
		}
		log.Debugf("getTonflowUser():\n%v\nisExisted: %v", pkg.AnyPrint(user), pkg.AnyPrint(isExisted))
	}
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
		switch {
		case user.StageData.Stage == model.AddressWait:
			bot.acceptSendingAddress(ctx, update, user)
		case user.StageData.Stage == model.AmountWait:
			bot.acceptSendingAmount(ctx, update, user)
		case user.StageData.Stage == model.CommentWait:
			bot.acceptComment(ctx, update, user)
		default:
			/// need implement this right
			err := bot.sendText(update.Message.Chat.ID, "Nothing to do...", nil)
			if err != nil {
				log.Error(err)
				return
			}
			///
		}
	}
}

func (bot *Bot) handleCallback(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	switch update.CallbackData() {
	case "receive":
		bot.inlineReceiveCoins(update, user)
	case "send":
		bot.inlineSendCoins(ctx, update, user)
	case "balance":
		bot.inlineBalance(ctx, update, user)
	case "update balance":
		bot.inlineBalanceUpdate(ctx, update, user)
	case "cancel":
		bot.inlineCancel(ctx, update, user)
		bot.inlineBalance(ctx, update, user)
	case "add comment":
		bot.inlineAddComment(ctx, update, user)
	case "send all":
		bot.inlineSendAll(ctx, update, user)
	case "confirm":
		bot.confirmSending(ctx, update, user)
		bot.inlineBalance(ctx, update, user)
	default:
		log.Error(fmt.Errorf("need handle this case"))
	}
}

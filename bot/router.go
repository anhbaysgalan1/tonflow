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
	var (
		message  = update.Message
		callback = update.CallbackQuery
		user     *model.User
		err      error
	)

	if message != nil || callback != nil {
		user, err = bot.getTonflowUser(ctx, update.SentFrom())
		if err != nil {
			log.Error(err)
			return
		}

		log.Debugf("getTonflowUser():\n%s", pkg.PrintAny(user))

		switch {
		case message != nil:
			bot.handleMessage(ctx, update, user)
		case callback != nil:
			bot.handleCallback(ctx, update, user)
		}
	}
}

func (bot *Bot) handleMessage(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	switch {
	case update.Message.IsCommand():
		switch update.Message.Command() {
		case "start":
			bot.start(update, user)
		case "balance":
			bot.balance(ctx, update, user)
		case "receive":
			bot.receiveCoins(update, user)
		case "send":
			bot.sendCoins(ctx, update, user)
		case "cancel":
			bot.cancel(ctx, update, user)
		default:
			err := bot.sendText(update.Message.Chat.ID, "There is no command like this", nil)
			if err != nil {
				log.Error(err)
				return
			}
		}
	default:
		switch user.StageData.Stage {
		case model.AddressWait:
			bot.setAddress(ctx, update, user)
		case model.AmountWait:
			bot.setAmount(ctx, update, user)
		case model.CommentWait:
			bot.setComment(ctx, update, user)
		default:
			err := bot.sendText(update.Message.Chat.ID, "Nothing to do with this message...", nil)
			if err != nil {
				log.Error(err)
				return
			}
		}
	}
}

func (bot *Bot) handleCallback(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	switch update.CallbackData() {
	case "balance":
		bot.balance(ctx, update, user)
	case "receive":
		bot.receiveCoins(update, user)
	case "send":
		bot.sendCoins(ctx, update, user)
	case "add comment":
		bot.addComment(ctx, update, user)
	case "send all":
		bot.sendAll(ctx, update, user)
	case "confirm":
		bot.confirm(ctx, update, user)
	case "cancel":
		bot.cancel(ctx, update, user)
	default:
		log.Warning(fmt.Errorf("unsupported callback data"))
	}
}

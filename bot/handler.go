package bot

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	tgBotAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
	"github.com/xssnick/tonutils-go/tlb"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"tonflow/model"
	"tonflow/pkg"
)

func (bot *Bot) getTonflowUser(ctx context.Context, tgUser *tgBotAPI.User) (*model.User, error) {
	user, err := bot.redis.GetUserCache(ctx, tgUser.ID)
	switch {
	// error
	case err != nil && err != redis.Nil:
		log.Error(err)
		return nil, err

	// no user cache
	case err == redis.Nil:
		user, err = bot.checkUserDB(ctx, tgUser)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		err = bot.redis.SetUserCache(ctx, user)
		if err != nil {
			log.Error(err)
			return nil, err
		}

	// has user cache
	default:
		user.IsExisted = true
	}

	return user, nil
}

func (bot *Bot) checkUserDB(ctx context.Context, tgUser *tgBotAPI.User) (*model.User, error) {
	user, err := bot.storage.GetUser(ctx, tgUser.ID)
	switch {
	// error
	case err != nil && !errors.Is(err, pgx.ErrNoRows):
		log.Error(err)
		return nil, err

	// user is not existed
	case errors.Is(err, pgx.ErrNoRows):
		wlt, err := bot.ton.NewWallet()
		if err != nil {
			log.Error(err)
			return nil, err
		}
		log.Debugf("seed: %v", wlt.Seed)

		wlt.Seed, err = pkg.EncodeAES(wlt.Seed, strconv.FormatInt(tgUser.ID, 10))
		if err != nil {
			log.Error(err)
			return nil, err
		}
		log.Debugf("encrypted seed: %v", wlt.Seed)

		err = bot.storage.AddUser(ctx, tgUser)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		err = bot.storage.AddWallet(ctx, wlt, tgUser.ID)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		user = &model.User{
			ID:           tgUser.ID,
			Username:     tgUser.UserName,
			FirstName:    tgUser.FirstName,
			LastName:     tgUser.LastName,
			LanguageCode: tgUser.LanguageCode,
			Wallet:       wlt,
			IsExisted:    false,
			StageData:    &model.EmptyStageData,
		}

	// has user in db
	default:
		wallet, err := bot.storage.GetWallet(ctx, user.Wallet.Address)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		user.Wallet = wallet
		user.IsExisted = true
		user.StageData = &model.EmptyStageData

	}

	return user, nil
}

func (bot *Bot) start(update tgBotAPI.Update, user *model.User) {
	address := user.Wallet.Address

	qr, err := pkg.EncodeQR(address)
	if err != nil {
		log.Errorf("failed to encode qr: %s", err)
		return
	}

	txt := ""
	firstName := update.Message.From.FirstName
	switch user.IsExisted {
	case true:
		txt = fmt.Sprintf(WelcomeExistedUser, firstName, address)
	case false:
		txt = fmt.Sprintf(WelcomeNewUser, firstName, address)
	}

	if err = bot.sendImage(update.Message.Chat.ID, qr, txt, inlineMainKeyboard); err != nil {
		log.Error(err)
	}
}

func (bot *Bot) balance(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	balance, err := bot.ton.GetWalletBalance(ctx, user.Wallet.Address)
	if err != nil {
		log.Error(err)
		return
	}

	if err = bot.sendText(user.ID, fmt.Sprintf(Balance, balance), inlineReceiveSendKeyboard); err != nil {
		log.Error(err)
	}

	if update.CallbackQuery != nil {
		cb := tgBotAPI.NewCallback(update.CallbackQuery.ID, "")
		_, err = bot.api.Request(cb)
		if err != nil {
			log.Error(err)
		}
	}
}

func (bot *Bot) receiveCoins(update tgBotAPI.Update, user *model.User) {
	address := user.Wallet.Address

	qr, err := pkg.EncodeQR(address)
	if err != nil {
		log.Errorf("failed to encode qr: %s", err)
		return
	}

	caption := fmt.Sprintf(ReceiveInstruction, address)
	if err := bot.sendImage(user.ID, qr, caption, inlineMainKeyboard); err != nil {
		log.Error(err)
	}

	if update.CallbackQuery != nil {
		cb := tgBotAPI.NewCallback(update.CallbackQuery.ID, "")
		_, err := bot.api.Request(cb)
		if err != nil {
			log.Error(err)
		}
	}
}

func (bot *Bot) sendCoins(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	address := user.Wallet.Address

	balance, err := bot.ton.GetWalletBalance(ctx, address)
	if err != nil {
		log.Error(err)
		return
	}

	chatID := user.ID

	if balance == "0" {
		if err = bot.sendText(chatID, NoFunds, nil); err != nil {
			log.Error(err)
			return
		}
		return
	}

	user.StageData.Stage = model.AddressWait
	if err = bot.redis.SetUserCache(ctx, user); err != nil {
		log.Error(err)
		return
	}

	if err = bot.sendText(chatID, AskWallet, nil); err != nil {
		log.Error(err)
	}

	if update.CallbackQuery != nil {
		cb := tgBotAPI.NewCallback(update.CallbackQuery.ID, "")
		_, err = bot.api.Request(cb)
		if err != nil {
			log.Error(err)
		}
	}
}

func (bot *Bot) setAddress(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	chatID := update.Message.Chat.ID
	photos := update.Message.Photo
	addr := update.Message.Text

	if len(photos) != 0 {
		index, size := 0, 0
		for i, v := range photos {
			if v.FileSize > size {
				index = i
			}
		}

		fileURL, err := bot.api.GetFileDirectURL(photos[index].FileID)
		if err != nil {
			log.Error(err)
			return
		}

		resp, err := http.Get(fileURL)
		if err != nil {
			log.Error(err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Error(err)
			return
		}

		addr, err = pkg.DecodeQR(resp.Body)
		if err != nil {
			log.Warning(err)
			if err = bot.sendText(chatID, InvalidQR, nil); err != nil {
				log.Error(err)
				return
			}
			return
		}
		log.Debugf("parsed QR: %s", addr)
	}

	err := bot.ton.ValidateWallet(addr)
	if err != nil {
		log.Warnf("failed to validate address %v: %v", addr, err)
		if err = bot.sendText(chatID, InvalidWallet, nil); err != nil {
			log.Error(err)
			return
		}
		return
	}

	user.StageData.Stage = model.AmountWait
	user.StageData.AddressToSend = addr

	err = bot.redis.SetUserCache(ctx, user)
	if err != nil {
		log.Error(err)
		return
	}

	if err = bot.sendText(chatID, AskAmount, inlineSendAllKeyboard); err != nil {
		log.Error(err)
	}
}

func (bot *Bot) sendAll(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	balance, err := bot.ton.GetWalletBalance(ctx, user.Wallet.Address)
	if err != nil {
		log.Error(err)
		return
	}

	user.StageData.Stage = model.ConfirmationWait
	user.StageData.AmountToSend = balance
	user.StageData.SendAll = true

	err = bot.redis.SetUserCache(ctx, user)
	if err != nil {
		log.Error(err)
		return
	}

	txt := fmt.Sprintf(SendingConfirmation, user.StageData.AddressToSend, balance, bot.blockchainTxFee)
	if err = bot.sendText(update.CallbackQuery.Message.Chat.ID, txt, inlineConfirmKeyboard); err != nil {
		log.Error(err)
	}

	if update.CallbackQuery != nil {
		cb := tgBotAPI.NewCallback(update.CallbackQuery.ID, "")
		_, err = bot.api.Request(cb)
		if err != nil {
			log.Error(err)
		}
	}

}

func (bot *Bot) setAmount(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	address := user.Wallet.Address

	balance, err := bot.ton.GetWalletBalance(ctx, address)
	if err != nil {
		log.Error(err)
		return
	}

	balanceInCoins, err := tlb.FromTON(balance)
	if err != nil {
		log.Error()
		return
	}

	feesInCoins, err := tlb.FromTON(bot.blockchainTxFee)
	if err != nil {
		log.Error()
		return
	}

	amount := strings.ReplaceAll(update.Message.Text, " ", "")
	amount = strings.ReplaceAll(amount, ",", ".")
	chatID := update.Message.Chat.ID

	amountInCoins, err := tlb.FromTON(amount)
	if err != nil {
		log.Warnf("failed to convert %v to coins: %v", amountInCoins, err)
		if err = bot.sendText(chatID, InvalidAmount, struct{}{}); err != nil {
			log.Error(err)
			return
		}
		return
	}

	amountWithFees := big.NewInt(0).Add(amountInCoins.NanoTON(), feesInCoins.NanoTON())

	if balanceInCoins.NanoTON().Cmp(amountWithFees) == -1 {
		txt := fmt.Sprintf(NotEnoughFunds, balance, bot.blockchainTxFee)
		if err = bot.sendText(chatID, txt, inlineSendAllKeyboard); err != nil {
			return
		}
		return
	}

	user.StageData.Stage = model.ConfirmationWait
	user.StageData.AmountToSend = amount
	err = bot.redis.SetUserCache(ctx, user)
	if err != nil {
		log.Error(err)
		return
	}

	txt := fmt.Sprintf(SendingConfirmation, user.StageData.AddressToSend, amount, bot.blockchainTxFee)
	if err = bot.sendText(chatID, txt, inlineConfirmKeyboard); err != nil {
		log.Error(err)
	}
}

func (bot *Bot) addComment(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	user.StageData.Stage = model.CommentWait
	err := bot.redis.SetUserCache(ctx, user)
	if err != nil {
		log.Error(err)
		return
	}

	if err = bot.sendText(update.CallbackQuery.Message.Chat.ID, AskComment, nil); err != nil {
		log.Error(err)
	}

	cb := tgBotAPI.NewCallback(update.CallbackQuery.ID, "")
	_, err = bot.api.Request(cb)
	if err != nil {
		log.Error(err)
	}
}

func (bot *Bot) setComment(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	comment := update.Message.Text

	user.StageData.Stage = model.ConfirmationWait
	user.StageData.Comment = comment
	err := bot.redis.SetUserCache(ctx, user)
	if err != nil {
		log.Error(err)
		return
	}

	txt := fmt.Sprintf(
		SendingConfirmation,
		user.StageData.AddressToSend,
		user.StageData.AmountToSend,
		bot.blockchainTxFee,
	) + fmt.Sprintf(Comment, comment)
	if err = bot.sendText(update.Message.Chat.ID, txt, inlineConfirmWithCommentKeyboard); err != nil {
		log.Error(err)
	}
}

func (bot *Bot) confirm(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	if user.StageData.Stage != model.ConfirmationWait {
		return
	}

	txt := strings.ReplaceAll(update.CallbackQuery.Message.Text, "address: ", "address: <code>")
	txt = strings.ReplaceAll(txt, "\n\nAmount", "</code>\n\nAmount")
	confirmed := txt + SendingCoins

	msg := tgBotAPI.EditMessageTextConfig{
		BaseEdit: tgBotAPI.BaseEdit{
			ChatID:    update.CallbackQuery.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.MessageID,
		},
		Text:      confirmed,
		ParseMode: "HTML",
	}

	_, err := bot.api.Request(msg)
	if err != nil {
		log.Error(err)
		return
	}

	if user.StageData.SendAll == true {
		err := bot.ton.SendAll(ctx, user)
		if err != nil {
			log.Error(err)
			return
		}
	} else {
		err = bot.ton.Send(ctx, user)
		if err != nil {
			log.Error(err)
			return
		}
	}

	user.StageData = &model.EmptyStageData
	err = bot.redis.SetUserCache(ctx, user)
	if err != nil {
		log.Error(err)
		return
	}

	msg.Text = txt + fmt.Sprintf(Sent)

	_, err = bot.api.Request(msg)
	if err != nil {
		log.Error(err)
		return
	}

	bot.balance(ctx, update, user)
}

func (bot *Bot) cancel(ctx context.Context, update tgBotAPI.Update, user *model.User) {
	if update.CallbackQuery != nil {
		txt := strings.ReplaceAll(update.CallbackQuery.Message.Text, "address: ", "address: <code>")
		txt = strings.ReplaceAll(txt, "\n\nAmount", "</code>\n\nAmount")
		txt += Canceled

		msg := tgBotAPI.EditMessageTextConfig{
			BaseEdit: tgBotAPI.BaseEdit{
				ChatID:    update.CallbackQuery.Message.Chat.ID,
				MessageID: update.CallbackQuery.Message.MessageID,
			},
			Text:      txt,
			ParseMode: "HTML",
		}

		user.StageData = &model.EmptyStageData
		err := bot.redis.SetUserCache(ctx, user)
		if err != nil {
			log.Error(err)
			return
		}

		_, err = bot.api.Request(msg)
		if err != nil {
			log.Error(err)
		}

		cb := tgBotAPI.NewCallback(update.CallbackQuery.ID, "")
		_, err = bot.api.Request(cb)
		if err != nil {
			log.Error(err)
		}
	} else {
		user.StageData = &model.EmptyStageData
		err := bot.redis.SetUserCache(ctx, user)
		if err != nil {
			log.Error(err)
			return
		}
	}

	bot.balance(ctx, update, user)
}

func (bot *Bot) Notify(ctx context.Context, tx *tlb.Transaction) {
	if tx.IO.In.MsgType == "INTERNAL" {
		address := tx.IO.In.AsInternal().DstAddr.String()
		addrList := bot.storage.GetInMemoryWallets()
		userID, exist := addrList[address]

		if exist {
			from := tx.IO.In.AsInternal().SrcAddr.String()
			comment := tx.IO.In.AsInternal().Comment()
			amount := tx.IO.In.AsInternal().Amount.TON()

			bal, err := bot.ton.GetWalletBalance(ctx, address)
			if err != nil {
				log.Error(err)
				return
			}

			txt := fmt.Sprintf(ReceivedCoins, amount, hex.EncodeToString(tx.Hash), pkg.ShortAddr(from))
			if comment != "" {
				txt += fmt.Sprintf(ReceivedComment, comment)
			}
			txt += fmt.Sprintf(ReceivedBalance, bal)

			if err := bot.sendNotification(userID, txt, nil); err != nil {
				log.Error(err)
			}
		}
	}
}

package bot

const (
	WelcomeNewUser     = "Hello, %s!\n\nYour new fast and secure wallet has been created. Here is its address:\n\n" + WelcomeCommon
	WelcomeExistedUser = "Welcome again, %s!\n\nYour wallet address:\n\n" + WelcomeCommon
	WelcomeCommon      = "<code>%s</code>\n\nJust click on it to copy. Or use the QR code to receive TON coins.\n\nYou can manage your wallet using following command:\n/balance - wallet balance\n/send - send coins\n/receive - receive coins\n/cancel - cancel current operation\n\nPlease remember that for security reasons, access to this wallet is only possible from your current Telegram account."

	Balance = "Balance: %s TON"

	ReceiveInstruction = "Show this QR code to receive TON coins or use this wallet address:\n\n<code>%s</code>"
	ReceivedCoins      = "+ %s TON\nFrom <a href=\"https://tonapi.io/transaction/%s\">%s</a>"
	ReceivedComment    = "\nComment: %s"
	ReceivedBalance    = "\nBalance: %s TON"

	ReceiveButton    = "Receive"
	SendButton       = "Send"
	BalanceButton    = "Balance"
	SendAllButton    = "Send all coins"
	AddCommentButton = "Add comment"
	ConfirmButton    = "Confirm"
	CancelButton     = "Cancel"

	AskWallet     = "OK, send me the recipient's address or take a picture of his QR code"
	InvalidQR     = "🤷 Unable to recognize the QR code, try another photo or send me recipient's address"
	InvalidWallet = "🤷️ This address is not valid, check it and try again"

	AskAmount      = "How many TON coins to send?"
	NoFunds        = "🤷️ You have no TON coins"
	NotEnoughFunds = "🤷️ Not enough funds on your balance, you just have %s TON coins. Also blockchain fee of ~ %s coins is charged for each transfer."
	InvalidAmount  = "️Only digits and one dot are allowed, try again, please"

	AskComment = "OK, send me your transfer comment"
	Comment    = "\n\nComment: %s"

	SendingConfirmation = "✋ Confirm sending?\n\nTo address: <code>%s</code>\n\nAmount: %s TON coins\nFee: ~ %s"
	SendingCoins        = "\n\n⏳ <i>Sending coins...</i>"
	Sent                = "\n\n✅ Successfully sent"
	Canceled            = "\n\nCanceled"
)

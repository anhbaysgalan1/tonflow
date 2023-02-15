package bot

const (
	WelcomeNewUser     = "Hello, %s!\n\nYour new fast and secure wallet has been created. " + WelcomeCommon
	WelcomeExistedUser = "Welcome again, %s!\n\n" + WelcomeCommon
	WelcomeCommon      = "Your wallet address is\n<code>%s</code>\n\nJust click on it to copy. The QR code of this address in the picture. You can manage your wallet using the buttons under this message.\n\n" + "Please remember that for security reasons, access to this wallet is only possible from your current Telegram account."
	ReceiveInstruction = "Show QR code to receive TON or use wallet address"

	Balance = "💎 Balance: %s TON"

	ReceiveButton       = "Receive"
	SendButton          = "Send"
	BalanceButton       = "Balance"
	UpdateBalanceButton = "Update balance"
	SendAllButton       = "Send all"
	SendMaxButton       = "Send maximum"
	AddCommentButton    = "Add comment"
	ConfirmButton       = "Confirm"
	CancelButton        = "Cancel"

	BalanceUpToDate = "Balance is up to date"
	BalanceUpdated  = "Balance updated"

	AskAmount      = "💰 How many TON to send?"
	NoFunds        = "🤷️ You have no TON"
	NotEnoughFunds = "🤷️ Not enough funds on your balance, you just have %s TON. Also blockchain fee of ~%s TON is charged for each transfer."
	InvalidAmount  = "☝️ Only digits and one dot are allowed, try again"

	AskWallet     = "📲 Enter the recipient's wallet address or take a photo of QR code"
	InvalidQR     = "🤷 Unable to recognize QR code, try another photo or enter wallet address"
	InvalidWallet = "🤷 There is no wallet with this address"

	AskComment = "💬 Ok, send your transfer comment"
	Comment    = "\n\nComment: %s"

	SendingConfirmation = "✋ Confirm sending?\n\nTo wallet address: <code>%s</code>\n\nAmount: %s TON\nFee: ~%s TON"
	Confirmed           = "\n\n⏳ Sending coins..."
	Sent                = "\n\n✅ Successfully sent"
	Canceled            = "\n\n✔ Canceled"

	ReceivedCoins  = "🎉 Received %s TON\nFrom: %s"
	SeeTransaction = "\n<a href=\"https://tonapi.io/transaction/%s\">see transaction</a>"
)

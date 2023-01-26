package bot

const (
	WelcomeNewUser      = "Hello, %s!\n\nYour new fast and secure wallet has been created. " + WelcomeCommon
	WelcomeExistedUser  = "Welcome again, %s!\n\n" + WelcomeCommon
	WelcomeCommon       = "Your wallet address is\n<code>%s</code>\n\nJust click on it to copy. The QR code of this address in the picture. You can manage your wallet using the buttons under this message.\n\n" + "Please remember that for security reasons, access to this wallet is only possible from your current Telegram account."
	Balance             = "💎 Balance %s TON"
	ReceiveInstruction  = "Show QR code to receive TON or use wallet address"
	AskAmount           = "💰 How many TON to send?"
	AskWallet           = "📲 Enter the recipient's wallet address or take a photo of QR code"
	SendingConfirmation = "✋ Confirm sending?\nTo wallet address: <code>%s</code>\nAmount: %s TON"
	NoFunds             = "🤷️ You have no TON"
	NotEnoughFunds      = "🤷️ Not enough funds on your balance, you just have %s TON"
	InvalidAmount       = "☝️ Only digits and one dot are allowed, try again"
	InvalidQR           = "🤷 Unable to recognize QR code, try another photo or enter wallet address"
	InvalidWallet       = "🤷 There is no wallet with this address"
	Canceled            = "✅ Canceled"
	ReceiveButton       = "📥 Receive"
	SendButton          = "📤 Send"
	BalanceButton       = "Balance"
	UpdateBalanceButton = "Update balance"
	AddCommentButton    = "💬 Add comment"
	ConfirmButton       = "Confirm"
	CancelButton        = "Cancel"
	BalanceUpToDate     = "Balance is up to date"
	BalanceUpdated      = "Balance updated"
)

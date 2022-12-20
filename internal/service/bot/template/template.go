package template

const (
	WelcomeNewUser     = "Hello, %s!\n\nYour new fast and secure wallet has been created. " + WelcomeCommon
	WelcomeExistedUser = "Welcome again, %s!\n\n" + WelcomeCommon
	WelcomeCommon      = "Your wallet address is\n<pre>%s</pre>\n\nJust click on it to copy. The QR code of this address in the picture. You can manage your wallet using the buttons at the bottom of the screen.\n\n" +
		"Please remember that for security reasons, access to this wallet is only possible from your current Telegram account."

	Balance = "💎 Your balance is %s TON"

	ReceiveInstruction  = "Show QR code to receive TON or use wallet address"
	AskAmount           = "💰 How many TON to send?"
	AskWallet           = "📲 Enter the recipient's wallet address or take a photo of QR code"
	SendingConfirmation = "✋ Confirm if you want to send <b>27</b> TON to wallet %s"

	NoFunds        = "🤷️ You have no TON"
	NotEnoughFunds = "🤷️ Not enough funds on your balance, you only have %s TON"
	InvalidAmount  = "☝️ Only digits and one dot are allowed, try again"
	InvalidQR      = "🤷 Unable to recognize QR code, try another photo or enter wallet address"
	InvalidWallet  = "🤷 There is no wallet with this address"
	Canceled       = "✖️ Canceled"

	ReceiveButton = "📥 Receive"
	SendButton    = "📤 Send"
	BalanceButton = "💎 Balance"
	ConfirmButton = "✅ Confirm"
	CancelButton  = "✖️ Cancel"
)

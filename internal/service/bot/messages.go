package bot

const (
	startCommonPart = "Please remember that for security reasons, access to this wallet is only possible from your current Telegram account.\n\n" +
		"You can manage your wallet using the buttons at the bottom of the screen." +
		"In addition, the following commands are available: /history - to see recent transactions and /news - to get latest news"
	startNewUser = "Hello, %s! Your new fast and secure wallet has been created. " +
		"The address is listed below. " + startCommonPart
	startRegisteredUser = "Welcome once again, %s! Your wallet address is listed below. " + startCommonPart
)

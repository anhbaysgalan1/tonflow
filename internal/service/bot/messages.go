package bot

const (
	startNewUser = "Hello, %s!\n\nYour new fast and secure wallet has been created. " + startCommonPart

	startRegisteredUser = "Welcome again, %s!\n\n" + startCommonPart

	startCommonPart = "Your wallet address is\n<pre>%s</pre>\n\nJust click on it to copy. The QR code of this address in the picture. You can manage your wallet using the buttons at the bottom of the screen.\n\n" +
		"Please remember that for security reasons, access to this wallet is only possible from your current Telegram account."
)

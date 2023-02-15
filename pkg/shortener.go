package pkg

func ShortAddress(a string) string {
	firstFour := a[:4]
	lastFour := a[len(a)-4:]
	return firstFour + "..." + lastFour
}

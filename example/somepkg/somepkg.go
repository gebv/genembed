package somepkg

//go:generate $EMBEDBIN somefile

func Value() string {
	return string(getLocalFile("somefile"))
}

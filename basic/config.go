package basic

var Config = make(map[string]string)

func SetConfig(option, value string) {
	Config[option] = value
}

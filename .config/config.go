package config

type ConfigType struct {
	URL   string
	Token string
}

var Config = ConfigType{
	URL:   "{{PLACEHOLDER1}}"
	Token: "{{PLACEHOLDER2}}"
}

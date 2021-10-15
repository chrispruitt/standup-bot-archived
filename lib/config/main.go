package config

import (
	"os"
	"strconv"
)

var (
	SlackBotToken    string
	SlackAppToken    string
	BotName          string
	ShellModeChannel string
	ShellMode        bool
	Debug            bool
)

func init() {
	SlackBotToken = getenv("SLACK_BOT_TOKEN", "").(string)
	SlackAppToken = getenv("SLACK_APP_TOKEN", "").(string)
	BotName = getenv("BOT_NAME", "slackbot").(string)
	ShellMode = getenv("SHELL_MODE", false).(bool)
	ShellModeChannel = getenv("SHELL_MODE_CHANNEL", "").(string)
	Debug = getenv("DEBUG", true).(bool)
}

func getenv(key string, fallback interface{}) interface{} {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	switch fallback.(type) {
	case string:
		v := os.Getenv(key)
		if len(value) == 0 {
			return fallback
		}
		return v
	case int:
		s := os.Getenv(key)
		v, err := strconv.Atoi(s)
		if err != nil {
			return fallback
		}
		return v

	case bool:
		s := os.Getenv(key)
		v, err := strconv.ParseBool(s)
		if err != nil {
			return fallback
		}
		return v
	default:
		return value
	}
}

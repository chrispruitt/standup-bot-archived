module github.com/chrispruitt/standup-bot

go 1.16

replace github.com/chrispruitt/standup-bot/lib/standup => ./standup

require (
	github.com/chrispruitt/go-slackbot v0.3.2
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/slack-go/slack v0.9.4
)

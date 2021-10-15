package standup

import (
	"fmt"
	"log"
	"os"

	"github.com/chrispruitt/go-slackbot/lib/bot"
	"github.com/chrispruitt/standup-bot/lib/config"
	"github.com/chrispruitt/standup-bot/lib/views"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

var (
	webApi     *slack.Client
	socketMode *socketmode.Client
)

func Start() {
	listner()
}

func listner() {

	// bot

	webApi = slack.New(
		config.SlackBotToken,
		slack.OptionAppLevelToken(config.SlackAppToken),
		slack.OptionDebug(config.Debug),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)
	socketMode = socketmode.New(
		webApi,
		socketmode.OptionDebug(config.Debug),
		socketmode.OptionLog(log.New(os.Stdout, "sm: ", log.Lshortfile|log.LstdFlags)),
	)
	authTest, authTestErr := webApi.AuthTest()
	if authTestErr != nil {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN is invalid: %v\n", authTestErr)
		os.Exit(1)
	}
	selfUserId := authTest.UserID

	go func() {
		for envelope := range socketMode.Events {
			switch envelope.Type {
			case socketmode.EventTypeEventsAPI:
				// Events API:

				// Acknowledge the eventPayload first
				socketMode.Ack(*envelope.Request)

				eventPayload, _ := envelope.Data.(slackevents.EventsAPIEvent)
				switch eventPayload.Type {
				case slackevents.CallbackEvent:
					switch event := eventPayload.InnerEvent.Data.(type) {
					case *slackevents.MessageEvent:
						if event.User != selfUserId {
							bot.HandleMessageEvent(event)
						}
					case *slackevents.AppMentionEvent:
						socketMode.Debugf("Skipped: %v", event)
					default:
						socketMode.Debugf("Skipped: %v", event)
					}
				default:
					socketMode.Debugf("unsupported Events API eventPayload received")
				}
			case socketmode.EventTypeSlashCommand:
				socketMode.Ack(*envelope.Request)
				cmd, _ := envelope.Data.(slack.SlashCommand)
				socketMode.Debugf("Slash command received: %+v", cmd)

				handleSlashCommand(cmd)
			case socketmode.EventTypeInteractive:
				socketMode.Ack(*envelope.Request)
				payload, _ := envelope.Data.(slack.InteractionCallback)

				switch payload.Type {
				case slack.InteractionTypeViewSubmission:
					switch payload.View.CallbackID {
					case views.SettingsModalCallBackId:
						submitSettingsModal(payload)
					default:
						socketMode.Debugf("Ignore Submission with CallbackID: %v", payload.View.CallbackID)
					}
				default:
					socketMode.Debugf("Ignore Payload Type: %v", payload.Type)
				}
			default:
				socketMode.Debugf("Skipped: %v", envelope.Type)
			}
		}
	}()

	socketMode.Run()
}

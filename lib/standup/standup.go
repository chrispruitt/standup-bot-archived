package standup

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/chrispruitt/go-slackbot/lib/bot"
	"github.com/chrispruitt/standup-bot/lib/types"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/slack-go/slack"
)

var brain = Brain{Standups: make(map[string]types.StandupSettings)}
var header = "Asynchronous Standups! Less time in meetings means more time getting things done. Keep the channel clean by using threads! To manage times, participants, etc, just type the `/standup` command."

type Brain struct {
	Standups map[string]types.StandupSettings `json:"standups"`
}

func init() {
	// TODO sync brain to/from s3 file - use bot.RegisterPeriodicScript

	bot.RegisterScript(bot.Script{
		Name:    "test",
		Matcher: "test",
		Function: func(context *bot.ScriptContext) {
			// declare settings within function scope
			users, err := bot.SlackClient.GetUsersInfo("U01TDLJ5AG3")

			fmt.Printf("USERS: %v\n", users)
			fmt.Printf("ERR: %v\n", err)
		},
	})
}

func RegisterStandup(settings types.StandupSettings) error {

	fmt.Printf("\n\nRegistering Standup \n %v\n\n", settings)

	err := bot.RegisterPeriodicScript(bot.PeriodicScript{
		Name:     fmt.Sprintf("standup-solicit-%s", settings.ChannelID),
		CronSpec: settings.SolicitCronSpec,
		Function: getSolicitStandupFunc(settings),
	})

	// bot.RegisterScript(bot.Script{
	// 	Name:        fmt.Sprintf("standup-solicit-%s", settings.ChannelID),
	// 	Matcher:     bot.Matcher(fmt.Sprintf("standup-solicit-%s", settings.ChannelID)),
	// 	Description: "Test standup bot for a given channel",
	// 	Function: func(context *bot.ScriptContext) {
	// 		// declare settings within function scope
	// 		settings := settings
	// 		fmt.Println("standup solicit")

	// 		for _, userId := range settings.Participants {
	// 			bot.PostMessage(userId, settings.SolicitMsg)
	// 		}
	// 	},
	// })

	if err != nil {
		return err
	}

	bot.RegisterPeriodicScript(bot.PeriodicScript{
		Name:     fmt.Sprintf("standup-share-%s", settings.ChannelID),
		CronSpec: settings.ShareCronSpec,
		Function: getShareStandupFunc(settings),
	})

	brain.Standups[settings.ChannelID] = settings

	return nil
}

func getSolicitStandupFunc(settings types.StandupSettings) func() {
	return func() {
		// declare settings within function scope
		settings := settings
		for _, userId := range settings.Participants {
			bot.PostMessage(userId, settings.SolicitMsg)
		}
	}
}

func getShareStandupFunc(settings types.StandupSettings) func() {
	return func() {
		// declare settings within function scope
		settings := settings

		bot.PostMessage(settings.ChannelID, header)

		users, err := bot.SlackClient.GetUsersInfo(settings.Participants...)

		fmt.Printf("USERS: %v\n", users)
		fmt.Printf("ERR: %v\n", err)

		conversations, err := getConversations()
		if err != nil {
			fmt.Println("Error getting conversations: ", err)
		}

		// seed our random colors
		rand.Seed(time.Now().UTC().UnixNano())

		for _, user := range *users {

			if channel, ok := conversations[user.ID]; ok {
				fmt.Println(channel.ID, user.Name)
				notes, err := getStandupNote(channel.ID, settings.SolicitMsg)
				if err != nil {
					fmt.Println(err)
					continue
				}

				// Shame people in the channel that aren't participating
				// TODO
				// if len(notes) == 0 {
				// 	if shameParticipants == "true" {
				// 		notes = append(notes, fmt.Sprintf(":poop: %s has no standup notes", user.Name))
				// 	} else {
				// 		continue
				// 	}
				// }

				attachments := []slack.Attachment{
					{
						Title: fmt.Sprintf("%s standup notes:", user.RealName),
						Text:  fmt.Sprintf(strings.Join(reverse(notes), "\n")),
						Color: colorful.FastHappyColor().Hex(),
					},
				}

				_, _, err = bot.SlackClient.PostMessage(
					settings.ChannelID,
					slack.MsgOptionText("", false),
					slack.MsgOptionEnableLinkUnfurl(),
					slack.MsgOptionAttachments(attachments...))

				if err != nil {
					fmt.Printf("Error posting standup report: %s\n", err)
				}

			} else {
				fmt.Printf("channel %s is orphaned", channel.ID)
			}
		}

		fmt.Println("standup share")
	}
}

func getConversations() (map[string]slack.Channel, error) {
	conversations := make(map[string]slack.Channel)
	r, _, err := bot.SlackClient.GetConversations(&slack.GetConversationsParameters{
		ExcludeArchived: true,
		Types:           []string{"im"},
	})
	if err != nil {
		return nil, err
	}
	for _, channel := range r {
		conversations[channel.User] = channel
	}
	return conversations, nil
}

func getStandupNote(channelID string, solicitStandupMessage string) ([]string, error) {
	conversationHistory, err := bot.SlackClient.GetConversationHistory(&slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		// Oldest:    fmt.Sprintf("%d.000001", (now.Add(time.Hour * time.Duration(-1))).Unix()),
	})
	if err != nil {
		return nil, fmt.Errorf("error getting standup info from channel %s: %s", channelID, err)
	}

	var txt []string

	for _, m := range conversationHistory.Messages {

		if m.Text == solicitStandupMessage {
			break
		}
		txt = append(txt, m.Text)
	}

	return txt, nil
}

func reverse(ss []string) []string {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
	return ss
}

// for _, mType := range []string{"solicit", "standup"} {

// 	cmdType := mType
// 	err := bot.RegisterPeriodicScript(bot.PeriodicScript{
// 		Name:     fmt.Sprintf("standup-%s-%s", cmdType, settings.ChannelID),
// 		CronSpec: "*/1 * * * *",
// 		Function: func() {
// 			// TODO Solicit standup
// 			fmt.Printf("%s standup solicit \n", cmdType)
// 		},
// 	})
// 	if err != nil {
// 		return err
// 	}
// }

package standup

import (
	"fmt"
	"log"
	"strings"

	"github.com/chrispruitt/standup-bot/lib/types"
	"github.com/chrispruitt/standup-bot/lib/views"
	"github.com/slack-go/slack"
)

func handleSlashCommand(cmd slack.SlashCommand) {
	switch cmd.Text {
	case "":
		openSettingsModal(cmd)
	case "solicit":
		settings := types.NewStandupSettings(cmd.ChannelID, cmd.ChannelName)
		settings.Participants = []string{cmd.UserID}
		getSolicitStandupFunc(*settings)()
	case "share":
		settings := types.NewStandupSettings(cmd.ChannelID, cmd.ChannelName)
		settings.Participants = []string{cmd.UserID}
		getShareStandupFunc(*settings)()
	default:
		socketMode.Debugf("Ignored: %v", cmd)
	}
}

func openSettingsModal(cmd slack.SlashCommand) {

	initialSettings := *types.NewStandupSettings(cmd.ChannelID, cmd.ChannelName)

	// If config for channel exists - populate form with saved settings
	if value, ok := brain.Standups[cmd.ChannelID]; ok {
		initialSettings = value
	}

	resp, err := webApi.OpenView(cmd.TriggerID, views.GetSettingsModal(initialSettings))
	if err != nil {
		log.Printf("Failed to opemn a modal: %v", err)
	}
	socketMode.Debugf("views.open response: %v", resp)
}

func submitSettingsModal(payload slack.InteractionCallback) {
	socketMode.Debugf("submit settings modal payload.View.State.Values: %v", payload.View.State.Values)

	channelID := payload.View.State.Values[views.ModalStandupChannelBlockId][views.ModalStandupChannelActionId].SelectedChannel
	participants := payload.View.State.Values[views.ModalParticipantsBlockId][views.ModalParticipantsActionId].SelectedUsers
	solicitTime := strings.Split(payload.View.State.Values[views.ModalSolicitTimeBlockId][views.ModalSolicitTimeActionId].SelectedTime, ":")
	shareTime := strings.Split(payload.View.State.Values[views.ModalShareTimeBlockId][views.ModalShareTimeActionId].SelectedTime, ":")
	meetingDays := payload.View.State.Values[views.ModalMeetingDaysBlockId][views.ModalMeetingDaysActionId].SelectedOptions
	solicitMsg := payload.View.State.Values[views.ModalSolicitMessgaeBlockId][views.ModalSolicitMessageActionId].Value

	settings := types.StandupSettings{
		ChannelID:       channelID,
		SolicitCronSpec: fmt.Sprintf("%s %s * * %s", solicitTime[1], solicitTime[0], meetingDaysToCron(meetingDays)),
		ShareCronSpec:   fmt.Sprintf("%s %s * * %s", shareTime[1], shareTime[0], meetingDaysToCron(meetingDays)),
		SolicitMsg:      solicitMsg,
		Participants:    participants,
	}

	socketMode.Debugf("standup submission settings: %v", settings)

	RegisterStandup(settings)
}

func meetingDaysToCron(meetingDays []slack.OptionBlockObject) string {
	values := []string{}
	for _, day := range meetingDays {
		values = append(values, day.Value)
	}
	return strings.Join(values, ",")
}

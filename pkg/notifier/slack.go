package notifier

import (
	"github.com/nlopes/slack"
	"github.com/rs/zerolog/log"
)

func (n *Notifier) sendNotification(a slack.Attachment) {

	msgOpts := slack.PostMessageParameters{Attachments: []slack.Attachment{}}
	msgOpts.Attachments = append(msgOpts.Attachments, a)

	id, _, err := n.slack.PostMessage(n.config.slack.Channel, "", msgOpts)
	if err != nil {
		log.Error().Err(err).Str("channel", id).Msg("failed to send Slack notification")
		return
	}

	log.Info().Str("channel", id).Msg("successfully sent Slack notification")
}

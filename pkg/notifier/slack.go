package notifier

import (
	"github.com/nlopes/slack"
	"github.com/rs/zerolog/log"
)

func (n *Notifier) sendNotification(nomadID string, msg slack.Attachment) {
	if val, ok := n.state.notifications[nomadID]; !ok {
		n.sendNewMsg(nomadID, msg)
	} else {
		n.sendUpdateMsg(nomadID, val.timestamp, msg)
	}
}

func (n *Notifier) sendNewMsg(nomadID string, msg slack.Attachment) {
	msgOpts := slack.PostMessageParameters{Attachments: []slack.Attachment{}}
	msgOpts.Attachments = append(msgOpts.Attachments, msg)

	chanID, ts, err := n.slack.PostMessage(n.config.slack.Channel, "", msgOpts)
	if err != nil {
		log.Error().Err(err).Str("channel", chanID).Msg("failed to send new Slack notification")
		return
	}
	if n.chanID == nil {
		n.chanID = &chanID
	}

	n.newNotifierState(nomadID, ts, msg)

	log.Info().Str("channel", chanID).Msg("successfully sent new Slack notification")
}

func (n *Notifier) sendUpdateMsg(id, ts string, msg slack.Attachment) {
	u := slack.MsgOptionUpdate(ts)
	attachOpts := slack.MsgOptionAttachments(msg)
	opts := []slack.MsgOption{u, attachOpts}

	chanID, ts, _, err := n.slack.SendMessage(*n.chanID, opts...)
	if err != nil {
		log.Error().Err(err).Str("channel", chanID).Msg("failed to send update Slack notification")
	}

	n.updateNotifierState(id, ts, msg)

	log.Info().Str("channel", chanID).Msg("successfully sent new Slack notification")
}

func (n *Notifier) newNotifierState(id, ts string, msg slack.Attachment) {
	log.Debug().
		Str("ts", ts).
		Str("id", id).
		Msg("adding new entry into notifier state tracking")

	m := []slack.Attachment{msg}
	n.state.Lock()
	n.state.notifications[id] = notification{ts, m}
	n.state.Unlock()
}

func (n *Notifier) updateNotifierState(id, ts string, msg slack.Attachment) {
	log.Debug().
		Str("ts", ts).
		Str("id", id).
		Msg("adding update into notifier state tracking")

	n.state.Lock()
	s := n.state.notifications[id]
	s.timestamp = ts
	s.messages = append(s.messages, msg)
	n.state.Unlock()
}

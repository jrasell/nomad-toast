package notifier

import (
	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/nomad-toast/pkg/config"
	"github.com/nlopes/slack"
	"github.com/rs/zerolog/log"
)

// Notifier is the main notifier struct which holds all config items for triggering notifications.
type Notifier struct {
	config *notifierConfig
	slack  *slack.Client

	MsgChan chan interface{}
}

type notifierConfig struct {
	slack *config.SlackConfig
}

// NewNotifier builds a new notifier struct in order to run the nomad-toast notifier task.
// The message channel is importantly where events are received from the watcher.
func NewNotifier(cfg *config.SlackConfig) (*Notifier, error) {
	return &Notifier{
		config:  &notifierConfig{slack: cfg},
		MsgChan: make(chan interface{}),
		slack:   slack.New(cfg.AuthToken),
	}, nil
}

// Run triggers the notifier to start listening for messages.
func (n *Notifier) Run() {
	log.Info().Msg("starting deployment notifier")

	for {
		select {
		case msg := <-n.MsgChan:

			switch v := msg.(type) {
			case *api.AllocationListStub:
				go n.formatAllocationMessage(v)
			case *api.Deployment:
				go n.formatDeploymentMessage(v)
			default:
				log.Error().Msg("notifier received unknown message type from watcher")
			}
		}
	}
}

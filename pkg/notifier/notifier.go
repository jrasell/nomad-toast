package notifier

import (
	"sync"

	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/nomad-toast/pkg/config"
	"github.com/jrasell/nomad-toast/pkg/watcher"
	"github.com/nlopes/slack"
	"github.com/rs/zerolog/log"
)

// Notifier is the main notifier struct which holds all config items for triggering notifications.
type Notifier struct {
	config *notifierConfig
	slack  *slack.Client
	state  *notifications
	chanID *string

	MsgChan chan interface{}

	nomadRegion string
}

type notifierConfig struct {
	slack *config.SlackConfig
	ui    *config.UI
}

type notifications struct {
	sync.RWMutex
	notifications map[string]notification
	nomadType     watcher.EndpointType
}

type notification struct {
	timestamp string
	messages  []slack.Attachment
}

// NewNotifier builds a new notifier struct in order to run the nomad-toast notifier task.
// The message channel is importantly where events are received from the watcher.
func NewNotifier(sCfg *config.SlackConfig, fCfg *config.UI, et watcher.EndpointType) (*Notifier, error) {
	return &Notifier{
		config:  &notifierConfig{slack: sCfg, ui: fCfg},
		MsgChan: make(chan interface{}),
		slack:   slack.New(sCfg.AuthToken),
		state:   &notifications{notifications: make(map[string]notification), nomadType: et},
	}, nil
}

// Run triggers the notifier to start listening for messages.
func (n *Notifier) Run(region string) {
	log.Info().Msgf("starting %s notifier", n.state.nomadType)

	n.nomadRegion = region

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

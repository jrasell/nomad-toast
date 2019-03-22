package watcher

import (
	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/nomad-toast/pkg/config"
	"github.com/rs/zerolog/log"
)

// Watcher is the main watcher struct which holds all config items for running a watcher.
type Watcher struct {
	config *watcherConfig
	EndpointType
	lastChangeIndex uint64
	initialRun      bool
	nomadClient     *api.Client

	mshChan chan interface{}
}

type watcherConfig struct {
	nomad *config.NomadConfig
}

// EndpointType is the Nomad API endpoint type support for watching.
type EndpointType string

const (
	// Allocations represents the Nomad Allocations API endpoint.
	Allocations EndpointType = "allocations"

	// Deployments represents the Nomad Deployments API endpoint.
	Deployments EndpointType = "deployments"
)

// NewWatcher builds a new watcher struct in order to run the nomad-toast watcher task.
// The message channel is importantly where events are sent for notification.
func NewWatcher(cfg *config.NomadConfig, et EndpointType, mChan chan interface{}) (*Watcher, string, error) {

	d := api.DefaultConfig()
	d.Address = cfg.NomadAddress

	if cfg.NomadRegion != "" {
		d.Region = cfg.NomadRegion
	}

	client, err := api.NewClient(d)
	if err != nil {
		return nil, "", err
	}

	var r string

	if cfg.NomadRegion == "" {
		r, err = getNomadRegion(client)
		if err != nil {
			return nil, "", err
		}
		cfg.NomadRegion = r
	}

	return &Watcher{
		config:          &watcherConfig{nomad: cfg},
		EndpointType:    et,
		initialRun:      true,
		lastChangeIndex: 1,
		nomadClient:     client,
		mshChan:         mChan,
	}, r, nil
}

func getNomadRegion(c *api.Client) (string, error) {
	info, err := c.Agent().Self()
	if err != nil {
		return "", err
	}
	return info.Config["Region"].(string), nil
}

// Run triggers the watcher of configured type to be run.
// This will monitor the endpoint for changes and pass them to the notifier for notifications.
func (w *Watcher) Run() {
	switch w.EndpointType {
	case Allocations:
		go w.runAllocationWatcher()
	case Deployments:
		go w.runDeploymentWatcher()
	default:
		log.Fatal().Msgf("unknown endpoint type %v", w.EndpointType)
	}
}

func (w *Watcher) indexHasChange(new, old uint64) bool {
	if new < old {
		return false
	}
	return true
}

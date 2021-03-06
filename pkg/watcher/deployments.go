package watcher

import (
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog/log"
)

func (w *Watcher) runDeploymentWatcher() {
	log.Info().Msg("starting deployment watcher")

	q := &api.QueryOptions{
		AllowStale: w.config.nomad.AllowStale,
		WaitIndex:  uint64(w.lastChangeIndex),
		WaitTime:   5 * time.Minute,
	}

	for {

		var maxFound uint64

		deployments, meta, err := w.nomadClient.Deployments().List(q)
		// If Nomad returned an error, log this and sleep.
		// The sleep is needed, otherwise we will just call the failing endpoint over and over.
		if err != nil {
			log.Error().Err(err).Msg("unable to fetch deployments list from Nomad")
			time.Sleep(10 * time.Second)
			continue
		}

		if w.initialRun {
			log.Info().Msg("running initial deployment watcher index update")
			q.WaitIndex = meta.LastIndex
			w.lastChangeIndex = meta.LastIndex
			w.initialRun = false
			continue
		}

		if !w.indexHasChange(meta.LastIndex, q.WaitIndex) {
			log.Debug().Msg("deployments index has not changed")
			continue
		}

		maxFound = meta.LastIndex
		log.Debug().Msg("deployments index has changed")

		for _, deployment := range deployments {

			if !w.indexHasChange(deployment.ModifyIndex, maxFound) {
				log.Debug().Str("job", deployment.JobID).Msg("job deployment index has not changed")
				continue
			}

			maxFound = deployment.ModifyIndex
			log.Info().Str("job", deployment.JobID).Msg("job deployment index has changed")

			w.mshChan <- deployment
		}
		q.WaitIndex = maxFound
		w.lastChangeIndex = maxFound
	}
}

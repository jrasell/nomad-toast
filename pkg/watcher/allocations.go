package watcher

import (
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog/log"
)

func (w *Watcher) runAllocationWatcher() {
	log.Info().Msg("starting allocation watcher")

	q := &api.QueryOptions{
		AllowStale: w.config.nomad.AllowStale,
		WaitIndex:  uint64(w.lastChangeIndex),
		WaitTime:   5 * time.Minute,
	}

	for {
		allocations, meta, err := w.nomadClient.Allocations().List(q)

		// If Nomad returned an error, log this and sleep.
		// The sleep is needed, otherwise we will just call the failing endpoint over and over.
		if err != nil {
			log.Error().Err(err).Msg("unable to fetch allocations list from Nomad")
			time.Sleep(10 * time.Second)
			continue
		}

		if w.initialRun {
			log.Info().Msg("running initial allocation watcher index update")
			q.WaitIndex = meta.LastIndex
			w.lastChangeIndex = meta.LastIndex
			w.initialRun = false
			continue
		}

		if !w.indexHasChange(q.WaitIndex, meta.LastIndex) {
			log.Debug().Msg("allocations index has not changed")
			continue
		}

		log.Debug().Msg("allocations index has changed")

		for _, alloc := range allocations {

			if !w.indexHasChange(w.lastChangeIndex, alloc.ModifyIndex) {
				log.Debug().Str("alloc-id", alloc.ID).Msg("allocation index has not changed")
				continue
			}

			log.Info().Str("alloc-id", alloc.ID).Msg("allocation index has changed")

			if len(alloc.TaskStates) == 0 {
				continue
			}

			w.mshChan <- alloc
		}
		q.WaitIndex = meta.LastIndex
		w.lastChangeIndex = meta.LastIndex
	}
}

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

		var maxFound uint64

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

		if !w.indexHasChange(meta.LastIndex, q.WaitIndex) {
			log.Debug().Msg("allocations index has not changed")
			continue
		}

		maxFound = meta.LastIndex
		log.Debug().Msg("allocations index has changed")

		for _, alloc := range allocations {

			if !w.indexHasChange(alloc.ModifyIndex, maxFound) {
				log.Debug().Str("alloc-id", alloc.ID).Msg("allocation index has not changed")
				continue
			}

			if len(alloc.TaskGroup) == 0 {
				continue
			}

			maxFound = alloc.ModifyIndex
			log.Info().Str("alloc-id", alloc.ID).Msg("allocation index has changed")

			w.mshChan <- alloc
		}
		q.WaitIndex = maxFound
		w.lastChangeIndex = maxFound
	}
}

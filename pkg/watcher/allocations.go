package watcher

import (
	"strings"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/jrasell/nomad-toast/pkg/config"
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
			// If the allocation's index hasn't changed, then there is nothing to notify about
			if !w.indexHasChange(alloc.ModifyIndex, maxFound) {
				log.Debug().Str("alloc-id", alloc.ID).Msg("allocation index has not changed")
				continue
			}

			// If the allocation index *has* changed, check if the change is of interest
			maxFound = alloc.ModifyIndex
			log.Info().Str("alloc-id", alloc.ID).Msg("allocation index has changed")
			if isFiltered(alloc, w.config.nomad.AllocCfg) {
				continue
			}

			w.mshChan <- alloc
		}
		q.WaitIndex = maxFound
		w.lastChangeIndex = maxFound
	}
}

// isFiltered checks whether a given update about a notification is to be filtered or not
func isFiltered(alloc *api.AllocationListStub, allocCfg *config.AllocConfig) bool {
	for i := range allocCfg.ExcludeStates {
		if strings.ToLower(alloc.ClientStatus) == strings.ToLower(allocCfg.ExcludeStates[i]) {
			log.Debug().Str("alloc-id", alloc.ID).Str("client-status", alloc.ClientStatus).Msg("allocation client status blacklisted, omitting")
			return true
		}
	}

	if allocCfg.IncludeStates != nil {
		for i := range allocCfg.IncludeStates {
			if strings.ToLower(alloc.ClientStatus) == strings.ToLower(allocCfg.IncludeStates[i]) {
				return false
			}
		}
		log.Debug().Str("alloc-id", alloc.ID).Str("client-status", alloc.ClientStatus).Msg("allocation client not whitelisted, omitting")
		return true
	}

	return false
}

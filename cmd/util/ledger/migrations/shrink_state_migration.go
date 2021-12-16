package migrations

import (
	lc "github.com/onflow/flow-go/cmd/util/ledger/common"
	"github.com/onflow/flow-go/fvm/state"
	"github.com/onflow/flow-go/ledger"
	"github.com/onflow/flow-go/model/flow"
	"github.com/rs/zerolog"
	"runtime"
	"sync"
)

// ShrinkStateMigration is a migration that shrinks the state, so that it can be more manageable to be used in tests.
type ShrinkStateMigration struct {
	Chain flow.Chain
	Log   zerolog.Logger
}

func (m *ShrinkStateMigration) Name() string {
	return "Shrink State Migration"
}

var _ ledger.Migration = &ShrinkStateMigration{}

func (m *ShrinkStateMigration) Migrate(payloads []ledger.Payload) ([]ledger.Payload, error) {

	l := NewView(payloads)
	st := state.NewState(l)
	sth := state.NewStateHolder(st)
	gen := state.NewStateBoundAddressGenerator(sth, m.Chain)

	progress := lc.NewProgressBar(m.Log, int64(gen.AddressCount()), "collecting account data:")

	workerCount := runtime.NumCPU()
	addressIndexes := make(chan uint64)
	filteredAddresses := make(chan flow.Address, workerCount)
	keepSet := make(map[flow.Address]struct{})
	wg := &sync.WaitGroup{}
	fwg := &sync.WaitGroup{}

	if workerCount == 0 {
		workerCount = 1
	}

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go m.accountFilter(wg, l, addressIndexes, filteredAddresses)
	}
	fwg.Add(1)
	go func() {
		for address := range filteredAddresses {
			keepSet[address] = struct{}{}
		}
		fwg.Done()
	}()

	for i := uint64(1); i <= gen.AddressCount(); i++ {
		addressIndexes <- i
		progress.Increment()
	}
	close(addressIndexes)
	wg.Wait()
	close(filteredAddresses)
	fwg.Wait()

	progress.Finish()

	progress = lc.NewProgressBar(m.Log, int64(len(keepSet)), "shrinking state:")

	newPayloads := make([]ledger.Payload, 0)

	for _, payload := range payloads {
		progress.Increment()

		id, err := KeyToRegisterID(payload.Key)
		if err != nil {
			return nil, err
		}
		if len(id.Owner) == 0 {
			newPayloads = append(newPayloads, payload)
			continue
		}
		a := flow.BytesToAddress([]byte(id.Owner))
		if _, ok := keepSet[a]; ok {
			newPayloads = append(newPayloads, payload)
		}
	}
	progress.Finish()

	return newPayloads, nil
}

func (m *ShrinkStateMigration) accountFilter(wg *sync.WaitGroup, view *view, indexes <-chan uint64, filteredAddresses chan<- flow.Address) {
	v := view.NewChild()
	st := state.NewState(v)
	sth := state.NewStateHolder(st)
	accounts := state.NewAccounts(sth)

	shouldTakeAccount := func(addressIndex uint64) (bool, flow.Address) {
		address, err := m.Chain.AddressAtIndex(addressIndex)
		if err != nil {
			m.Log.
				Err(err).
				Uint64("index", addressIndex).
				Msgf("Error getting address")
			return false, address
		}
		contracts, err := accounts.GetContractNames(address)
		if err != nil {
			m.Log.
				Err(err).
				Uint64("index", addressIndex).
				Msgf("Error getting address contracts")
			return false, address
		}

		// Condition 1. keep all accounts that have contracts
		if len(contracts) > 0 {
			return true, address
		}
		// more conditions go here
		return false, address
	}

	for index := range indexes {
		yes, address := shouldTakeAccount(index)
		if yes {
			filteredAddresses <- address
		}
	}
	wg.Done()
}

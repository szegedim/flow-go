package migrations

import (
	"github.com/onflow/flow-go/ledger"
)

// PruneMigration removes all the payloads with empty value
// this prunes the trie for values that has been deleted
type PruneMigration struct{}

func (m *PruneMigration) Name() string {
	return "Prune Migration"
}

var _ ledger.Migration = &PruneMigration{}

func (m *PruneMigration) Migrate(payloads []ledger.Payload) ([]ledger.Payload, error) {
	newPayload := make([]ledger.Payload, 0, len(payloads))
	for _, p := range payloads {
		if len(p.Value) > 0 {
			newPayload = append(newPayload, p)
		}
	}
	return newPayload, nil
}

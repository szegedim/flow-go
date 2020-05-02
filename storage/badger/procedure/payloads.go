// (c) 2019 Dapper Labs - ALL RIGHTS RESERVED

package procedure

import (
	"fmt"

	"github.com/dgraph-io/badger/v2"

	"github.com/dapperlabs/flow-go/model/flow"
	"github.com/dapperlabs/flow-go/storage/badger/operation"
)

func InsertPayload(payload *flow.Payload) func(*badger.Txn) error {
	return func(tx *badger.Txn) error {

		// insert the block guarantees
		for _, guarantee := range payload.Guarantees {
			err := operation.SkipDuplicates(operation.InsertGuarantee(guarantee))(tx)
			if err != nil {
				return fmt.Errorf("could not insert guarantee (%x): %w", guarantee.CollectionID, err)
			}
		}

		// insert the block seals
		for _, seal := range payload.Seals {
			err := operation.SkipDuplicates(operation.InsertSeal(seal))(tx)
			if err != nil {
				return fmt.Errorf("could not insert seal (%x): %w", seal.ID(), err)
			}
		}

		return nil
	}
}

func IndexPayload(blockID flow.Identifier, payload *flow.Payload) func(*badger.Txn) error {
	return func(tx *badger.Txn) error {

		// index guarantees
		err := IndexGuarantees(blockID, flow.GetIDs(payload.Guarantees))(tx)
		if err != nil {
			return fmt.Errorf("could not index guarantees: %w", err)
		}

		// index seals
		err = IndexSeals(blockID, flow.GetIDs(payload.Seals))(tx)
		if err != nil {
			return fmt.Errorf("could not index seals: %w", err)
		}

		return nil
	}
}

func RetrievePayload(blockID flow.Identifier, payload *flow.Payload) func(tx *badger.Txn) error {
	return func(tx *badger.Txn) error {
		// make sure there is a nil value on error
		*payload = flow.Payload{}

		// get guarantees
		var guarantees []*flow.CollectionGuarantee
		err := RetrieveGuarantees(blockID, &guarantees)(tx)
		if err != nil {
			return fmt.Errorf("could not retrieve guarantees: %w", err)
		}

		// get seals
		var seals []*flow.Seal
		err = RetrieveSeals(blockID, &seals)(tx)
		if err != nil {
			return fmt.Errorf("could not retrieve seals: %w", err)
		}

		// create the block content
		*payload = flow.Payload{
			Identities: nil,
			Guarantees: guarantees,
			Seals:      seals,
		}

		return nil
	}
}

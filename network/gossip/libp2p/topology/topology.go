package topology

import "github.com/onflow/flow-go/model/flow"

// Topology provides a subset of nodes which a given node should directly connect to for 1-k messaging
type Topology interface {
	// Subset returns a random subset of the identity list that is passed. shouldHave represents set of
	// identities that should be included in the returned subset.
	Subset(ids flow.IdentityList, shouldHave flow.IdentityList, topic string, fanout uint) (flow.IdentityList, error)
}

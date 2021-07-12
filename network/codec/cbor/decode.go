// (c) 2019 Dapper Labs - ALL RIGHTS RESERVED

package cbor

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/pkg/errors"

	"github.com/onflow/flow-go/binstat"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/model/libp2p/message"
	"github.com/onflow/flow-go/model/messages"
)

func switchenv2v(code uint8) (interface{}, error) {
	var v interface{}

	switch code {

	// consensus
	case CodeBlockProposal:
		v = &messages.BlockProposal{}
	case CodeBlockVote:
		v = &messages.BlockVote{}

	// cluster consensus
	case CodeClusterBlockProposal:
		v = &messages.ClusterBlockProposal{}
	case CodeClusterBlockVote:
		v = &messages.ClusterBlockVote{}
	case CodeClusterBlockResponse:
		v = &messages.ClusterBlockResponse{}

	// protocol state sync
	case CodeSyncRequest:
		v = &messages.SyncRequest{}
	case CodeSyncResponse:
		v = &messages.SyncResponse{}
	case CodeRangeRequest:
		v = &messages.RangeRequest{}
	case CodeBatchRequest:
		v = &messages.BatchRequest{}
	case CodeBlockResponse:
		v = &messages.BlockResponse{}

	// collections, guarantees & transactions
	case CodeCollectionGuarantee:
		v = &flow.CollectionGuarantee{}
	case CodeTransactionBody:
		v = &flow.TransactionBody{}
	case CodeTransaction:
		v = &flow.Transaction{}

	// core messages for execution & verification
	case CodeExecutionReceipt:
		v = &flow.ExecutionReceipt{}
	case CodeResultApproval:
		v = &flow.ResultApproval{}

	// execution state synchronization
	case CodeExecutionStateSyncRequest:
		v = &messages.ExecutionStateSyncRequest{}
	case CodeExecutionStateDelta:
		v = &messages.ExecutionStateDelta{}

	// data exchange for execution of blocks
	case CodeChunkDataRequest:
		v = &messages.ChunkDataRequest{}
	case CodeChunkDataResponse:
		v = &messages.ChunkDataResponse{}

	case CodeApprovalRequest:
		v = &messages.ApprovalRequest{}
	case CodeApprovalResponse:
		v = &messages.ApprovalResponse{}

	// generic entity exchange engines
	case CodeEntityRequest:
		v = &messages.EntityRequest{}
	case CodeEntityResponse:
		v = &messages.EntityResponse{}

	// testing
	case CodeEcho:
		v = &message.TestMessage{}

	default:
		return nil, errors.Errorf("invalid message code (%d)", code)
	}

	return v, nil
}

func switchenv2what(code uint8) (string, error) {
	var what string

	switch code {

	// consensus
	case CodeBlockProposal:
		what = "CodeBlockProposal"
	case CodeBlockVote:
		what = "CodeBlockVote"

	// cluster consensus
	case CodeClusterBlockProposal:
		what = "CodeClusterBlockProposal"
	case CodeClusterBlockVote:
		what = "CodeClusterBlockVote"
	case CodeClusterBlockResponse:
		what = "CodeClusterBlockResponse"

	// protocol state sync
	case CodeSyncRequest:
		what = "CodeSyncRequest"
	case CodeSyncResponse:
		what = "CodeSyncResponse"
	case CodeRangeRequest:
		what = "CodeRangeRequest"
	case CodeBatchRequest:
		what = "CodeBatchRequest"
	case CodeBlockResponse:
		what = "CodeBlockResponse"

	// collections, guarantees & transactions
	case CodeCollectionGuarantee:
		what = "CodeCollectionGuarantee"
	case CodeTransactionBody:
		what = "CodeTransactionBody"
	case CodeTransaction:
		what = "CodeTransaction"

	// core messages for execution & verification
	case CodeExecutionReceipt:
		what = "CodeExecutionReceipt"
	case CodeResultApproval:
		what = "CodeResultApproval"

	// execution state synchronization
	case CodeExecutionStateSyncRequest:
		what = "CodeExecutionStateSyncRequest"
	case CodeExecutionStateDelta:
		what = "CodeExecutionStateDelta"

	// data exchange for execution of blocks
	case CodeChunkDataRequest:
		what = "CodeChunkDataRequest"
	case CodeChunkDataResponse:
		what = "CodeChunkDataResponse"

	case CodeApprovalRequest:
		what = "CodeApprovalRequest"
	case CodeApprovalResponse:
		what = "CodeApprovalResponse"

	// generic entity exchange engines
	case CodeEntityRequest:
		what = "CodeEntityRequest"
	case CodeEntityResponse:
		what = "CodeEntityResponse"

	// testing
	case CodeEcho:
		what = "CodeEcho"

	default:
		return "", errors.Errorf("invalid message code (%d)", code)
	}

	return what, nil
}

// decode will decode the envelope into an entity.
func env2vDecode(data []byte, code uint8, via string) (interface{}, error) {

	// create the desired message
	v, err1 := switchenv2v(code)
	what, err2 := switchenv2what(code)

	if nil != err1 {
		return nil, err1
	}

	if nil != err2 {
		return nil, err2
	}

	// unmarshal the payload
	p := binstat.EnterTimeVal(fmt.Sprintf("%s%s:%d", via, what, code), "", int64(len(data)))
	err := cbor.Unmarshal(data, v)
	binstat.Leave(p)
	if err != nil {
		return nil, fmt.Errorf("could not decode payload: %w", err)
	}

	return v, nil
}

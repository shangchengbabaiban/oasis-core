package staking

import (
	"encoding/hex"

	"github.com/cometbft/cometbft/abci/types"

	"github.com/oasisprotocol/oasis-core/go/common/crypto/signature"
	abciAPI "github.com/oasisprotocol/oasis-core/go/consensus/cometbft/api"
	registryState "github.com/oasisprotocol/oasis-core/go/consensus/cometbft/apps/registry/state"
	registry "github.com/oasisprotocol/oasis-core/go/registry/api"
)

func (app *Application) resolveEntityIDsFromVotes(
	ctx *abciAPI.Context,
	regState *registryState.MutableState,
	lastCommitInfo types.CommitInfo,
) ([]signature.PublicKey, error) {
	var entityIDs []signature.PublicKey
	for _, a := range lastCommitInfo.Votes {
		if !a.SignedLastBlock {
			continue
		}
		valAddr := a.Validator.Address

		// Map address to node/entity.
		node, err := regState.NodeByConsensusAddress(ctx, valAddr)
		switch err {
		case nil:
		case registry.ErrNoSuchNode:
			ctx.Logger().Warn("failed to get validator node",
				"err", err,
				"address", hex.EncodeToString(valAddr),
			)
			continue
		default:
			return nil, err
		}

		entityIDs = append(entityIDs, node.EntityID)
	}

	return entityIDs, nil
}

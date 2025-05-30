package churp

import (
	"fmt"

	"github.com/oasisprotocol/oasis-core/go/common/cbor"
	"github.com/oasisprotocol/oasis-core/go/consensus/api"
	"github.com/oasisprotocol/oasis-core/go/consensus/api/transaction"
	tmapi "github.com/oasisprotocol/oasis-core/go/consensus/cometbft/api"
	"github.com/oasisprotocol/oasis-core/go/keymanager/churp"
)

// Ensure that the CHURP extension implements the Extension interface.
var _ tmapi.Extension = (*churpExt)(nil)

type churpExt struct {
	appName string
	state   tmapi.ApplicationState
}

// New creates a new CHRUP extension for the key manager application.
func New(appName string, state tmapi.ApplicationState) tmapi.Extension {
	return &churpExt{
		appName: appName,
		state:   state,
	}
}

// Methods implements api.Extension.
func (ext *churpExt) Methods() []transaction.MethodName {
	return churp.Methods
}

// ExecuteTx implements api.Extension.
func (ext *churpExt) ExecuteTx(ctx *tmapi.Context, tx *transaction.Transaction) error {
	switch tx.Method {
	case churp.MethodCreate:
		var cfg churp.CreateRequest
		if err := cbor.Unmarshal(tx.Body, &cfg); err != nil {
			return api.ErrInvalidArgument
		}
		return ext.create(ctx, &cfg)
	case churp.MethodUpdate:
		var cfg churp.UpdateRequest
		if err := cbor.Unmarshal(tx.Body, &cfg); err != nil {
			return api.ErrInvalidArgument
		}
		return ext.update(ctx, &cfg)
	case churp.MethodApply:
		var reg churp.SignedApplicationRequest
		if err := cbor.Unmarshal(tx.Body, &reg); err != nil {
			return api.ErrInvalidArgument
		}
		return ext.apply(ctx, &reg)
	case churp.MethodConfirm:
		var reg churp.SignedConfirmationRequest
		if err := cbor.Unmarshal(tx.Body, &reg); err != nil {
			return api.ErrInvalidArgument
		}
		return ext.confirm(ctx, &reg)

	default:
		panic(fmt.Sprintf("keymanager: churp: invalid method: %s", tx.Method))
	}
}

// BeginBlock implements api.Extension.
func (ext *churpExt) BeginBlock(ctx *tmapi.Context) error {
	changed, epoch := ext.state.EpochChanged(ctx)
	if !changed {
		return nil
	}

	return ext.onEpochChange(ctx, epoch)
}

// EndBlock implements api.Extension.
func (*churpExt) EndBlock(*tmapi.Context) error {
	return nil
}

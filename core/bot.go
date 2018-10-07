package core

import (
	"vngo/core/protocol"
	"vngo/pkg/gateway/mockgateway"
	"vngo/pkg/module/mockmodule"

	"go.uber.org/zap"
)

func newCoordinators(ctx *protocol.ApplicationContext) []protocol.VtModule {
	// This order is important - it makes sure that the things taking requests start up before things sending requests
	return []protocol.VtModule{
		&mockmodule.MockModule{
			Ctx: ctx,
			Log: ctx.Logger.With(
				zap.String("type", "module"),
				zap.String("name", "mock"),
			),
		},
	}
}

func newGateways(ctx *protocol.ApplicationContext) []protocol.VtGateway {
	return []protocol.VtGateway{
		&mockgateway.MockGateway{
			Ctx: ctx,
			Log: ctx.Logger.With(
				zap.String("type", "gateway"),
				zap.String("name", "mock"),
			),
		},
	}
}

func SetupBot(app *protocol.ApplicationContext) error {

	// Wait until we're told to exit

	// Stop the coordinators in the reverse order. This assures that request senders are stopped before request servers

	// Exit cleanly
	return nil
}

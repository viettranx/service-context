package core

import sctx "github.com/viettranx/service-context"

func Recover() {
	if r := recover(); r != nil {
		sctx.GlobalLogger().GetLogger("recovered").Errorln(r)
	}
}

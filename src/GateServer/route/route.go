package route

import (
	"GateServer/pack"
	"gameLog"
	"fmt"
)

// 必须保证是线程安全的
func RouteIn(p *pack.Pack) {
	//gameLog.Debug("routing: " + fmt.Sprintf("sid: %v data: %v", p.Sid, p.Data))
	RouteOut(p)
}

// 必须保证是线程安全的
func RouteOut(p *pack.Pack) {
	if lk, ok := Instance.GateServer.GetLink(p.Sid); ok {
		lk.WtSyncChan <- p.Data
	} else {
		gameLog.Warn("a pack has benn drop " + fmt.Sprintf("sid: %v data: %v", p.Sid, p.Data))
	}
}
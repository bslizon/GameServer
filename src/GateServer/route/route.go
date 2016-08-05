package route

import (
	"GateServer/pack"
	"gameLog"
)

// 必须保证是线程安全的
// 必须保证无阻塞
func RouteIn(p *pack.Pack) {
	gameLog.Debug(p.Data)
}

// 必须保证是线程安全的
// 必须保证无阻塞
func RouteOut(p *pack.Pack) {

}
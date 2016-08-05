package route

import (
	"GateServer/pack"
	"gameLog"
)

// 必须保证是线程安全的
// 必须保证无阻塞
// 所有的分包goroutine都会调用这个方法分发自己的包
func Route(p *pack.Pack) {
	gameLog.Debug(p.Data)
}
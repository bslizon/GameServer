package Server

import (
	"GateServer/pack"
	"gameLog"
	"fmt"
	"errors"
)

// 必须保证是线程安全的
func RouteIn(p *pack.Pack) (err error) {
	// panic转error
	defer func() {
		if x := recover(); x != nil {
			switch value := x.(type) {
			case error:
				err = value
			case string:
				err = errors.New(value)
			default:
				err = errors.New(fmt.Sprintf("unknown panic: %#v. ", value))
			}
		}
	}()
	////////////////////////////////////////////////////////////////////

	//gameLog.Debug(fmt.Sprintf("routing. sid: %d data: %v", p.Sid, p.Data))
	RouteOut(p)

	////////////////////////////////////////////////////////////////////
	err = nil
	return
}

// 必须保证是线程安全的
func RouteOut(p *pack.Pack) (err error) {
	// panic转error
	defer func() {
		if x := recover(); x != nil {
			switch value := x.(type) {
			case error:
				err = value
			case string:
				err = errors.New(value)
			default:
				err = errors.New(fmt.Sprintf("unknown panic: %#v ", value))
			}
		}
	}()
	////////////////////////////////////////////////////////////////////

	if lk, ok := GateServer.GetLink(p.Sid); ok {
		lk.wtSyncChan <- p.Data
	} else {
		gameLog.Warn(fmt.Sprintf("link missing, a pack has benn drop. sid: %d data: %v. ", p.Sid, p.Data))
	}

	////////////////////////////////////////////////////////////////////
	err = nil
	return
}
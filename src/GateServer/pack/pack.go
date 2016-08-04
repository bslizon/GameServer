package pack

//数据包结构
type Pack struct {
	Data []byte
}

func New(b []byte) *Pack {
	p := new(Pack)
	p.Data = b
	return p
}
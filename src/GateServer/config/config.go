package config

import "math"

const(
	// 监听端口
	EXTERNAL_LISTEN_PORT = "8080"
	INTERNAL_LISTEN_PORT = "10000"

	//允许的最大同时TCP连接数
	MAX_TCP_CONN = 100000

	//外部TCP读写超时
	TCP_READ_TIMEOUT = 600	// sec
	TCP_WRITE_TIMEOUT = 300	// sec

	//进出协议的规格，尺寸以及解析相关
	PACK_DATA_SIZE_TYPE_LEN = 4	// sizeof(PackDataSizeType)
	MAX_INBOUND_PACK_DATA_SIZE = 1 << 14	// 16KB PackDataSizeType表示
	MAX_OUTBOUND_PACK_DATA_SIZE = 1 << 20	// 1MB

	//TcpLink接收chan相关参数
	WRITE_PACK_SYNC_CHAN_SIZE = 10
	WRITE_PACK_SYNC_CHAN_TIMEOUT = 20 // sec

	//sid相关
	BROCASTING_SID = math.MaxUint64 // 这个和 SocketIdType 对应
	DROP_SID = 0

	//性能分析相关
	PROFILE_FILE = "E:/GateServerProfile.prof"
)

type SocketIdType uint64 // 这个和 BROCASTING_SID 对应
type PackDataSizeType int32 // 关系到 PACK_DATA_SIZE_TYPE_LEN
package common

type MemPool struct {
	BlockSize int
	PoolSize  int
	memCh     chan []byte
}

func NewMemPool(blockSize, poolSize int) *MemPool {
	return &MemPool{
		BlockSize: blockSize,
		PoolSize:  poolSize,
		memCh:     make(chan []byte, poolSize),
	}
}

func (pool *MemPool) Get() []byte {
	var block []byte
	select {
	case block = <-pool.memCh:
	default:
		block = make([]byte, pool.BlockSize)
	}
	return block
}

func (pool *MemPool) Put(block []byte) {
	select {
	case pool.memCh <- block:
	default:
	}
}

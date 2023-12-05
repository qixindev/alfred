package cache

import (
	"github.com/allegro/bigcache"
	"time"
)

func NewBigCache(exp time.Duration) (*bigcache.BigCache, error) {
	config := bigcache.Config{
		Shards:             1024,            // 设置分区的数量
		LifeWindow:         exp,             // LifeWindow后,缓存对象被认为不活跃,但并不会删除对象
		CleanWindow:        5 * time.Second, // CleanWindow后，会删除被认为不活跃的对象，0代表不操作；
		MaxEntriesInWindow: 1024 << 2,       // 设置最大存储对象数量，仅在初始化时可以设置
		MaxEntrySize:       1024,            // 缓存对象的最大字节数，仅在初始化时可以设置
		Verbose:            true,            // 是否打印内存分配信息
		HardMaxCacheSize:   0,               // 设置缓存最大值(单位为MB),0表示无限制
		OnRemove:           nil,             // 在缓存过期或者被删除时,可设置回调函数，参数是(key、val)，默认是nil不设置
		OnRemoveWithReason: nil,             // 在缓存过期或者被删除时,可设置回调函数，参数是(key、val,reason)默认是nil不设置
	}
	bigCache, err := bigcache.NewBigCache(config)
	if err != nil {
		return nil, err
	}

	return bigCache, nil
}

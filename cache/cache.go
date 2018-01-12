package cache

import (
	"time"

	cache "github.com/patrickmn/go-cache"
)

func init() {
	//初始化cache
	GetCache()
}

var c *cache.Cache

//GetCache 获取cache对象
func GetCache() *cache.Cache {
	if c != nil {
		return c
	}
	c = cache.New(5*time.Minute, 10*time.Minute)
	return c
}

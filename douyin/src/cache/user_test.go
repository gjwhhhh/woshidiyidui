package cache

import (
	"douyin/src/pojo/entity"
	"fmt"
	"testing"
)

func TestUserLruCache1(t *testing.T) {
	userLruCache := UserCacheConstructor(50)
	fmt.Println(userLruCache.size)
	userLruCache.Put(1, &entity.DyUser{Username: "张三"})
	fmt.Println(userLruCache.size)
	fmt.Println(userLruCache.Head().Username)
	userLruCache.Put(2, &entity.DyUser{Username: "李四"})
	fmt.Println(userLruCache.size)
	fmt.Println(userLruCache.Head().Username)
	userLruCache.Get(1)
	fmt.Println(userLruCache.size)
	fmt.Println(userLruCache.Head().Username)
}

func TestUserLruCache2(t *testing.T) {
	userLruCache := UserCacheConstructor(10)
	for i := 0; i < 20; i++ {
		fmt.Println(userLruCache.size)
		userLruCache.Put(int64(i), &entity.DyUser{Username: fmt.Sprintf("用户%d", i)})
	}
	fmt.Println(userLruCache.Head().Username)
	userLruCache.Get(15)
	fmt.Println(userLruCache.Head().Username)
}

package user_uname_pwd

import (
	"douyin/src/pojo/entity"
	"sync"
)

var lock sync.Mutex

// UserLruCache 用户缓存
type UserLruCache struct {
	size       int
	capacity   int
	cache      map[string]*node
	head, tail *node
}

// node 内部节点
type node struct {
	k         string
	v         *entity.DyUser
	pre, next *node
}

func initNode(key string, value *entity.DyUser) *node {
	return &node{
		k: key,
		v: value,
	}
}

func UserCacheConstructor(capacity int) *UserLruCache {
	lruCache := &UserLruCache{
		size:     0,
		capacity: capacity,
		cache:    map[string]*node{},
		head:     initNode("", nil),
		tail:     initNode("", nil),
	}
	lruCache.head.next = lruCache.tail
	lruCache.tail.pre = lruCache.head
	return lruCache
}

func (this *UserLruCache) Head() *entity.DyUser {
	return this.head.next.v
}

func (this *UserLruCache) Get(key string) *entity.DyUser {
	if _, ok := this.cache[key]; !ok {
		return nil
	}
	node := this.cache[key]
	this.moveToHead(node)
	return node.v
}

func (this *UserLruCache) Put(key string, value *entity.DyUser) {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := this.cache[key]; !ok {
		node := initNode(key, value)
		this.cache[key] = node
		this.addToHead(node)
		this.size++
		if this.size > this.capacity {
			tail := this.removeTail()
			delete(this.cache, tail.k)
			this.size--
		}
	} else {
		node := this.cache[key]
		node.v = value
		this.moveToHead(node)
	}
}

func (this *UserLruCache) addToHead(node *node) {
	node.pre = this.head
	node.next = this.head.next
	this.head.next.pre = node
	this.head.next = node
}

func (this *UserLruCache) removeNode(node *node) {
	node.pre.next = node.next
	node.next.pre = node.pre
}

func (this *UserLruCache) moveToHead(node *node) {
	this.removeNode(node)
	this.addToHead(node)
}

func (this *UserLruCache) removeTail() *node {
	node := this.tail.pre
	this.removeNode(node)
	return node
}

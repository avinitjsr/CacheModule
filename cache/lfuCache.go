package cache

import "container/heap"

type LFUCache struct {
	capacity int
	cache    map[int]*lfuItem
	heap     *lfuHeap
}
type lfuItem struct {
	key       int
	value     int
	frequency int
	index     int
}

type lfuHeap []*lfuItem

func (h lfuHeap) Len() int           { return len(h) }
func (h lfuHeap) Less(i, j int) bool { return h[i].frequency < h[j].frequency }
func (h lfuHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *lfuHeap) Push(x interface{}) {
	n := len(*h)
	item := x.(*lfuItem)
	item.index = n
	*h = append(*h, item)
}

func (h *lfuHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*h = old[0 : n-1]
	return item
}

func (h *lfuHeap) update(item *lfuItem, value int, frequency int) {
	item.value = value
	item.frequency = frequency
	heap.Fix(h, item.index)
}

// NewLFUCache creates a new LFUCache with the given capacity
func NewLFUCache(capacity int) *LFUCache {
	h := &lfuHeap{}
	heap.Init(h)
	return &LFUCache{
		capacity: capacity,
		cache:    make(map[int]*lfuItem),
		heap:     h,
	}
}

// Get retrieves a value from the cache
func (c *LFUCache) Get(key int) (int, bool) {
	if item, found := c.cache[key]; found {
		item.frequency++
		heap.Fix(c.heap, item.index)
		return item.value, true
	}
	return 0, false
}

// Put inserts a value into the cache
func (c *LFUCache) Put(key, value int) {
	if c.capacity == 0 {
		return
	}

	if item, found := c.cache[key]; found {
		item.frequency++
		c.heap.update(item, value, item.frequency)
		return
	}

	if len(c.cache) == c.capacity {
		item := heap.Pop(c.heap).(*lfuItem)
		delete(c.cache, item.key)
	}

	item := &lfuItem{
		key:       key,
		value:     value,
		frequency: 1,
	}
	c.cache[key] = item
	heap.Push(c.heap, item)
}

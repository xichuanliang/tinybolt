package page

import "fmt"

// 在文件中使用bucket表示
type bucket struct {
	root     pgid
	sequence uint64
}

type Bucket struct {
	*bucket
	page     *page              // page优化子Bucket
	buckets  map[string]*Bucket // 子Bucket
	rootNode *node              // root node
	nodes    map[pgid]*node     // node的cache
}

// 根据pgid 获取 page / node
func (b *Bucket) pageNode(id pgid) (*page, *node) {
	// 从根节点中获取node
	if b.root == 0 {
		if id != 0 {
			panic(fmt.Sprintf("inline bucket non-zero page access(2): %d != 0", id))
		}
		if b.rootNode != nil {
			return nil, b.rootNode
		}
		return b.page, nil
	}
	// 从nodes中获取node
	if b.nodes != nil {
		if n := b.nodes[id]; n != nil {
			return nil, n
		}
	}
	return nil, nil
}

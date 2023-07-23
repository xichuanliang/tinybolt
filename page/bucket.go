package page

// 在文件中使用bucket表示
type bucket struct {
	root     pgid
	sequence uint64
}

type Bucket struct {
	*bucket
	page     *page              // 暂时不用inline page优化子Bucket
	buckets  map[string]*Bucket // 子Bucket
	rootNode *node              // root node
	nodes    map[pgid]*node     // node的cache
}


// 获取node
func (b *Bucket) pageNode(id pgid) (*page, *node) {
	// 从根节点中获取node
	if b.rootNode != nil {
		return nil, b.rootNode
	}
	// 从nodes中获取node
	if b.nodes != nil {
		if n := b.nodes[id]; n != nil {
			return nil, n
		}
	}
	return nil, nil
}

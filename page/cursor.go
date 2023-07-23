package page

import (
	"bytes"
	"sort"
)

type Cursor struct {
	bucket *Bucket   // 对应的bucket
	stack  []elemRef // 将递归遍历到的值存储到栈中
}

// 移动cursor找到第一个在Bucket中的元素，并返回它的key-value
func (c *Cursor) First(key []byte, value []byte) {
	// 将c.stack清零，但是底层的数组并没有释放
	c.stack = c.stack[:0]
	// 根据根root找到page或者node
	p, n := c.bucket.pageNode(c.bucket.root)
	// 将root节点所在的page、node信息放在栈顶，index为0，表示第一个从子节点开始遍历
	c.stack = append(c.stack, elemRef{
		page:  p,
		node:  n,
		index: 0,
	})

}

// curse移动到栈中最后一页下的第一个叶子元素,
func (c *Cursor) first() {
	for { // 找到最左边第一个叶子节点
		// 每次循环取出最后一个元素
		ref := &c.stack[len(c.stack)-1]
		if ref.isLeaf() {
			break
		}
		var pgid pgid
		if ref.node != nil {
			pgid = ref.node.innodes[ref.index].pgid
		} else {
			pgid = ref.page.branchPageElement(uint16(ref.index)).pgid
		}
		// 根据pgid获得page，或者node
		p, n := c.bucket.pageNode(pgid)
		c.stack = append(c.stack, elemRef{
			page:  p,
			node:  n,
			index: 0,
		})
	}
}

// 在stack顶的leaf node搜索一个Key
func (c *Cursor) nserach(key []byte) {
	e := c.stack[len(c.stack)-1]
	_, n := *e.page, *e.node
	// func(i int)为false时，二分搜索向右搜索
	index := sort.Search(len(n.innodes), func(i int) bool {
		return bytes.Compare(n.innodes[i].key, key) != -1
	})
	e.index = index
}

func (c *Cursor) searchNode(key []byte, n *node) {
	var exact bool
	index := sort.Search(len(n.innodes), func(i int) bool {
		com := bytes.Compare(n.innodes[i].key, key)
		if com == 0 {
			exact = true
		}
		return com != -1
	})
	if exact && index > 0 {
		index--
	}
	c.stack[len(c.stack)-1].index = index

}

// elemRef 表示在一个给定的page/node上针对某个元素的引用
type elemRef struct {
	page  *page // 页面
	node  *node // 内存中页面的信息
	index int   // 保存在当前page、node遍历到那个节点
}

// 判断 ref是否指向一个 leaf node/ page
func (e *elemRef) isLeaf() bool {
	if e.page != nil {
		return (e.page.flags & leafPageFlag) != 0
	}
	return e.node.isLeaf
}

// 返回node/page的elem数
func (e *elemRef) count() int {
	if e.node != nil {
		return len(e.node.innodes)
	}
	return int(e.page.count)
}

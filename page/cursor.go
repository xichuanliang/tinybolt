package page

import (
	"bytes"
	"fmt"
	"sort"
)

type Cursor struct {
	bucket *Bucket   // 对应的bucket
	stack  []elemRef // 将递归遍历到的值存储到栈中
}

// 找到指定的Key-value
func (c *Cursor) Seek(seek []byte) (key []byte, value []byte) {
	k, v, flags := c.seek(seek)
	if ref := &c.stack[len(c.stack)-1]; ref.index >= ref.count() {
		k, v, flags = c.next()
	}
	if k == nil {
		return nil, nil
		// 是子桶
	} else if (flags & uint32(bucketLeafFlag)) != 0 {
		return k, nil
	}
	return k, v
}

// 找到指定的key，并返回key value
func (c *Cursor) seek(seek []byte) (key []byte, value []byte, flags uint32) {
	// 清空stack
	c.stack = c.stack[:0]
	// 根据seek的key值搜索root
	c.search(seek, c.bucket.root)
	// c.stack保存了所有遍历过的节点
	ref := &c.stack[len(c.stack)-1]
	// cursor指向page/node的末尾，返回nil，说明没找到
	if ref.index >= ref.count() {
		return nil, nil, 0
	}
	return c.keyValue()
}

// 从leaf elemRef中获取 key / value
func (c *Cursor) keyValue() ([]byte, []byte, uint32) {
	// 去最后一个elemRef
	ref := &c.stack[len(c.stack)-1]
	if ref.count() == 0 || ref.index >= ref.count() {
		return nil, nil, 0
	}
	// 存在leaf node，返回节点对应index的key value
	if ref.node != nil {
		innode := ref.node.innodes[ref.index]
		return innode.key, innode.value, innode.flags
	}
	// 从page中 返回对应节点的Key value
	elem := ref.page.leafPageElement(uint16(ref.index))
	return elem.key(), elem.value(), elem.flags
}

// 递归搜索，对指定的page/node执行二分搜索，直到找到给定的Key
func (c *Cursor) search(key []byte, pgid pgid) {
	// 根据pgid获得 page / node
	p, n := c.bucket.pageNode(pgid)
	// 如果返回page，并且该page既不是leaf page ，也不是 brach page，则panic
	if p != nil && (p.flags&(branchPageFlag|leafPageFlag)) == 0 {
		panic(fmt.Sprintf("invalid page type: &d: %x", p.id, p.flags))
	}
	e := elemRef{
		page: p,
		node: n,
	}
	// 记录遍历过的路径
	c.stack = append(c.stack, e)
	if e.isLeaf() {
		c.nserach(key)
		return
	}
}

// 移动到下一个叶子节点，返回key value，中序遍历
func (c *Cursor) next() (key []byte, value []byte, flags uint32) {
	for {
		var i int
		// 遍历stack中的每个索引，如果索引所在的page/node不是最后一个元素，就往后移动一个位置
		for i = len(c.stack) - 1; i >= 0; i-- {
			elem := &c.stack[i]
			if elem.index < elem.count() {
				elem.index++
				break
			}
			// 遍历完了所有页面
			if i == -1 {
				return nil, nil, 0
			}
			// 剩余的节点里面找，跳过原先遍历过的节点
			c.stack = c.stack[:i+1]
			// 如果是叶子节点，first()啥都不做，直接退出。返回elem.index+1的数据
			// 非叶子节点的话，需要移动到stack中最后一个路径的第一个元素
			c.first()
			if c.stack[len(c.stack)-1].count() == 0 {
				continue
			}
			return c.keyValue()
		}
	}
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

// 找到最后一个非叶子节点的第一个叶子节点。index=0的节点
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

// 在stack顶部的 leaf node/ page中 搜索Key
func (c *Cursor) nserach(key []byte) {
	e := &c.stack[len(c.stack)-1]
	p, n := e.page, e.node
	// 搜索node
	if n != nil {
		index := sort.Search(len(n.innodes), func(i int) bool {
			return bytes.Compare(n.innodes[i].key, key) != -1
		})
		e.index = index
		return
	}
	// 获取leaf page中的所有elem
	inodes := p.leafPageElements()
	// 二分搜索page中的elem
	index := sort.Search(len(inodes), func(i int) bool {
		return bytes.Compare(inodes[i].key(), key) != -1
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

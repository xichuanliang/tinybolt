package page

import (
	"unsafe"

	"github.com/xichuan/tinybolt"
)

// branch Page和leaf Page被抽象为一个node表示，用isLeaf区分

type nodes []node

// leaf page 和 branch page 在内存中都是使用node表示的
type node struct {
	isLeaf     bool
	unbalanced bool
	spilled    bool   // 是否内存溢出
	key        []byte // 保存第一个key
	pgid       pgid
	parent     *node
	children   nodes
	innodes    innodes
}

// 通过page初始化node
func (n *node) read(p *page) {
	n.pgid = p.id
	n.isLeaf = ((p.flags & leafPageFlag) != 0)
	n.innodes = make(innodes, int(p.count))
	for i := 0; i < int(p.count); i++ {
		inode := n.innodes[i]
		if n.isLeaf {
			// leaf node
			elem := p.leafPageElement(uint16(i))
			inode.flags = elem.flags
			inode.key = elem.key()
			inode.value = elem.value()
		} else {
			// branch node
			elem := p.branchPageElement((uint16(i)))
			inode.key = elem.key()
			inode.pgid = elem.pgid
		}
		tinybolt.Assert(len(inode.key) > 0, "read: zero-length inode key")
	}
	// 保存第一个Key
	if len(n.innodes) > 0 {
		n.key = n.innodes[0].key
		tinybolt.Assert(len(n.key) > 0, "read: zero-length node key")
	} else {
		n.key = nil
	}
}

// 将node写入page
func (n *node) write(p *page) {
	if n.isLeaf {
		p.flags |= leafPageFlag
	} else {
		p.flags |= branchPageFlag
	}
	p.id = n.pgid
	p.count = uint16(len(n.innodes))
	// buf为存储的key-value的切片
	buf := (*[tinybolt.MaxAllocSize]byte)(unsafe.Pointer(&p.ptr))[n.pageElementSize()*len(n.innodes):]
	// 逐个将node中的ionode封装为对应的element,
	for idx, innode := range n.innodes {
		if n.isLeaf {
			elem := p.leafPageElement(uint16(idx))
			elem.pos = uint32(uintptr(unsafe.Pointer(&buf[0])) - uintptr(unsafe.Pointer(elem)))
			elem.flags = innode.flags
			elem.ksize = uint32(len(innode.key))
			elem.vsize = uint32(len(innode.value))
		} else {
			elem := p.branchPageElement(uint16(idx))
			elem.pos = uint32(uintptr(unsafe.Pointer(&buf[0])) - uintptr(unsafe.Pointer(elem)))
			elem.ksize = uint32(len(innode.key))
			elem.pgid = innode.pgid
		}
		// 将element后面加入Key-val
		klen, vlen := len(innode.key), len(innode.value)
		copy(buf[0:], innode.key)
		buf = buf[klen:]
		copy(buf[0:], innode.value)
		buf = buf[vlen:]
	}
}

// 返回page中leaf元素/branch元素所占空间的大小
func (n *node) pageElementSize() int {
	if n.isLeaf {
		return leafPageElementSize
	}
	return branchPageElementSize
}

type innodes []inode

// 节点内部的一个内部节点，指向page中元素，或者指向尚未添加到page的元素
type inode struct {
	flags uint32 // 是 leaf节点/brach节点
	pgid  pgid
	key   []byte
	value []byte // inode为branch节点时，value为空 inode为leaf 节点时，value不为空
}

type branchPageElement struct {
	pos   uint32 // key相对于当前page数据部分的偏移量
	ksize uint32 // key的大小
	pgid  pgid   // page的id
}

// 返回b的ksize，将其封装为[]byte
func (b *branchPageElement) key() []byte {
	// buf 为一个指向*[utils.MaxAllocSize]byte类型的指针
	buf := (*[tinybolt.MaxAllocSize]byte)(unsafe.Pointer(b))
	return (*[tinybolt.MaxAllocSize]byte)(unsafe.Pointer(&buf[b.pos]))[:b.ksize:b.ksize]
}

// 保存key-value的元素信息
type leafPageElement struct {
	flags uint32 // 0：普通的叶子节点元素  1：子bucket
	pos   uint32 //key距离leafPageElement的偏移位置
	ksize uint32 // key的尺寸
	vsize uint32 // value的尺寸
}

// key的字节切片
func (l *leafPageElement) key() []byte {
	buf := (*[tinybolt.MaxAllocSize]byte)(unsafe.Pointer(l))
	return (*[tinybolt.MaxAllocSize]byte)(unsafe.Pointer(&buf[l.pos]))[:l.ksize:l.ksize]
}

// value的字节切片
func (l *leafPageElement) value() []byte {
	buf := (*[tinybolt.MaxAllocSize]byte)(unsafe.Pointer(l))
	return (*[tinybolt.MaxAllocSize]byte)(unsafe.Pointer(&buf[l.pos+l.ksize]))[:l.vsize:l.vsize]
}

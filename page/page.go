package page

// page是物理页和内存page之间的映射
import (
	"fmt"
	"unsafe"

	"github.com/xichuan/tinybolt"
)

// 4种page
const (
	branchPageFlag   = 0x01 //0000 0001
	leafPageFlag     = 0x02 //0000 0010
	metaPageFlag     = 0x04 //0000 0100
	freelistPageFlag = 0x10 //0000 1010
)

const bucketLeafFlag = 0x01

// branch元素所占大小
const branchPageElementSize = int(unsafe.Sizeof(branchPageElement{}))

// leaf 元素所占大小
const leafPageElementSize = int(unsafe.Sizeof(leafPageElement{}))

type pgid uint64

type page struct {
	id       pgid    // 页id
	flags    uint16  // 4种page类型
	count    uint16  // 节点数目 统计叶子节点、非叶子节点的个数
	overflow uint32  // 数据是否溢出，当一个page存不下时就会溢出
	ptr      uintptr // 保存数据的首指针
}

func (p *page) typ() string {
	if (p.flags & branchPageFlag) != 0 {
		return "branch"
	} else if (p.flags & leafPageFlag) != 0 {
		return "leaf"
	} else if (p.flags & metaPageFlag) != 0 {
		return "meta"
	} else if (p.flags & freelistPageFlag) != 0 {
		return "freelist"
	}
	return fmt.Sprintf("unknown<%02x>", p.flags)
}

// 返回page中meta的首地址
func (p *page) meta() *meta {
	return (*meta)(unsafe.Pointer(&p.ptr))
}

// 获得freelist page 中的下标为index的branchPageElement
func (p *page) branchPageElement(index uint16) *branchPageElement {
	return &(*[tinybolt.MaxElementSize]branchPageElement)(unsafe.Pointer(&p.ptr))[index]
}

// 获取page的整个branchPageElements
func (p *page) branchPageElements() []branchPageElement {
	if p.count == 0 {
		return nil
	}
	return (*[tinybolt.MaxElementSize]branchPageElement)(unsafe.Pointer(&p.ptr))[:]
}

// 获取第index个leafpageElement
func (p *page) leafPageElement(index uint16) *leafPageElement {
	return &(*[tinybolt.MaxElementSize]leafPageElement)(unsafe.Pointer(&p.ptr))[index]
}

// 获取整个leafPageElements
func (p *page) leafPageElements() []leafPageElement {
	if p.count == 0 {
		return nil
	}
	return (*[tinybolt.MaxElementSize]leafPageElement)(unsafe.Pointer(&p.ptr))[:]
}

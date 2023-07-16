package page

import (
	"fmt"
	"hash/fnv"
	"unsafe"

	"github.com/xichuan/tinybolt"
)

// 保存page的元数据信息
type meta struct {
	magic    uint32 // magic number，值为“0xED0CDAED”,验证数据库文件的有效性和正确性
	version  uint32 // 版本号
	pagesize uint32 // page的页面大小
	freelist pgid   // 保存freelist页面的id
	pgid     pgid   // 保存当前总的页面数量，即最大页面号+1
	checksum uint64 // 校验码，用于校验meta页面是否出错
}

// 检查数据库的完整性和一致性
func (m *meta) validate() error {
	if m.magic != tinybolt.Magic {
		return tinybolt.ErrInvalid
	} else if m.version != tinybolt.Version {
		return tinybolt.ErrVersionMismatch
	} else if m.checksum != 0 && m.checksum != m.sum64() {
		return tinybolt.ErrChecksum
	}
	return nil
}

// 从磁盘上读取meta，并封装到meta中
func (m *meta) write(p *page) {
	if m.freelist > m.pgid {
		panic(fmt.Sprintf("freelist pgid (%d) above high water mark (%d)", m.freelist, m.pgid))
	}
	// 设置page的标识为meta
	p.flags |= metaPageFlag
	// 将page中的meta页置换到内存中的meta
	m.copy(p.meta())
}

// 将page.ptr指针转换为meta类型返回
func (m *meta) copy(dest *meta) {
	*m = *dest
}

func (m *meta) sum64() uint64 {
	// 获取新的 FNV-1a 64 位哈希对象
	var h = fnv.New64a()
	// unsafe.Offsetof() 结构体字段的偏移量
	_, _ = h.Write((*[unsafe.Offsetof(meta{}.checksum)]byte)(unsafe.Pointer(m))[:])
	return h.Sum64()
}

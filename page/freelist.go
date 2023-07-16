package page

// freelist page中的数据部分仅存储可用的page ids
import (
	"fmt"
	"unsafe"

	"github.com/xichuan/tinybolt"
)

type freelist struct {
	ids   []pgid        // 已经可以被分配的空闲页
	cache map[pgid]bool // 快速查找出空闲页
}

// 创建一个freelist
func newFreelist() *freelist {
	return &freelist{
		ids:   make([]pgid, 0),
		cache: make(map[pgid]bool),
	}
}

// 将freelist写入到page中
func (f *freelist) write(p *page) error {
	// 设置page的表示为freelist
	p.flags |= freelistPageFlag
	lenids := f.count()
	p.count = uint16(lenids)
	// 将page.ptr的指针指向指针数组，并将freelist的ids复制到该数组中
	f.copyall(((*[tinybolt.MaxAllocSize]pgid)(unsafe.Pointer(&p.ptr)))[:])
	return nil
}

// free pages的个数
func (f *freelist) count() int {
	return len(f.ids)
}

func (f *freelist) copyall(dst []pgid) {
	copy(dst, f.ids)
}

// 从page中读取freelist，并转换为freelist
func (f *freelist) read(p *page) error {
	idx, count := 0, int(p.count)
	// 如果freelist的page中count==0，说明当前page中无空闲页
	if count == 0 {
		f.ids = nil
	} else {
		// 将freslist的page中的pgidx存储到 freelist中
		ids := ((*[tinybolt.MaxAllocSize]pgid)(unsafe.Pointer(&p.ptr)))[idx:count]
		f.ids = make([]pgid, len(ids))
		copy(f.ids, ids)
	}
	// 重建freelist 的 cache
	f.reindex()
	return nil
}

// 基于free page 重新创建freelist的cache
func (f *freelist) reindex() {
	f.cache = make(map[pgid]bool, len(f.ids))
	for _, id := range f.ids {
		f.cache[id] = true
	}
}

// 根据给定的尺寸分配连续的page，如果没有连续的page，返回0
func (f *freelist) allocate(n int) pgid {
	if len(f.ids) == 0 {
		return 0
	}
	var initial, previd pgid
	for idx, id := range f.ids {
		if id <= 1 {
			panic(fmt.Sprintf("invalid page allocation: %d", id))
		}
		// id-previd==1判断是否连续
		if previd == 0 || id-previd != 1 {
			// 记录每一次不连续的位置
			initial = id
		}
		if (id-initial)+1 == pgid(n) {
			// 找到前n个空闲位置，freelist中的ids记录空闲的pid，就直接后移n个
			if (idx + 1) == n {
				f.ids = f.ids[idx+1:]
			} else {
				copy(f.ids[idx-n+1:], f.ids[idx+1:])
				f.ids = f.ids[:len(f.ids)-n]
			}

			// 更新cache
			for i := pgid(0); i < pgid(n); i++ {
				delete(f.cache, initial+i)
			}
			return initial
		}
		previd = id
	}
	return 0
}

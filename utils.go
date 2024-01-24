package tinybolt

import "fmt"

// 创建数组指针时使用的大小
const (
	MaxAllocSize   = 0x7FFFFFFF // int32 的最大值
	MaxElementSize = 0x7FFFFFF  // int16的最大值
	Magic          = 0xED0CDAED
	Version        = 2
)

func Assert(conditon bool, msg string, v ...interface{}) {
	if !conditon {
		panic(fmt.Sprintf("assertion failed: "+msg, v...))
	}
}

func a() {

}

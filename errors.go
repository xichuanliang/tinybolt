package tinybolt

import "errors"

var (
	// 当数据库上的两个元数据页都无效时，会返回 ErrInvalid
	// 这通常发生在文件不是 BoltDB 数据库时
	ErrInvalid = errors.New("invalid database")
	// 当数据文件是使用不同版本的 Bolt 创建时，将返回 ErrVersionMismatch。
	ErrVersionMismatch = errors.New("version mismatch")
	// 当元数据页的校验和不匹配
	ErrChecksum = errors.New("checksum error")
)

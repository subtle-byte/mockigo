package mockery

import "unsafe"

type UnsafeInterface interface {
	Do(ptr *unsafe.Pointer)
}

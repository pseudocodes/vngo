package core

import (
	"syscall"
)

func internalDup2(oldfd uintptr, newfd uintptr) error {
	return syscall.Dup2(int(oldfd), int(newfd))
}

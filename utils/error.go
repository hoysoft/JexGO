package utils

import (
	"fmt"
)

func CatchPanic() error{
	if err := recover(); err != nil {
		return err.(error)
	}
}

package utils


func CatchPanic() error{
	if err := recover(); err != nil {
		return err.(error)
	}
	return nil
}

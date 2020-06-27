package pidfile

import (
	"errors"
	"os"
	"strconv"
	"syscall"
)

var emptyPidFile = errors.New("empty pidfile path")

func Write(fileName string) (err error) {
	var f *os.File
	if fileName == "" {
		return emptyPidFile
	}
	if f, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0640); err != nil {
		return
	}
	if err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		return
	}
	if _, err = f.Write([]byte(strconv.Itoa(os.Getpid()))); err != nil {
		return
	}
	err = f.Sync()
	return
}

func Unlink(fileName string) {
	_ = os.Remove(fileName)
}

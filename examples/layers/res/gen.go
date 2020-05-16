package res

import (
	"github.com/gabstv/tau/io"
	"github.com/gabstv/tau/io/broccolifs"
)

//go:generate broccoli -var=rx

func FS() io.Filesystem {
	return broccolifs.New(rx)
}

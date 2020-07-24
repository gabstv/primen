package res

import (
	"github.com/gabstv/primen/io"
	"github.com/gabstv/primen/io/broccolifs"
)

//go:generate broccoli -var=rx

func FS() io.Filesystem {
	return broccolifs.New(rx)
}

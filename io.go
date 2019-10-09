package glogrotation_cli

import (
	"io"
)

type ChanRW struct {
	C chan []byte
}

func NewChanRW() *ChanRW {
	return &ChanRW{make(chan []byte)}
}

func (this *ChanRW) Read(p []byte) (n int, err error) {
	d := <-this.C
	if d == nil || len(d) == 0 {
		err = io.EOF
		return
	}
	n = copy(p, d)
	return
}

func (this *ChanRW) Write(p []byte) (n int, err error) {
	this.C <- p
	return len(p), nil
}

func (this *ChanRW) Close() error {
	close(this.C)
	return nil
}

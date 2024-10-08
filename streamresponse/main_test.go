package main

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

// 使用管道传递数据
func TestPipe(t *testing.T) {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		for i := 0; i < 10; i++ {
			_, _ = pw.Write([]byte(fmt.Sprintf("hello world: %d\n", i)))
			time.Sleep(time.Second)
		}
	}()

	io.Copy(os.Stdout, pr)
}

type data struct {
	ch chan string
}

// 使用 channel 传递数据
func TestChannel(t *testing.T) {
	d := &data{
		ch: make(chan string),
	}
	go func() {
		for i := 0; i < 10; i++ {
			d.ch <- fmt.Sprintf("hello world: %d\n", i)
			time.Sleep(time.Second)
		}
		close(d.ch)
	}()

	io.Copy(os.Stdout, d)
}

func (d data) Read(p []byte) (n int, err error) {
	s, ok := <-d.ch
	if !ok {
		return 0, io.EOF
	}
	return copy(p, s), nil
}

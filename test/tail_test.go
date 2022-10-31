package test

import (
	"fmt"
	"github.com/hpcloud/tail"
	"os"
	"testing"
	"time"
)

func TestTail(t *testing.T) {

	testMap := make(map[int]int, 5)
	testMap[1] = 1 + 1
	testMap[1+1] = 2
	fmt.Println(testMap)
	os.Exit(1)

	filename := "./test.log"
	cfg := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	tails, err := tail.TailFile(filename, cfg)
	if err != nil {
		return
	}
	var (
		line *tail.Line
		ok   bool
	)
	for {
		line, ok = <-tails.Lines
		if !ok {
			fmt.Printf("tail file close reopen, filename:%s\n", tails.Filename)
			time.Sleep(time.Second)
			continue
		}
		fmt.Println("line:", line.Text)
	}
}

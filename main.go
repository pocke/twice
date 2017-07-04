package main

import (
	"errors"
	"io"
	"os"
	"strconv"
	"time"
)

var Usage = errors.New("Usage: twice [count]")

func main() {
	if err := Main(os.Stdin, os.Stdout, os.Args); err != nil {
		panic(err)
	}
}

func Main(in, out *os.File, args []string) error {
	if len(args) > 2 {
		return Usage
	}
	count := 2
	if len(args) == 2 {
		var err error
		count, err = strconv.Atoi(args[1])
		if err != nil {
			return err
		}
	}
	count--

	rec := NewRecorder()
	tee := io.MultiWriter(rec, out)
	io.Copy(tee, in)
	for i := 0; i < count || count < 0; i++ {
		err := rec.CopyWithWait(out)
		if err != nil {
			return err
		}
	}
	return nil
}

type Recorder struct {
	records    []Record
	beforeTime time.Time
}

func NewRecorder() *Recorder {
	return &Recorder{
		beforeTime: time.Now(),
	}
}

func (r *Recorder) Write(p []byte) (int, error) {
	now := time.Now()
	t := now.Sub(r.beforeTime)
	r.beforeTime = now
	bytes := make([]byte, len(p))
	copy(bytes, p)
	r.records = append(r.records, Record{bytes: bytes, time: t})
	return len(p), nil
}

func (r *Recorder) CopyWithWait(out io.Writer) error {
	for _, rec := range r.records {
		time.Sleep(rec.time)
		_, err := out.Write(rec.bytes)
		if err != nil {
			return err
		}
	}
	return nil
}

type Record struct {
	bytes []byte
	time  time.Duration
}

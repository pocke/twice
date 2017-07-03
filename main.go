package main

import (
	"io"
	"os"
	"time"
)

func main() {
	if err := Main(os.Stdin, os.Stdout); err != nil {
		panic(err)
	}
}

func Main(in, out *os.File) error {
	rec := NewRecorder()
	tee := io.MultiWriter(rec, out)
	io.Copy(tee, in)
	if err := rec.CopyWithWait(out); err != nil {
		return err
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

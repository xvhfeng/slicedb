package fileio

import (
	"../logger"

	"fmt"
	"io"
	"sync/atomic"
	"time"
)

type BufWriter struct {
	err      error
	locker   int32
	d        time.Duration
	buffSize int
	offset   int
	wr       io.Writer
	buf      []byte
	t        *time.Timer
	log      *logger.Log
}

const (
	defaultBufSize = 1000
)

func NewBufWriter(wr io.Writer,
	buffSize int, d time.Duration,
	log *logger.Log) (b *BufWriter) {
	if buffSize < defaultBufSize {
		buffSize = defaultBufSize
	}
	b = new(BufWriter)
	b.d = d
	b.buffSize = buffSize
	b.wr = wr
	b.log = log
	b.buf = make([]byte, buffSize)
	b.t = time.AfterFunc(d, func() {
		if !atomic.CompareAndSwapInt32(&b.locker, 0, 1) {
			return
		}
		b.t.Stop()
		for {
			err := b.flush()
			if nil == err && 0 == b.offset {
				break
			} else {
				if nil != log {
					log.Brush(logger.Error,
						"flush the buffer to disk is fail.err %v.",
						err)
				}
			}
		}
		b.t.Reset(b.d)
		atomic.CompareAndSwapInt32(&b.locker, 1, 0)
	})
	return
}

func (b *BufWriter) Close() {
	b.t.Stop()
	/* close(b.t.C) */
}

func (b *BufWriter) Reset(wr io.Writer) {
	b.err = nil
	b.offset = 0
	b.wr = wr
	b.t.Reset(b.d)
}

// Flush writes any buffered data to the underlying io.Writer.
func (b *BufWriter) Flush() error {
	err := b.flush()
	return err
}

func (b *BufWriter) flush() error {
	if b.err != nil {
		return b.err
	}
	if b.offset == 0 {
		return nil
	}
	n, err := b.wr.Write(b.buf[0:b.offset])
	if n < b.offset && err == nil {
		err = io.ErrShortWrite
	}
	if err != nil {
		if n > 0 && n < b.offset {
			copy(b.buf[0:b.offset-n], b.buf[n:b.offset])
		}
		b.offset -= n
		b.err = err
		return err
	}
	b.offset = 0
	return nil
}

func (b *BufWriter) Available() int { return len(b.buf) - b.offset }

func (b *BufWriter) Buffered() int { return b.offset }

// Write writes the contents of p into the buffer.
// It returns the number of bytes written.
// If nn < len(p), it also returns an error explaining
// why the write is short.
func (b *BufWriter) Write(p []byte) (nn int, err error) {
	times := 0
	for 3 > times {
		if atomic.CompareAndSwapInt32(&b.locker, 0, 1) {
			defer atomic.CompareAndSwapInt32(&b.locker, 1, 0)
		} else {
			times++
			if nil != b.log {
				b.log.Brush(logger.Warn, "get write locker is fail,sleep 1 ms and try again.")
			}
			time.Sleep(1 * time.Millisecond)
			continue
		}
		for len(p) > b.Available() && b.err == nil {
			var n int
			if b.Buffered() == 0 {
				// Large write, empty buffer.
				// Write directly from p to avoid copy.
				n, b.err = b.wr.Write(p)
			} else {
				n = copy(b.buf[b.offset:], p)
				b.offset += n
				b.flush()
			}
			nn += n
			p = p[n:]
		}
		if b.err != nil {
			return nn, b.err
		}
		n := copy(b.buf[b.offset:], p)
		b.offset += n
		nn += n
		return nn, nil
	}
	if nil != b.log {
		b.log.Brush(logger.Error, "write is fail,try 3 time later.")
	}
	return 0, fmt.Errorf("buffer writer is fail.")
}

// WriteString writes a string.
// It returns the number of bytes written.
// If the count is less than len(s), it also returns an error explaining
// why the write is short.
func (b *BufWriter) WriteString(s string) (int, error) {
	times := 0
	for 3 > times {
		if atomic.CompareAndSwapInt32(&b.locker, 0, 1) {
			defer atomic.CompareAndSwapInt32(&b.locker, 1, 0)
		} else {
			times++
			if nil != b.log {
				b.log.Brush(logger.Warn, "get write locker is fail,sleep 1 ms and try again.")
			}
			time.Sleep(1 * time.Millisecond)
			continue
		}
		nn := 0
		for len(s) > b.Available() && b.err == nil {
			n := copy(b.buf[b.offset:], s)
			b.offset += n
			nn += n
			s = s[n:]
			b.flush()
		}
		if b.err != nil {
			return nn, b.err
		}
		n := copy(b.buf[b.offset:], s)
		b.offset += n
		nn += n
		return nn, nil
	}
	if nil != b.log {
		b.log.Brush(logger.Error, "write is fail,try 3 time later.")
	}
	return 0, fmt.Errorf("buffer writer is fail.")
}

/*
// WriteByte writes a single byte.
func (b *Writer) WriteByte(c byte) error {
	if b.err != nil {
		return b.err
	}
	if b.Available() <= 0 && b.flush() != nil {
		return b.err
	}
	b.buf[b.n] = c
	b.n++
	return nil
}

// WriteRune writes a single Unicode code point, returning
// the number of bytes written and any error.
func (b *Writer) WriteRune(r rune) (size int, err error) {
	if r < utf8.RuneSelf {
		err = b.WriteByte(byte(r))
		if err != nil {
			return 0, err
		}
		return 1, nil
	}
	if b.err != nil {
		return 0, b.err
	}
	n := b.Available()
	if n < utf8.UTFMax {
		if b.flush(); b.err != nil {
			return 0, b.err
		}
		n = b.Available()
		if n < utf8.UTFMax {
			// Can only happen if buffer is silly small.
			return b.WriteString(string(r))
		}
	}
	size = utf8.EncodeRune(b.buf[b.n:], r)
	b.n += size
	return size, nil
}
*/

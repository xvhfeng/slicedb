package idxfile

import (
	"../logger"

	"time"
)

type IdxFileWriter struct {
	path string
	t    time.Time
	idx  int
	bufw BufWriter
	log  *logger.Log
}

type IdxFileReader struct {
	path string
	t    time.Time
	idx  int
	bufr BufReader
	log  *logger.Log
}

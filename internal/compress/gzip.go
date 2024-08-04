package compress

import (
	"compress/gzip"
	"io"
	"net/http"
)

type CompressWriter struct {
	w http.ResponseWriter
	gw *gzip.Writer
}

func NewCompressWrite (w http.ResponseWriter) *CompressWriter {
	return &CompressWriter{w: w, gw: gzip.NewWriter(w),}
}

func (c *CompressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *CompressWriter) Write(p []byte) (int, error) {
	return c.gw.Write(p)
}

func (c *CompressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *CompressWriter) Close() error {
	return c.gw.Close()
}


type CompressReader struct {
	r io.ReadCloser
	gr *gzip.Reader
}

func NewCompressReader(r io.ReadCloser) (*CompressReader, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &CompressReader{r: r, gr: gr}, nil
}

func (c *CompressReader) Read (p []byte) (int, error) {
	return c.gr.Read(p)
}

func (c *CompressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}

	return c.gr.Close()
}
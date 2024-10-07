package objectstream

import (
	"fmt"
	"io"
	"net/http"
)

type PutStream struct {
	writer *io.PipeWriter
	ch     chan error
}

func NewPutStream(server, object string) *PutStream {
	reader, writer := io.Pipe()
	ch := make(chan error)
	go func() {
		request, err := http.NewRequest(
			"PUT",
			"http://"+server+"/objects/"+object,
			reader,
		)
		if err != nil {
			fmt.Println("request dataServer error:", err)
			err = fmt.Errorf("request dataServer error: %d", err)
			ch <- err
		}

		client := http.Client{}
		res, err := client.Do(request)
		if err == nil && res.StatusCode != http.StatusOK {
			err = fmt.Errorf("dataServer return http code %d", res.StatusCode)
		}
		ch <- err

	}()

	return &PutStream{writer, ch}
}

func (w *PutStream) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

func (w *PutStream) Close() error {
	w.writer.Close()
	return <-w.ch
}

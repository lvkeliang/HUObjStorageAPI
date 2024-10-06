package objectstream

import (
	"fmt"
	"io"
	"net/http"
)

type GetStream struct {
	reader io.Reader
}

func newGetStream(url string) (*GetStream, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dataServer return http code %d", r.StatusCode)
	}

	return &GetStream{r.Body}, nil
}

func NewGetStream(server, objectName string) (*GetStream, error) {
	if server == "" || objectName == "" {
		return nil, fmt.Errorf("invalid server %s object %s", server, objectName)
	}
	return newGetStream("http://" + server + "/objects/" + objectName)
}

func (r *GetStream) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

package objectstream

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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

type TempPutStream struct {
	Server string
	Uuid   string
}

// 定义结构体来解析 POST 的 JSON 响应
type PostResponseData struct {
	Info string `json:"info"`
	UUID string `json:"uuid"`
}

func NewTempPutStream(server, hash string, size int64) (*TempPutStream, error) {

	request, err := http.NewRequest(
		"POST",
		"http://"+server+"/temp/"+hash,
		nil,
	)
	if err != nil {
		return nil, err
	}

	request.Header.Set("size", fmt.Sprintf("%d", size))

	client := http.Client{}

	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	// 读取响应体
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("读取响应体失败: %v", err)
		return nil, err
	}

	// 将 JSON 响应解码到结构体中
	var data PostResponseData
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("解析 JSON 失败: %v", err)
		return nil, err
	}

	return &TempPutStream{server, data.UUID}, nil
}

func (w *TempPutStream) Write(p []byte) (n int, err error) {
	request, err := http.NewRequest(
		"PATCH",
		"http://"+w.Server+"/temp/"+w.Uuid,
		strings.NewReader(string(p)),
	)
	if err != nil {
		return 0, err
	}
	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("dataServer return http code %d", res.StatusCode)
	}
	return len(p), nil
}

func (w *TempPutStream) Commit(good bool) {
	method := "DELETE"
	if good {
		method = "PUT"
	}

	request, _ := http.NewRequest(
		method,
		"http://"+w.Server+"/temp/"+w.Uuid, nil,
	)

	client := http.Client{}
	client.Do(request)
}

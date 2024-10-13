package es

import (
	"HUObjStorageAPI/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// 定义一个包装结构体来处理 "_source" 字段
var sourceResponse struct {
	Source Metadata `json:"_source"`
}

type Metadata struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
	Size    int64  `json:"size"`
	Hash    string `json:"hash"`
}

func getMetadata(name string, versionId int) (meta Metadata, err error) {
	addr := fmt.Sprintf("http://%s/metadata/_doc/%s_%d", config.Configs.Elasticsearch.EsServer, name, versionId)
	res, err := http.Get(addr)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to get %s_%d: %d", name, versionId, res.StatusCode)
		return
	}

	result, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	// 解码 JSON 到 sourceResponse 结构体
	err = json.Unmarshal(result, &sourceResponse)
	if err != nil {
		return
	}

	meta = sourceResponse.Source
	return
}

type hit struct {
	Source Metadata `json:"_source"`
}

type Total struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

type searchResult struct {
	Hits struct {
		Total Total `json:"total"`
		Hits  []hit `json:"hits"`
	} `json:"hits"`
}

func SearchLatestVersion(name string) (meta Metadata, err error) {
	query := url.Values{}
	query.Set("q", fmt.Sprintf("name:%s", name))
	query.Set("size", "1")
	query.Set("sort", "version:desc")

	addr := fmt.Sprintf("http://%s/metadata/_search?%s", config.Configs.Elasticsearch.EsServer, query.Encode())
	res, err := http.Get(addr)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to search latest metadata: %d", res.StatusCode)
		result, _ := io.ReadAll(res.Body)
		fmt.Println(string(result))
		return
	}

	result, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	var sr searchResult
	err = json.Unmarshal(result, &sr)
	if err != nil {
		return
	}

	if len(sr.Hits.Hits) != 0 {
		meta = sr.Hits.Hits[0].Source
	}
	return
}

func GetMetadata(name string, version int) (Metadata, error) {
	if version == 0 {
		return SearchLatestVersion(name)
	}
	return getMetadata(name, version)
}

func PutMetadata(name string, version int, size int64, hash string) error {
	doc := fmt.Sprintf(`{"name":"%s","version":%d,"size":%d,"hash":"%s"}`, name, version, size, hash)
	client := http.Client{}
	addr := fmt.Sprintf("http://%s/metadata/_doc/%s_%d?op_type=create", config.Configs.Elasticsearch.EsServer, name, version)
	request, _ := http.NewRequest("PUT", addr, strings.NewReader(doc))
	request.Header.Set("Content-Type", "application/json")

	res, err := client.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusConflict {
		return PutMetadata(name, version+1, size, hash)
	}
	if res.StatusCode != http.StatusCreated {
		result, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to put metadata: %d %s", res.StatusCode, string(result))
	}
	return nil
}

func AddVersion(name, hash string, size int64) error {
	version, err := SearchLatestVersion(name)
	if err != nil {
		return err
	}
	return PutMetadata(name, version.Version+1, size, hash)
}

func SearchAllVersions(name string, from, size int) ([]Metadata, error) {
	query := url.Values{}
	query.Set("sort", "name,version")
	query.Set("from", fmt.Sprintf("%d", from))
	query.Set("size", fmt.Sprintf("%d", size))

	if name != "" {
		query.Set("q", fmt.Sprintf("name:%s", name))
	}

	addr := fmt.Sprintf("http://%s/metadata/_search?%s", config.Configs.Elasticsearch.EsServer, query.Encode())
	res, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	metas := make([]Metadata, 0)
	result, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var sr searchResult
	err = json.Unmarshal(result, &sr)
	if err != nil {
		return nil, err
	}

	for _, hit := range sr.Hits.Hits {
		metas = append(metas, hit.Source)
	}

	return metas, nil
}

type Bucket struct {
	Key         string
	Doc_count   int
	Min_version struct {
		Value float32
	}
}

type aggregateResult struct {
	Aggregations struct {
		Group_by_name struct {
			Buckets []Bucket
		}
	}
}

func SearchVersionStatus(min_doc_count int) ([]Bucket, error) {
	client := http.Client{}
	addr := fmt.Sprintf("http://%s/metadata/_search", config.Configs.Elasticsearch.EsServer)

	body := fmt.Sprintf(`{"size":0,"aggs":{"group_by_name":{"terms":{"field":"name","min_doc_count":%d},"aggs":{"min_version":{"min":{"field":"version"}}}}}}`, min_doc_count)
	request, _ := http.NewRequest("GET", addr, strings.NewReader(body))
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	b, _ := io.ReadAll(res.Body)
	var ar aggregateResult
	json.Unmarshal(b, &ar)
	return ar.Aggregations.Group_by_name.Buckets, nil
}

func DelMetadata(name string, version int) {
	client := http.Client{}
	addr := fmt.Sprintf("http://%s/metadata/_doc/%s_%d", config.Configs.Elasticsearch.EsServer, name, version)
	request, _ := http.NewRequest("DELETE", addr, nil)
	client.Do(request)
}

func HasHash(hash string) (bool, error) {
	addr := fmt.Sprintf("http://%s/metadata/_search?q=hash:%s&size=0", config.Configs.Elasticsearch.EsServer, hash)
	res, err := http.Get(addr)
	if err != nil {
		return false, err
	}
	b, _ := io.ReadAll(res.Body)
	var sr searchResult
	json.Unmarshal(b, &sr)
	return sr.Hits.Total.Value != 0, nil
}

func SearchHashSize(hash string) (size int64, err error) {
	addr := fmt.Sprintf("http://%s/metadata/_search?q=hash:%s&size=1", config.Configs.Elasticsearch.EsServer, hash)
	res, err := http.Get(addr)
	if err != nil {
		return 0, err
	}
	if res.StatusCode != http.StatusOK {
		err := fmt.Errorf("failed to search hash size: %d", res.StatusCode)
		return 0, err
	}
	result, _ := io.ReadAll(res.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		size = sr.Hits.Hits[0].Source.Size
	}
	return size, nil
}

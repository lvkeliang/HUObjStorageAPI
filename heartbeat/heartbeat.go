package heartbeat

import (
	"HUObjStorageAPI/config"
	"HUObjStorageAPI/rabbitmq"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var dataServers = make(map[string]time.Time)
var mutex sync.Mutex

func ListenHeartbeat() {
	q := rabbitmq.New(config.Configs.Rabbitmq.RabbitmqServer)
	defer q.Close()

	q.Bind("apiServers")
	c := q.Consume()

	go removeExpiredDataServer()

	for msg := range c {
		dataServer, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			panic(err)
		}

		mutex.Lock()
		dataServers[dataServer] = time.Now()
		mutex.Unlock()
	}
}

func removeExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for s, t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}
}

func GetDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()
	ds := make([]string, 0)
	for s, _ := range dataServers {
		ds = append(ds, s)
	}
	return ds
}

func ChooseRandomDataServers(n int, exclude map[int]string) (dataServers []string) {
	candidates := make([]string, 0)
	reverseExcludeMap := make(map[string]int)

	for id, addr := range exclude {
		reverseExcludeMap[addr] = id
	}

	servers := GetDataServers()

	for i := range servers {
		s := servers[i]
		_, excluded := reverseExcludeMap[s]
		if !excluded {
			candidates = append(candidates, s)
		}
	}

	length := len(candidates)
	if length < n {
		return
	}
	p := rand.Perm(length)
	for i := 0; i < n; i++ {
		dataServers = append(dataServers, candidates[p[i]])
	}
	return
}

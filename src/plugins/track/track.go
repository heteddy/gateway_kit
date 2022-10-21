// @Author : detaohe
// @File   : track
// @Description:
// @Date   : 2022/10/20 20:11

package track

import (
	"bytes"
	"encoding/json"
	"gateway_kit/config"
	"gateway_kit/util"
	"gateway_kit/util/kafka"
	"github.com/Shopify/sarama"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	TrackHeader = "X-TRACK-PATH"
)

var teOnce sync.Once
var TeInstance *EventTracker

type Event struct {
	UserID    string                 `json:"user_id"`
	Uri       string                 `json:"uri"`
	Method    string                 `json:"method"`
	Path      string                 `json:"path"`
	Body      map[string]interface{} `json:"body"`
	CreatedAt string                 `json:"created_at"`
}

type EventTracker struct {
	sink   *kafka.KafkaSink
	topic  string
	eventC chan *Event
	stopC  chan struct{}
}

func NewEventTracker(c config.KafkaSinkConfig) *EventTracker {
	teOnce.Do(func() {
		addrs := strings.Split(c.Broker, ",")
		conf := sarama.NewConfig()
		conf.ClientID = c.ClientID
		TeInstance = &EventTracker{
			sink:   kafka.NewKafkaSink(addrs, conf, c.Topics),
			topic:  c.Topics,
			eventC: make(chan *Event, 100),
			stopC:  make(chan struct{}),
		}
	})
	return TeInstance
}

func getBody(req *http.Request) []byte {
	bodies, err := ioutil.ReadAll(req.Body)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodies))
	if err != nil {
		return nil
	}
	return bodies
}

func (tracker *EventTracker) Track(req *http.Request) {
	header := req.Header
	var path string
	var userID string
	// 首先判断是否登录，获取登录的header信息
	for k, v := range header {
		switch k {
		case TrackHeader:
			path = v[0]
		case util.JwtKeyUserID:
			userID = v[0]
		default:
		}
	}
	if userID == "" {
		return
	}
	//queries := req.URL.Query()
	body := getBody(req)
	bodyM := make(map[string]interface{})
	if err := json.Unmarshal(body, &bodyM); err != nil {

	}
	e := Event{
		UserID:    userID,
		Uri:       req.URL.Path + "?" + req.URL.RawQuery,
		Method:    req.Method,
		Path:      path,
		Body:      bodyM,
		CreatedAt: time.Now().Format("2006-04-02 15:04:05"),
	}
	tracker.eventC <- &e
}

func (tracker *EventTracker) runLoop() {
loop:
	for {
		select {
		case e, ok := <-tracker.eventC:
			if !ok {
				break loop
			}
			if s, err := json.Marshal(e); err != nil {

			} else {
				// todo: 写入kafka和mongo
				tracker.sink.In() <- &sarama.ProducerMessage{
					Topic: tracker.topic,
					//Key: sarama.StringEncoder(""), // 这里不填Key
					Value: sarama.StringEncoder(s),
				}
			}
		case <-tracker.stopC:
			break loop
		}
	}
}

func (tracker *EventTracker) Start() {
	go tracker.runLoop()
}

func (tracker *EventTracker) Stop() {
	close(tracker.stopC)
	close(tracker.eventC)
}

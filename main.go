package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"hangout-go/output"
	"time"
)

type ServicePv struct {
	Time     int64  `json:"time"`
	Logid    string `json:"logid"`
	Hostname string `json:"hostname"`
	Prelogid string `json:"prelogid"`
	Pageid   string `json:"pageid"`
	Sessid   string `json:"sessid"`
	Appid    uint64  `json:"appid"`
	Token    string `json:"token"`
	Ver      string `json:"ver"`
	Os       string `json:"os"`
	Ua       string `json:"ua"`
	Ip       string `json:"ip"`
	Xesid    string `json:"xesid"`
	Userid   string `json:"userid"`
	Devid    string `json:"devid"`
	// Data     []interface{} `json:"data"`
	Data []struct {
		Req  map[string]interface{} `json:"req"`
		Resp interface{}            `json:"resp"`
	}
}

func main() {
	logFile, err := os.OpenFile("/tmp/err_line.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Group.Topics.Whitelist = regexp.MustCompile(`shuaige.*`)
	// config.Group.Heartbeat.Interval = 1

	// brokers := []string{"10.97.14.111:9092", "10.97.14.112:9092", "10.97.14.112:9092"}
	brokers := []string{"10.99.1.151:19092", "10.99.1.148:19092", "10.99.1.48:19092"}
	// topics := []string{"flume-test-.*"}

	consumer, err := cluster.NewConsumer(brokers, "groupid", nil, config)
	if err != nil {
		panic(err)
	}
	// dbConn := output.NewDBConn("10.97.14.111", "9000", "test")
	dbConn := output.NewDBConn("10.99.1.151", "9001", "wangdazhuang")

	defer consumer.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		for err := range consumer.Errors() {
			fmt.Println("consumer error ", err)
		}
	}()

	go func() {
		for ntf := range consumer.Notifications() {
			fmt.Println("consumer not ", ntf)
		}
	}()

	for {
		select {
		case line, ok := <-consumer.Messages():
			if ok {
                fmt.Println(string(line.Value))
				msg, err := GetMsg(line.Value)
				if err != nil {
                    logFile.WriteString(err.Error() + "|")
					logFile.Write(line.Value)
					logFile.WriteString("\n")
					fmt.Println(err)
				} else {
					param, err := DealMsg(msg)
					if err != nil {
                        logFile.WriteString(err.Error() + "|")
						logFile.Write(line.Value)
						logFile.WriteString("\n")
						fmt.Println(err)
					} else {
						// fmt.Println(param)
						// fmt.Println(1)
						if err = dbConn.Insert("xes_service_pv", param); err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		case <-signals:
			return
		}
	}

}

func GetMsg(line []byte) ([]byte, error) {
	var msg []byte
	index := strings.Index(string(line), "{\"time")
	if index == -1 {
		return nil, errors.New("line error")
	}

	msg = line[index:]

	return msg, nil
}

func DealMsg(log []byte) (map[string]interface{}, error) {
	content := &ServicePv{}
	var param map[string]interface{}
	err := json.Unmarshal(log, &content)
	if err != nil {
		return nil, err
	} else {
		param = make(map[string]interface{})
		param["url"] = ""
		param["req"] = ""
		param["resp"] = ""

		for _, val := range content.Data {
			if val.Req != nil {
				if url, ok := val.Req["url"].(string); ok {
					param["url"] = url
				} else {
					return nil, errors.New("get url fail")
				}

				if req, err := json.Marshal(content.Data[0].Req); err != nil {
					return nil, err
				} else {
					param["req"] = string(req)
				}
			}

			if val.Resp != nil {
				if resp, err := json.Marshal(content.Data[1].Resp); err != nil {
					return nil, err
				} else {
					param["resp"] = string(resp)
				}
			}
		}

		param["date"] = time.Now().Format("2006-01-02")
		param["time"] = content.Time
		param["logid"] = content.Logid
		param["hostname"] = content.Hostname
		param["prelogid"] = content.Prelogid
		param["pageid"] = content.Pageid
		param["sessid"] = content.Sessid
		param["appid"] = content.Appid
		param["token"] = content.Token
		param["ver"] = content.Ver
		param["os"] = content.Os
		param["ua"] = content.Ua
		param["xesid"] = content.Xesid
		param["userid"] = content.Userid
		param["devid"] = content.Devid

		return param, nil
	}
}

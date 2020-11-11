package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func FixPage() {
	fmt.Println("我执行了队列任务")
	conn := GetMq()
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"page", false, false, false, false, nil)
	ch.Qos(1, 0, false)
	msg, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	for m := range msg {
		data := make(map[string]string)
		json.Unmarshal(m.Body, data)
		//http://www.diyibanzhu6.me/
		re, _ := http.Get("http://localhost:8090/get")
		d, _ := ioutil.ReadAll(re.Body)
		if re != nil {
			var i IpInfo
			json.Unmarshal(d, &i)
			GetTitle("http://"+i.IP+":"+strconv.Itoa(i.Port), "http://www.diyibanzhu6.me/"+data["href"])
		}
		GetTitle("", "http://www.diyibanzhu6.me/"+data["href"])
		fmt.Println("队列消息", data["href"])
		time.Sleep(time.Second)
	}
	fmt.Println("别告诉我你退出了")
}

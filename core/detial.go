package core

import (
	"encoding/json"
	"fmt"
	"github.com/AceDarkknight/GoProxyCollector/server"
	"strconv"
	"time"
)

func GetDetial(){
	conn:=GetMq()
	ch, err := conn.Channel()
	if err!=nil{
		panic(err)
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"novel",false,false,false,false,nil)
	ch.Qos(1,0,false)
	msg,err:=ch.Consume(q.Name,"",true,false,false,false,nil)
	if err!=nil{
		panic(err)
	}
	for m:=range msg{
		data:=make(map[string]string)
		json.Unmarshal(m.Body,data)
		//http://www.diyibanzhu6.me/
		re:=server.GetServer()
		if re!=nil{
			var i IpInfo
			json.Unmarshal(re,&i)
			GetTitle("http://"+i.IP+":"+strconv.Itoa(i.Port),"http://www.diyibanzhu6.me/"+data["href"])
		}
		GetTitle("","http://www.diyibanzhu6.me/"+data["href"])
		fmt.Println("队列消息",data["href"])
		time.Sleep(time.Second)
	}
	fmt.Println("别告诉我你退出了")
}

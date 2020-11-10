package core

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/streadway/amqp"
	"time"
)
var conn *amqp.Connection
func GetMq()*amqp.Connection{
	defer log.Flush()
	if conn==nil{
		conn,err:=amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err!=nil{
			panic(err)
			log.Info("连接rabbitmq失败",err)
		}
		return conn
	}
	return conn
}
func AddNovelPage(by []byte,exchange string){
	fmt.Println("我tm执行了吗")
	c:=GetMq()
	ch, err := c.Channel()
	if err!=nil{
		log.Info(err)
	}
	defer ch.Close()
	//创建交换机
	q,err:=ch.QueueDeclare(exchange,false,false,false,false,nil)
	//默认交换机
	err=ch.Publish("",q.Name,false,false,amqp.Publishing{
		Headers:         nil,
		ContentType:     "text/plain",
		ContentEncoding: "",
		DeliveryMode:    0,
		Priority:        0,
		CorrelationId:   "",
		ReplyTo:         "",
		Expiration:      "",
		MessageId:       "",
		Timestamp:       time.Time{},
		Type:            "",
		UserId:          "",
		AppId:           "",
		Body:            by,
	})
	if err!=nil{
		panic(err)
		log.Info(err)
	}
}

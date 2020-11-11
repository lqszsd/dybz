package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"novel/conf"
	"strconv"
)

func GetContent() {
	db, _ := conf.NewDb()
	var list = []Novel{}
	//"/0/2/7.html"
	db.Table("lq_novel").Where("link!=?", "").Limit(100).Find(&list)
	for _, g := range list {
		//http://www.diyibanzhu6.me/+g.link
		re, _ := http.Get("http://localhost:8090/get")
		d, _ := ioutil.ReadAll(re.Body)
		if re != nil {
			var i IpInfo
			json.Unmarshal(d, &i)
			GetPage("http://"+i.IP+":"+strconv.Itoa(i.Port), "http://www.diyibanzhu6.me"+g.Link)
		}
		GetPage("", "http://www.diyibanzhu6.me"+g.Link)
	}
	fmt.Println("看看数据库的值", list)
}

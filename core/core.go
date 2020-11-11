package core

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
	"log"
	"math/rand"
	"novel/conf"
	"strconv"
	"strings"
)

type Novel struct {
	Id        int64  `json:"id"`
	Title     string `json:"title"`
	ChapterId int    `json:"chapter_id"`
	ListId    int64  `json:"list_id"`
	Link      string `json:"link"`
	Content   string `json:"content"`
}
type List struct {
	Id    int64  `json:"id"`
	Link  string `json:"link"`
	Title string `json:"title"`
}
type IpInfo struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Location string `json:"location"`
	Source   string `json:"source"`
	Speed    int    `json:"speed"`
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func GetTitle(str string, u string) {
	fmt.Println("采集数据开始")
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("body > div.container > div.mod.page > div.pagelistbox > a.nextPage", func(e *colly.HTMLElement) {
		var data map[string]interface{}
		data = make(map[string]interface{})
		//fmt.Println(string(e.Response.Body))
		fmt.Println("我找不到a标签的值", e.Attr("href"))
		data["type"] = "1"
		data["href"] = e.Attr("href")
		j, _ := json.Marshal(data)
		AddNovelPage(j, "page")
	})
	c.OnHTML("body > div.container > div.mod.block.book-all-list > div.bd > ul > li", func(e *colly.HTMLElement) {
		//fmt.Println(string(e.Response.Body))
		u := e.ChildAttrs("a", "href")
		fmt.Println("我找得到链接吗", u)
		for _, s := range u {
			if !strings.Contains(s, "author") {
				var data map[string]interface{}
				data = make(map[string]interface{})
				data["type"] = "2"
				data["href"] = s
				j, _ := json.Marshal(data)
				AddNovelPage(j, "novel")
			}
		}
	})
	if str != "" {
		rp, err := proxy.RoundRobinProxySwitcher(str)
		if err != nil {
			log.Fatal(err)
		}
		c.SetProxyFunc(rp)
	}
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
	})

	c.Visit(u)
}
func GetList(str string, u string) {
	fmt.Println("采集数据开始")
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("body", func(e *colly.HTMLElement) {
		//判断当天章节是否已经存在数据库
		db, err := conf.NewDb()
		if err != nil {
			fmt.Println("这是错误信息", err)
		}
		var l = &List{}
		db.Debug().Table("lq_list").Where("link", u).Limit(1).Find(&l)
		if l.Id == 0 {
			l.Link = u
			fmt.Println("我找不到title", e.ChildText("div.container > div.mod.detail > div.bd.column-2 > div.right > h1"), 111111)
			l.Title = e.ChildText("div.container > div.mod.detail > div.bd.column-2 > div.right > h1")
			db.Debug().Table("lq_list").Create(l)
			//获取章节信息
			re := e.ChildAttrs("body > div.container > div:nth-child(5) > div.bd > ul > li>a", "href")
			if len(re) > 0 {
				for j, s := range re {
					var n = &Novel{}
					n.ListId = l.Id
					n.Link = s
					n.Title = e.ChildText("body > div.container > div:nth-child(5) > div.bd > ul > li>a:nth(" + strconv.Itoa(j+1) + ")")
					n.ChapterId = j
					db.Debug().Table("lq_novel").Create(&n)
				}
			}
		} else {
			//不做任何处理
		}
	})
	if str != "" {
		rp, err := proxy.RoundRobinProxySwitcher(str)
		if err != nil {
			log.Fatal(err)
		}
		c.SetProxyFunc(rp)
	}
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
	})

	c.Visit(u)
}
func GetPage(str string, u string) {
	fmt.Println("采集数据开始")
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("#ChapterView > div.bd > div", func(e *colly.HTMLElement) {
		//判断当天章节是否已经存在数据库
		db, err := conf.NewDb()
		if err != nil {
			fmt.Println("这是错误信息", err)
		}
		var re = Novel{}
		db.Where("link=?", u).First(&re)
		re.Content += re.Content + e.Text
	})
	if str != "" {
		rp, err := proxy.RoundRobinProxySwitcher(str)
		if err != nil {
			log.Fatal(err)
		}
		c.SetProxyFunc(rp)
	}
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
	})

	c.Visit(u)
}

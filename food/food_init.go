package food

import (
	"log"
)

var cli Client

func init() {

}
func Init() {
	log.Println("food start")
	conMongo() // 连接数据库
	//defer saveMeta() // 保存元信息
	cli = NewClient(2) // 创建全局客户端
	//go IpPoolStart() // 启动代理ip池
	//loadMeta()
	//go saveMetaInterval()
	//go fetchCats()              // 开始抓取分类
	//go searchKeywordLoop() // 开始搜索关键词
	go getDetailLoop() // 获取详情
	select {}
}

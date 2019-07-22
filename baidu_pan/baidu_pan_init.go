package baidu_pan

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Init() {
	log.Println("hello world")
	err := addTask("1",
		".",
		"http://pcs.dcdn.baidu.com/file/e47f77544ab7145da47a504ab49057e0?bkt=p3-0000a15f7e48b3796d61f707db5b3d321a01&xcode=197eef6f0fcbc0dd090d26bf39a511e2dd2dcb2e6d7f99b7da8f365422d4ea3a4ec35904a24b1de9dfceed65911c34a19717ec4418c70769&fid=2737778138-250528-1033811237124489&time=1541751292&sign=FDTAXGERQBHSKa-DCb740ccc5511e5e8fedcff06b081203-AuUrYsIbv9WFndXk4Ok0uTBCHTY%3D&to=z2&size=1929960833&sta_dx=1929960833&sta_cs=32357&sta_ft=mkv&sta_ct=7&sta_mt=0&fm2=MH%2CNanjing02%2CAnywhere%2C%2Chunan%2Cct&ctime=1420773040&mtime=1541732236&resv0=cdnback&resv1=0&vuk=2737778138&iv=2&htype=&newver=1&newfm=1&secfm=1&flow_ver=3&pkey=0000a15f7e48b3796d61f707db5b3d321a01&expires=8h&rt=pr&r=522045895&mlogid=7251861952334001417&vbdid=1745629073&fin=fate+stay+night%E5%89%A7%E5%9C%BA%E7%89%88.mkv&rtype=1&dp-logid=7251861952334001417&dp-callid=0.1.1&tsl=0&csl=0&csign=Ql1MyhBqm0E3t1Rkdn%2FsUfwU6Hs%3D&so=1&ut=1&uter=-1&serv=1&uc=3252767510&ti=76168191086d6f291a05685e8e79b9cd4cf2231e3f246fe6&by=themis")
	log.Println(err)
	select {}
}

var httpClient = http.Client{
	Timeout: 15 * time.Second,
}

var settings = struct {
	ua string
	segSize,
	coNum int // 每个任务协程数量
}{
	ua:      "netdisk;4.6.2.0;PC;PC-Windows;10.0.10240;WindowsBaiduYunGuanJia",
	segSize: 1024 * 1024 * 4, // 4 MB 为一个分段
	coNum:   16,
}

type tTask struct {
	id,
	filename,
	tempName,
	localPath,
	tempPath,
	link string
	url     *url.URL
	size    int
	segLock *sync.RWMutex
	state   downloadState
}

type seg struct {
	start, len  int
	finish      bool
	distributed bool
}
type downloadState struct {
	distributedSeg bool
	dataOffset     int
	seg            []*seg
}

var taskMap = new(sync.Map)

func addTask(id, localPath, link string) error {
	err, size, filename, url := checkLink(link)
	log.Println(size, filename)
	if err != nil {
		return errors.Wrap(err, "获取文件信息失败")
	}
	p, err := filepath.Abs(localPath)
	if err != nil {
		return errors.Wrap(err, "下载路径错误")
	}
	t := &tTask{
		id:        id,
		link:      link,
		url:       url,
		size:      size,
		filename:  filename,
		tempName:  "." + filename + ".tmp",
		localPath: p,
		segLock:   new(sync.RWMutex),
	}
	t.state.dataOffset = 1 * 1024 // 1KB头部信息
	taskMap.Store(id, t)
	go handleTask(t)
	return nil
}

func checkLink(link string) (err error, size int, filename string, u *url.URL) {
	u, err = url.Parse(link)
	if err != nil {
		err = errors.Wrap(err, "url 解析错误")
		return
	}
	var res *http.Response
	res, err = httpClient.Do(&http.Request{
		URL: u,
		Header: map[string][]string{
			"User-Agent": {settings.ua},
		},
	})
	if err != nil || res.StatusCode != 200 {
		err = errors.New("访问链接错误, 链接可能已失效")
		return
	}
	log.Println(res.Header)
	if cd, ok := res.Header["Content-Disposition"]; ok && len(cd) > 0 {
		if strings.IndexAny(cd[0], "attachment;filename=") != 0 {
			err = errors.New("不是文件链接")
			return
		}
		filename = strings.Trim(cd[0][len("attachment;filename="):], "\"")
	} else {
		err = errors.New("不是文件链接")
		return
	}
	if cl, ok := res.Header["Content-Length"]; ok && len(cl) > 0 {
		size, err = strconv.Atoi(cl[0])
	} else {
		err = errors.New("不是文件链接")
		return
	}
	return
}

func handleTask(task *tTask) error {
	if !(*task).state.distributedSeg {
		// 段分配
		for start := 0; start < task.size; start += settings.segSize {
			lg := settings.segSize
			if task.size-start < lg {
				lg = 0
			}
			(*task).state.seg = append((*task).state.seg, &seg{start: start, len: lg, finish: false})
		}
	}
	// 临时文件句柄
	fullName := task.localPath + string(filepath.Separator) + task.tempName
	f, err := os.OpenFile(fullName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return errors.Wrap(err, "无法创建文件: "+fullName)
	}

	type mini struct {
		task   *tTask
		seg    *seg
		buffer *bytes.Buffer
		others []io.Writer
		err    error
	}
	miniTaskCh := make(chan mini, 1)
	miniTaskFinishCh := make(chan mini, 1)
	wg := new(sync.WaitGroup)
	// 分配任务的协程
	go func() {
		for {
			// 遍历seg状态得到未下载的
			(*task).segLock.RLock()
			unlocked := false
			allDone := true
			for _, seg := range (*task).state.seg {
				if !seg.finish {
					allDone = false
				}
				if !seg.finish && !seg.distributed {
					(*task).segLock.RUnlock()
					unlocked = true
					seg.distributed = true
					miniTaskCh <- mini{task: task, seg: seg, buffer: bytes.NewBuffer([]byte{})}
					break
				}
			}
			if !unlocked {
				(*task).segLock.RUnlock()
			}
			if allDone {
				close(miniTaskCh)
				break
			}
		}
	}()
	// 接收数据写入磁盘的协程
	go func() {
		for mini := range miniTaskFinishCh {
			func() {
				defer func() { mini.seg.distributed = false }()
				if mini.err != nil {
					return
				}
				_, err := f.Seek(int64(mini.seg.start), io.SeekStart)
				if err != nil {
					log.Println(err)
					return
				}
				_, err = mini.buffer.WriteTo(f)
				if err != nil {
					log.Println(err)
					return
				}
				mini.seg.finish = true
				log.Println("mini finish", *mini.seg)
			}()
		}
		log.Println("下载完成")
	}()
	// 下载数据的协程, 开启多个
	for i := 0; i < settings.coNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for mini := range miniTaskCh {
				log.Println("new mini", *mini.seg)
				rg := fmt.Sprintf("bytes=%d-%d", mini.seg.start, mini.seg.start+mini.seg.len-1)
				if mini.seg.len == 0 {
					rg = fmt.Sprintf("bytes=%d-", mini.seg.start)
				}
				req := &http.Request{
					Method: "GET",
					URL:    mini.task.url,
					Header: map[string][]string{
						"Range":      {rg},
						"User-Agent": {settings.ua},
					},
				}
				res, err := httpClient.Do(req)
				if err != nil {
					log.Println(err)
					mini.err = err
					miniTaskFinishCh <- mini
					continue
				}
				l, err := mini.buffer.ReadFrom(res.Body)
				if err != nil {
					log.Println(err)
					mini.err = err
					miniTaskFinishCh <- mini
					continue
				}
				log.Println("小段下载完毕", l, err, res.Header)
				miniTaskFinishCh <- mini
			}
		}()
	}
	// 下载协程全部退出, 关闭任务完成通道, 以便接收协程退出
	go func() {
		wg.Wait()
		close(miniTaskFinishCh)
	}()
	return nil
}

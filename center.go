package guy_rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	nt "github.com/shirou/gopsutil/v3/net"
	"guy-rpc/server"
	"log"
	"net/http"
	"time"
)

//func StartCenter()  {
//	center.StartCenter()
//}
//
//func HandleHTTP()  {
//	register.HandleHTTP()
//}
//
//func Heartbeat(registry, addr string, duration time.Duration)  {
//	register.Heartbeat(registry,addr,duration)
//}
//

var State server.Service
var CenterIP string

func Heartbeat(ServiceIp string, Name string, CenterIp string) {
	State.Ip = ServiceIp
	State.Name = Name
	CenterIP = CenterIp

	c := cron.New()
	spec := "*/5 * * * * *" // 每隔5s执行一次，cron格式（秒，分，时，天，月，周）
	// 添加一个任务
	err := c.AddFunc(spec, SendHeart)
	if err != nil {
		log.Println(err)
		return
	}

	c.Start()

}

func SendHeart() {
	v, _ := mem.VirtualMemory()
	//fmt.Println(v.Total, v.UsedPercent, v.Used, v.Free)

	State.Info.MemAll = v.Total
	State.Info.MemFree = v.Free
	State.Info.MemUsed = v.Used
	State.Info.MemUsedPercent = v.UsedPercent

	c1, _ := cpu.Percent(time.Duration(time.Second), false)
	//fmt.Println(c1)

	State.Info.CpuPercent = c1[0]

	d, _ := disk.Usage("/")
	//fmt.Println(d.Total, d.Used, d.UsedPercent)

	State.Info.DiskTotal = d.Total
	State.Info.DiskUsedPercent = d.UsedPercent
	State.Info.DiskUsed = d.Used

	info, _ := nt.IOCounters(false)
	//fmt.Println(info[0].BytesSent, info[0].BytesRecv, info[0])

	State.Info.IOBytesSent = info[0].BytesSent
	State.Info.IOBytesRev = info[0].BytesRecv

	State.RefreshTime = time.Now()

	fmt.Println(State)

	js, err := json.Marshal(State)
	if err != nil {
		log.Println("json marshal err : ", err)
		return
	}

	req, err := http.NewRequest("POST", "http://"+CenterIP+"/service", bytes.NewBuffer(js))

	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println(err)
		panic(err)
	}
	defer response.Body.Close()
	return
}

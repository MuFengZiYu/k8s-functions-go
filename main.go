package main

import (
	"Calculate/pkg/setting"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var hostname = ""       // hostname used for auto-configure dns
var wait sync.WaitGroup //  syn thread is similiar to countdownlatch
var simMemory map[int]int

//*************** multi-thread calculate (SUM) *************************
func Calculate(times int) int {
	var result int

	cpunum, err := cpu.Counts(true)
	memory, _ := mem.VirtualMemory()
	memorysize := memory.Available / 1024 / 1024
	//cpunum = 1
	if err != nil {
		fmt.Printf("出错了，%s", err)
	}
	fmt.Printf("cpunum : %d\n mem : %d \n ", cpunum, memorysize)

	//××××××××× 模拟内存存取 (未实现)×××××××××××××

	//初始化
	//for i := 0; i < int(memorysize); i++ {
	//	simMemory[i] = 1
	//}

	//××××××××× 模拟CPU计算×××××××××××××

	times = times / cpunum
	c := make(chan int, cpunum)
	for i := 0; i < cpunum; i++ {
		go Count(times, c)
		wait.Add(1)
	}
	time.Sleep(1) // at least wait for one milli otherwise it will occur to error
	wait.Wait()
	for i := 0; i < cpunum; i++ {
		result += <-c
	}

	return result

}

func Count(times int, c chan int) {
	sum := 0
	times *= 10000
	for i := 0; i < times; i++ {
		sum += i
	}
	c <- sum
	defer wait.Done()
	//close(c)
}

// ********* test-fibonacci ***************
func fibonacci(times int) int {
	if times == 1 || times == 0 {
		return 1
	}

	return Calculate(times-1) + Calculate(times-2)
}

func handlefunc(c *gin.Context) {
	path := c.Query("path")
	value := c.Query("value")
	//RandMemory := c.Query("mem")

	result := ""
	result += request(path, value)
	c.Writer.WriteString(result)
}

func request(path, value string) string {
	result := ""
	if strings.Contains(path, "-") {
		client := http.DefaultClient
		pathIndex := strings.Index(path, "-")
		path = path[pathIndex+1:]
		valueIndex := strings.Index(value, "-")
		valueInt, err := strconv.Atoi(value[0:valueIndex])
		value = value[valueIndex+1:]
		urlStr := "http://" + hostname + "/calculate?path=" + path + "&value=" + value
		resp, err := client.Get(urlStr)
		//resp, err := http.Get("http://localhost:8080/calculate?path=1&value=2")
		if err != nil && resp.StatusCode != 200 {
			fmt.Printf("error : %s \n", err)
			client.CloseIdleConnections()
			return "error ..\n"
		}
		len := resp.ContentLength
		buf := make([]byte, len)
		resp.Body.Read(buf)
		result += string(buf[0:])
		if err != nil {
			fmt.Printf("format error1 : %s \n", err)
			return "format error \n"
		}

		result += strconv.Itoa(Calculate(valueInt))
		fmt.Printf("result1: -- %s --\n", result)
		return result + "\n"
	} else {
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			fmt.Printf("fomat error2 : %s \n", err)
			return "error \n"
		}
		result += strconv.Itoa(Calculate(valueInt))
		fmt.Printf("result2: -- %s --\n", result)
		return result + "\n"
	}

}

func main() {
	setting.InitSetting()
	hostname = setting.Setting.Redis.Host
	router := gin.Default()
	router.GET("calculate", handlefunc)
	err := router.Run(":8080")
	if err != nil {
		return
	}

}

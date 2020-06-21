package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
	"unsafe"
	"log"
)

func Ips() (map[string]string, error) {

	ips :=  make(map[string]string)
	log.Printf("get all card info before time")
	interfaces, err := net.Interfaces() // 获取网卡所有信息
	log.Printf("get all card info after time")
	if err != nil {
		return nil, err
	}
	for _, netCard := range interfaces {
		flags := netCard.Flags.String() // 获取网卡属性， "up" "broadcast"
		ipList, _ := netCard.Addrs()   // 获取网卡所有的IP xxxx/xx， ipv4 ipv6
		//获取up、非lo、有ip的网卡
		if strings.Contains(flags, "up") && strings.Contains(flags, "broadcast") && len(ipList) != 0 {
			for _, i := range ipList{
				if ipnet, ok := i.(*net.IPNet); ok {//通过 类型断言 来得到这个类型（net.IPNe 类型）的实例
				    // 获取ipv4
			    	if ipnet.IP.To4() != nil {
			    		ipLocal := ipnet.IP.String()
					    ips[netCard.Name] = ipLocal
				    }
				}
			}
		}
		if err != nil {
			return nil, err
		}
	}
	return ips, nil
}




func Shell(paramList []string ) (string, error) {


    //paramT := strings.Join(paramList, " ")

	cmd := exec.Command(paramList[0], paramList[1:len(paramList)]...)
	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stdout := "Error:can not obtain stdout pipe for command,"
		return  stdout, err
	}

	//执行命令
	if err := cmd.Start(); err != nil {
		stdout := "Error:The command is err,"
		return  stdout, err
	}

	//读取所有输出
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		stdout := "ReadAll Stdout:"
		return  stdout, err
	}

	if err := cmd.Wait(); err != nil {
		stdout := "wait:"
		return  stdout, err
	}

	return *(*string)(unsafe.Pointer(&bytes)), err
}



func main() {

	bondPath := "/proc/net/bonding/bond1"
	var bondState []string
	_, err := os.Stat(bondPath)
	if err != nil {
		bondState = append(bondState, "cat /proc/net/bonding/bond1")
	} else {
		bondState = append(bondState, "cat /proc/net/bonding/bond0")
	}


	var ch chan int
	//定时任务
	ticker := time.NewTicker(time.Second * 1)
	go func() {
		for range ticker.C {
			bondinfo, _ := Shell(bondState)
			count := strings.Count(bondinfo, "Interface")
			if count != 2 {
				log.Printf("get ip before time")
				ipList, _ := Ips()
				log.Printf("get ip after time")
				for device, ip := range ipList{
					arpSend := []string{"arping", "-A", ip, "-c", "1", "-I", device}
					stdout, err := Shell(arpSend)
					fmt.Println(stdout, err)
				}
			}
		}
		ch <- 1
	}()
	<-ch
}
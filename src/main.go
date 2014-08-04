package main

import (
	"bufio"
	"flag"
	"fmt"
	z "github.com/nutzam/zgo"
	serial "github.com/tarm/goserial"
	"io"
	"log"
	"runtime"
	"strings"
	"time"
)

// 设备
var device = flag.String("dev", "", "device address")

// 查看支持类型
var show = flag.Bool("show", false, "show phone book type")

// 读取电话薄
var read = flag.String("read", "", "ON 1")

// 写入电话薄
var write = flag.String("write", "", "ON 1 18600000000 Emergency")

// 主
func main() {

	// 解析程序参数
	flag.Parse()

	// 设置CPU核心数量
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 设置日志的结构
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.Lmicroseconds)

	// 判断是否指定设备
	if z.IsBlank(*device) {
		log.Panic("not set device address")
	}

	// 查看支持电话簿
	if *show {
		v, e := Com(*device, "AT+CPBS=?")
		if e != nil {
			log.Panic(e)
		}
		log.Println(v)
	}

	// 读取
	if !z.IsBlank(*read) {
		if cmd := strings.Fields(*read); len(cmd) == 2 {
			Com(*device, fmt.Sprintf("AT+CPBS=%s", cmd[0]))
			v, e := Com(*device, fmt.Sprintf("AT+CPBR=%s", cmd[1]))
			if e != nil {
				log.Panic(e)
			}
			log.Println(v)
		} else {
			log.Panic("parameter exception ...")
		}
	}

	// 写入
	if !z.IsBlank(*write) {
		if cmd := strings.Fields(*write); len(cmd) == 4 {
			Com(*device, fmt.Sprintf("AT+CPBS=%s", cmd[0]))
			v, e := Com(*device, fmt.Sprintf("AT+CPBW=%s,\"%s\",129,\"%s\")", cmd[1], cmd[2], cmd[3]))
			if e != nil {
				log.Panic(e)
			}
			log.Println(v)
		} else {
			log.Panic("parameter exception ...")
		}
	}

}

// Com 接口
func Com(dev string, data string) ([]string, error) {

	c := &serial.Config{Name: dev, Baud: 115200}
	s, e := serial.OpenPort(c)
	if e != nil {
		return make([]string, 0), e
	}
	defer s.Close()

	wd := bufio.NewWriter(s)
	wd.Write([]byte(data + "\r"))
	wd.Flush()

	time.Sleep(1 * time.Second)

	content := make([]string, 0)

	rd := bufio.NewReader(s)
	for {

		line, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		content = append(content, Trim(line))
		if strings.Contains(line, "OK") || strings.Contains(line, "ERROR") {
			break
		}

	}

	return content, nil

}

// 去掉一个字符串左右的空白串，即（0x00 - 0x20 之内的字符均为空白字符）
// 与strings.TrimSpace功能一致
func Trim(s string) string {
	size := len(s)
	if size <= 0 {
		return s
	}
	l := 0
	for ; l < size; l++ {
		b := s[l]
		if !IsSpace(b) {
			break
		}
	}
	r := size - 1
	for ; r >= l; r-- {
		b := s[r]
		if !IsSpace(b) {
			break
		}
	}
	return string(s[l : r+1])
}

// 是不是空字符
func IsSpace(c byte) bool {
	if c >= 0x00 && c <= 0x20 {
		return true
	}
	return false
}

package main

import (
	"bufio"
	serial "github.com/tarm/goserial"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

// 主
func main() {

	// 设置CPU核心数量
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 设置日志的结构
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.Lmicroseconds)

	v, e := Com(os.Args[1], "AT+CFUN=1,1")
	if e != nil {
	}
	log.Println(v)

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
		if strings.Contains(line, "OK") {
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

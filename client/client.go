package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"strconv"
	"../protocol"
	"bytes"
	"bufio"
)

//客户端状态
var clientStatus         = "close"

func send(conn net.Conn) {
	running := true
	reader := bufio.NewReader(os.Stdin)
	for running {
		data, _, _ := reader.ReadLine()
		command := string(data)
		if command == "stop" {
			running = false
		}
		conn.Write(protocol.Enpack([]byte("command:" +command)))
		fmt.Println("command", command)
	}
	fmt.Println("send over")
	//defer conn.Close()
}

func GetSession() string{
	gs1:=time.Now().Unix()
	gs2:=strconv.FormatInt(gs1,10)
	return gs2
}

/**
获取主机名
 */
func GetHostName()string{
	host, err := os.Hostname()
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		fmt.Println(host)
	}
	return  host;
}

/**
接受服务端操作指令
 */
func HandleInstruct(instruct string ,conn net.Conn)string{
	if instruct == "open"{
		clientStatus = "open";
	}
	if instruct == "close"{
		clientStatus = "close";
	}
	if instruct == "status"{
		conn.Write(protocol.Enpack([]byte("HostName:"+GetHostName()+";IP:" +GetIPAddress()+";status:"+clientStatus)))
	}
	return "";
}

/**
获取本机IP地址
 */
func GetIPAddress() string {
	var buffer bytes.Buffer
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				buffer.Write([]byte(ipnet.IP.String()));
			}
		}
	}
	return buffer.String();
}

/**
心跳检测
 */
func checkHeartBeat(conn net.Conn)  {
	go func() {
		for  {
			conn.Write(protocol.Enpack([]byte("心跳")))
			time.Sleep(time.Second * time.Duration(10))
		}
	}()
}

func main() {
	server := "localhost:6060"
	conn, err := net.DialTimeout("tcp",server,2 * time.Second);
	if err != nil {
		fmt.Println("连接超时")
		os.Exit(1)
	}
	fmt.Println("connect success")
	conn.Write(protocol.Enpack([]byte(GetHostName() +";"+"注册")))
	checkHeartBeat(conn)
	handleConnection(conn)
	send(conn)

}

func handleConnection(conn net.Conn) {
	go func() {
		for  {
			fmt.Println("接收服务端指令")
			// 缓冲区，存储被截断的数据
			tmpBuffer := make([]byte, 0)
			//接收解包
			readerChannel := make(chan []byte, 16)
			go reader(readerChannel,conn)
			buffer := make([]byte, 1024)
			for {
				n, err := conn.Read(buffer)
				if err != nil {
					Log(conn.RemoteAddr().String(), " connection error: ", err)
					return
				}
				tmpBuffer = protocol.Depack(append(tmpBuffer, buffer[:n]...), readerChannel)
			}
		}
	}()

	//defer conn.Close()
}

func reader(readerChannel chan []byte,conn net.Conn) {
	for {
		select {
		case data := <-readerChannel:
			Log(string(data))
			HandleInstruct(string(data),conn)
		}
	}
}

func Log(v ...interface{}) {
	fmt.Println(v...)
}


package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"strconv"
	"../protocol"
	"bytes"
)

func send(conn net.Conn) {
	//for i := 0; i < 100; i++ {
	//	session:=GetSession()
	//	words := "{\"ID\":"+ strconv.Itoa(i) +"\",\"Session\":"+session +"2015073109532345\",\"Meta\":\"golang\",\"Content\":\"message\"}"
	//	conn.Write(protocol.Enpack([]byte(words)))
	//}
	conn.Write(protocol.Enpack([]byte(GetHostName() +";"+"注册")))
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
		fmt.Printf("%s", host)
	}
	return  host;
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
	//ticker := time.NewTicker(time.Minute * 10)
	//go func() {
	//	for _ = range ticker.C {
	//		conn.Write([]byte("心跳"))
	//		fmt.Printf("ticked at %v", time.Now())
	//	}
	//}()
	conn.Write([]byte("心跳"))
}

func main() {
	server := "localhost:6060"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	fmt.Println("connect success")
	send(conn)
	//defer func(){ // 必须要先声明defer，否则不能捕获到panic异常
	//	fmt.Println("send fail")
	//	if err:=recover();err!=nil{
	//		fmt.Println(err) // 这里的err其实就是panic传入的内容
	//		conn.Close();
	//		os.Exit(1);
	//	}
	//	fmt.Println("d")
	//
	//}()
	checkHeartBeat(conn)

}
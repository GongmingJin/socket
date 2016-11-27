package main
import (
	"fmt"
	"net"
	"os"
	"../protocol"
	"bufio"
)

func main() {
	netListen, err := net.Listen("tcp", "localhost:6060")
	CheckError(err)
	defer netListen.Close()
	Log("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		Log(conn.RemoteAddr().String(), " tcp connect success")
		go handleConnection(conn)
		AcceptExternalStruct(conn);
	}
}

func handleConnection(conn net.Conn) {
	// 缓冲区，存储被截断的数据
	tmpBuffer := make([]byte, 0)
	//接收解包
	readerChannel := make(chan []byte, 16)
	go reader(readerChannel)
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}
		tmpBuffer = protocol.Depack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}
	//defer conn.Close()
}

func reader(readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel:
			Log(string(data))
		}
	}
}

func Log(v ...interface{}) {
	fmt.Println(v...)
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
/**
接受外部指令
 */
func AcceptExternalStruct (conn net.Conn)  {
	go func() {
		running := true
		reader := bufio.NewReader(os.Stdin)
		for running {
			data, _, _ := reader.ReadLine()
			command := string(data)
			if command == "status"{
				conn.Write(protocol.Enpack([]byte(command)))
			}
			if command == "open"{
				conn.Write(protocol.Enpack([]byte(command)))
			}
			if command == "close"{
				conn.Write(protocol.Enpack([]byte(command)))
			}
			if command == "stop" {
				running = false
			}
			fmt.Println("command", command)
		}
	}()
}
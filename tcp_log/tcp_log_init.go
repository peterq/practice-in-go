package tcp_log

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func Init() {
	port := "8000"
	fmt.Println("server has been start===>, listen on" + port)
	tcpAddr, _ := net.ResolveTCPAddr("tcp", ":"+port)
	//服务器端一般不定位具体的客户端套接字
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
	for {

		tcpConn, _ := tcpListener.AcceptTCP()
		go onConnect(tcpConn)

	}

}

func onConnect(con *net.TCPConn) {
	defer con.Close()
	log.Println("连接的客户端信息：", con.RemoteAddr().String())
	io.Copy(os.Stdout, con)
	log.Println("断开: ", con.RemoteAddr().String())
}

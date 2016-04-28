package main

import (
	"bufio"
	"fmt"
	"encoding/gob"
	"github.com/Kuldeep/cs733/assignment4/Assignment1Combined/fs"
	"net"
	"os"
	"strconv"
	"encoding/json"
	"strings"
)

var myid int
var pt int


var crlf = []byte{'\r', '\n'}

func check1(obj interface{}) {
	if obj != nil {
		fmt.Println(obj)
		os.Exit(1)
	}
}

type message struct{
	Cid int64
	Mssg fs.Msg
}


func reply(conn *net.TCPConn, msg *fs.Msg) bool {
	var err error
	write := func(data []byte) {
		if err != nil {
			return
		}
		_, err = conn.Write(data)
	}
	var resp string
	switch msg.Kind {
	case 'C': // read response
		resp = fmt.Sprintf("CONTENTS %d %d %d", msg.Version, msg.Numbytes, msg.Exptime)
	case 'O':
		resp = "OK "
		if msg.Version > 0 {
			resp += strconv.Itoa(msg.Version)
		}
	case 'F':
		resp = "ERR_FILE_NOT_FOUND"
	case 'V':
		resp = "ERR_VERSION " + strconv.Itoa(msg.Version)
	case 'M':
		resp = "ERR_CMD_ERR"
	case 'I':
		resp = "ERR_INTERNAL"
	case 'N':
		resp = "ERR_REDIRECT"+strconv.Itoa(msg.Version)
	default:
		fmt.Printf("Unknown response kind '%c'", msg.Kind)
		return false
	}
	resp += "\r\n"
	write([]byte(resp))
	if msg.Kind == 'C' {
		write(msg.Contents)
		write(crlf)
	}
	return err == nil
}

func serve(conn *net.TCPConn, c int64, clientch chan *fs.Msg, rn Node ) {
	//fmt.Println("Hellllooo In SErve")
	var s message
	reader := bufio.NewReader(conn)
	for {
		msg1, msgerr, fatalerr := fs.GetMsg(reader)
		//fmt.Println("client msg",msg1)
		if fatalerr != nil || msgerr != nil {
			reply(conn, &fs.Msg{Kind: 'M'})
			conn.Close()
			break
		}

		if msgerr != nil {
			if (!reply(conn, &fs.Msg{Kind: 'M'})) {
				conn.Close()
				break
			}
		}
		
		s=message{Cid:c,Mssg:*msg1}
		//fmt.Println(s)
		
		k,err := json.Marshal(s)
		if err!=nil {
			fmt.Println(err)
		}
		//fmt.Println("after marshal",string(k))
		rn.(*RaftNode).Append(k)
		resp := <-clientch
		
		if !reply(conn, resp) {
			conn.Close()
			break
		}
	}
}


func runfsMap(node Node, fs1 *fs.FS){
	for{
		
			ci := <-node.CommitChannel()
			//fmt.Println("Msg Rcvd : ", ci)
			var respMsg message
			ErrSplit:=strings.Split(ci.Err," ")
			if ErrSplit[1]=="ERR_REDIRECT"{
				Leader,_:=strconv.Atoi(ErrSplit[1])
				response:=&fs.Msg{
					Kind: 'N',
					Version: Leader,
				}
				if err :=json.Unmarshal(ci.Data, &respMsg); err != nil {
	        			panic(err)
	    			}
	    	//fmt.Println("After Unmarshal",respMsg)
	    	fmt.Println("Msg Rcvd : ", ci)
	    	//response=fs1.ProcessMsg(&(respMsg.Mssg))
	    	
	    	Clientchan := fs1.FSmap[respMsg.Cid]
	    	Clientchan <- response
			
			}else{

			 if err :=json.Unmarshal(ci.Data, &respMsg); err != nil {
	        			panic(err)
	    			}
	    	//fmt.Println("After Unmarshal",respMsg)
	    	fmt.Println("Msg Rcvd : ", ci)
	    	response1:=fs1.ProcessMsg(&(respMsg.Mssg))
	    	
	    	Clientchan := fs1.FSmap[respMsg.Cid]
	    	Clientchan <- response1

	    	}
	}

}


func serverMain(rn Node, fsport int) {
	var cid int64
	
	cid=0
	tcpaddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d",fsport))
	check1(err)
	
	tcp_acceptor, err := net.ListenTCP("tcp", tcpaddr)

	check1(err)
	var fs1 = &fs.FS{Dir: make(map[string]*fs.FileInfo, 1000)}
	fs1.FSmap=make(map[int64](chan *fs.Msg),100)
	go runfsMap(rn, fs1)
	

	for {
		
		tcp_conn, err := tcp_acceptor.AcceptTCP()
		
		check1(err)
		cid=cid+1
		
		clientCh := make(chan *fs.Msg,1)
		fs1.FSmap[cid]=clientCh
		//fmt.Println(fs1.FSmap[cid])
		go serve(tcp_conn, cid, clientCh, rn )

	}
}



func main() {
	gob.Register(VoteRequestEv{})
	argsConfig:=os.Args[1:]
	//myid = myid+1
	myid,_ = strconv.Atoi(argsConfig[0])
	pt,_ = strconv.Atoi(argsConfig[1])
	var rn Node
	rn, _ = New(myid, "config1.json")
	serverMain(rn, pt)

}

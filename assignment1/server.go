//server.go

package main

import (  
    "fmt"
    "net"
    "os"
    "strings"
    "io/ioutil"
    "strconv"
    "sync"
    "time"
)

const (  
    CONN_HOST = "127.0.0.1"
    CONN_PORT = "8080"
    CONN_TYPE = "tcp"
)

type fileinfo [500]struct{
filename string
version int64
numbytes int
exptime int64
createtime int64
}
var mutex = &sync.Mutex{}
var c fileinfo
var clength int
var d *fileinfo

func serverMain() {  
    // Listen for incoming connections.
    c[0].filename="file1.txt"
    c[0].version=0
    c[0].numbytes=113
    c[0].exptime=0
    c[1].filename="file2.txt"
    c[1].version=0
    c[1].numbytes=89
    c[1].exptime=0
    clength=2
    d = &c
    l, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
   
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
        }
        
        // Handle connections in a new goroutine.
        go commandExecutor(conn)
        go updateFiles(d)     
    }
}

// Handles incoming requests.
func commandExecutor(conn net.Conn) {  
     // Make a buffer to hold incoming data.
  for{
  buf := make([]byte, 2048)
  // Read the incoming connection into the buffer.
  Len, err := conn.Read(buf)
  if err != nil {
    fmt.Println("Error reading:", err.Error())
  }
  s := string(buf[:Len])
  tok := strings.SplitN(s,"\r\n",2)
  stringtok := strings.Split(tok[0]," ")
  tok[1]=strings.Trim(tok[1],"\r\n")
  toklen := len(stringtok)
  
  
if stringtok[0]=="write"{
  if toklen==3||toklen==4{
  mutex.Lock()
  if _, err := os.Stat(stringtok[1]); err == nil {
  for i := 0 ; i<=clength;i++{
  if c[i].filename== strings.TrimSpace(stringtok[1]){
  c[i].numbytes, err=strconv.Atoi(stringtok[2])
  if toklen==4{
  c[i].exptime, err=strconv.ParseInt(stringtok[3], 10, 64)
  c[i].createtime = time.Now().Unix()
  }
  
  }
   }
  }else {
  c[clength].filename=strings.TrimSpace(stringtok[1])
  c[clength].numbytes, err=strconv.Atoi(strings.TrimSpace(stringtok[2]))
  c[clength].version=int64(0)
  if toklen==4{  
  c[clength].exptime, err=strconv.ParseInt(strings.TrimSpace(stringtok[3]), 10, 64)
  }else{
  c[clength].exptime = int64(0)
  }
  c[clength].createtime = time.Now().Unix()
  clength=clength+1
  }
  err := ioutil.WriteFile(stringtok[1],[]byte(tok[1]), 0777)
  
  if err!=nil{
  conn.Write([]byte("ERR_INTERNAL\r\n"))
  }
  conn.Write([]byte("OK"+" "+strconv.FormatInt(c[clength-1].version,10)+"\r\n"))
  mutex.Unlock()
  } else {
  conn.Write([]byte("ERR_CMD_ERR\r\n"))
  }
  
  
}else  if strings.TrimSpace(stringtok[0])=="read"{
    
   if toklen==2{
        t := strings.TrimSpace(stringtok[1])
        mutex.Lock()
  	dat, err := ioutil.ReadFile(t)
 
     if err!=nil{
        conn.Write([]byte("ERR_FILE_NOT_FOUND\r\n"))
       }
       
  for i := 0; i <clength; i++{
  if c[i].filename == strings.TrimSpace(stringtok[1]){
  s := strconv.FormatInt(c[i].version,10)
  p := strconv.Itoa(c[i].numbytes)
  conn.Write([]byte("CONTENTS"+" "))
  conn.Write([]byte(s+" "+p+"\r\n"))
	}
  }
  
  conn.Write([]byte(dat))
  conn.Write([]byte("\r\n"))
  mutex.Unlock()
  }else{
  conn.Write([]byte("ERR_CMD_ERR\r\n"))
  }
  
  
  }else if stringtok[0]=="cas"{
 if toklen==4 || toklen ==5{
  mutex.Lock()
if _, err := os.Stat(stringtok[1]); err == nil {
  for i := 0 ; i<clength;i++{
  if c[i].filename== strings.TrimSpace(stringtok[1]){
  if strconv.FormatInt(c[i].version,10) == stringtok[2]{
  c[i].numbytes, err=strconv.Atoi(stringtok[3])
  if toklen == 5{
  c[i].exptime, err=strconv.ParseInt(stringtok[4], 10, 64)
  c[i].createtime=time.Now().Unix()
  }
  c[i].version++
  err := ioutil.WriteFile(stringtok[1], []byte(tok[1]), 0777)
   if err!=nil{
  conn.Write([]byte("ERR_INTERNAL\r\n"))
  }
  }else{
  conn.Write([]byte("ERR_VERSION"+" "))
  conn.Write([]byte(strconv.FormatInt(c[i].version,10)))
  conn.Write([]byte("\r\n"))
  }
  }
  }
} else {
  conn.Write([]byte("ERR_FILE_NOT_FOUND\r\n"))
  }
 conn.Write([]byte("OK"+" "))
  for i:=0;i<clength;i++{
  if c[i].filename== strings.TrimSpace(stringtok[1]){
  conn.Write([]byte(strconv.FormatInt(c[i].version,10)+"\r\n"))
   }
   }
   mutex.Unlock()
  }else{
  conn.Write([]byte("ERR_CMD_ERR\r\n"))
  //conn.Close()
  }
  
  
}else if stringtok[0]=="delete"{
if toklen==2{
mutex.Lock()
err := os.Remove(strings.TrimSpace(stringtok[1]))

      if err != nil {
          conn.Write([]byte("ERR_FILE_NOT_FOUND\r\n"))
      }
      for i:=0;i<clength;i++{
  if c[i].filename== strings.TrimSpace(stringtok[1]){
   c[i].filename="$"
  }
  }
  conn.Write([]byte("OK\r\n"))
mutex.Unlock()
}else{
conn.Write([]byte("ERR_CMD_ERR\r\n"))
   }
  }
 }
}



func main(){
serverMain()
}

func updateFiles(c *fileinfo) {

t := time.Now().Unix()

for i:=0; i< clength; i++{
if c[i].exptime!=0{
x:= c[i].createtime+c[i].exptime
if x < t{
mutex.Lock()
 err := os.Remove(c[i].filename)
 if err != nil {
          
      }
mutex.Unlock()      
 }
 
 }
 }
 }


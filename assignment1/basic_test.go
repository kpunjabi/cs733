package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"
	"sync"
)


// Simple serial check of getting and setting
func TestTCPSimple(t *testing.T) {
	go serverMain()
	var wg sync.WaitGroup
	wg.Add(5)
   
	for j:=0;j<5;j++{
	go func(k int){
	name := "abc.txt"
	contents1 := "by vhajhj klkle this is new version"
	contents := "Hiii this is CS733 course here u go with assignment-1"
	exptime := 300000
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Error(err.Error()) // report error through testing framework
	}
        //fmt.Println("Hello")
	scanner := bufio.NewScanner(conn)


	// Write a file
	fmt.Fprintf(conn, "write %v %v %v\r\n%v\r\n", name, len(contents), exptime, contents+":    Client no-"+strconv.Itoa(k))
	scanner.Scan() // read first line
	resp := scanner.Text() // extract the text from the buffer
	arr := strings.Split(resp," ") // split into OK and <version>
	expect(t, arr[0], "OK")
	//fmt.Println(resp)
	//fmt.Println(arr[1])
	ver, err := strconv.Atoi(arr[1]) // parse version as number
	if err != nil {
		t.Error("Non-numeric version found")
	}
	version := int64(ver)


	fmt.Fprintf(conn, "read %v\r\n", name) // try a read now
	scanner.Scan()
	arr = strings.Split(scanner.Text(), " ")
	expect(t, arr[0], "CONTENTS")
	expect(t, arr[1], fmt.Sprintf("%v", version)) // expect only accepts strings, convert int version to string
	expect(t, arr[2], fmt.Sprintf("%v", len(contents)))	
	scanner.Scan()
	//fmt.Println("Hello")
	expect(t, scanner.Text(), contents)
	//fmt.Println("Before cas")
	//ver = 2
	
	fmt.Fprintf(conn, "cas %v %v %v %v\r\n%v\r\n", name, ver, len(contents), exptime, contents1+":    Client no-"+strconv.Itoa(k))
	scanner.Scan()
	resp = scanner.Text()
	arr = strings.Split(resp," ") // split into OK and <version>
	expect(t, arr[0], "OK")
//	vers, err := strconv.Atoi(arr[1]) // parse version as number
	if err != nil {
		t.Error("Non-numeric version found")
	}
//	version1 := int64(vers)
//	fmt.Println(version1)
	/*
	fmt.Fprintf(conn, "delete %v\r\n", name) // try a read now
	scanner.Scan()
	arr = strings.Split(scanner.Text(), " ")
	expect(t, arr[0], "OK") 
	
	*/
	wg.Done()
	}(j)
	
	}
	wg.Wait()
    fmt.Println("Done")
}

// Useful testing function
func expect(t *testing.T, a string, b string) {
	if a != b {
		t.Error(fmt.Sprintf("Expected %v, found %v", b, a)) // t.Error is visible when running `go test -verbose`
	}
}

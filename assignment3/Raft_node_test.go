package main

import (
	"fmt"
	//cl "github.com/cs733-iitb/cluster"
	//lgd "github.com/cs733-iitb/log"
	//"os"
	"io/ioutil"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var rafts [5]Node
var ldr Node
var max int
//MakeRafts function creates all the nodes and initializes them

func makeRafts() [5]Node {
	var k int
    
    var val int
    val=0
    for i:=1;i<=5;i++{
     k=rand.Intn(10)
     if k>val{
     	val=k
     	max=i
     }
     data := strconv.Itoa(k)+" "+strconv.Itoa(0)
     dt:=[]byte(data)
    switch(i) {
  
	case 1:
		err := ioutil.WriteFile("persdata1",dt, 0644)
		check(err)
	case 2:
		err := ioutil.WriteFile("persdata2",dt, 0644)
		check(err)
	case 3:
		err := ioutil.WriteFile("persdata3",dt, 0644)
		check(err)
	case 4:
		err := ioutil.WriteFile("persdata4",dt, 0644)
		check(err)
	case 5:
		err := ioutil.WriteFile("persdata5",dt, 0644)
		check(err)

	}
}
	var raftnodes [5]Node
	var myid = 1

	raftnodes[0], _ = New(myid, "config1.json")
	myid++

	raftnodes[1], _ = New(myid, "config1.json")

	myid++
	raftnodes[2], _ = New(myid, "config1.json")
	myid++
	raftnodes[3], _ = New(myid, "config1.json")
	myid++
	raftnodes[4], _ = New(myid, "config1.json")

	
	return raftnodes

}

//GetLeader function returns leader from among the five nodes: Here I have configured initially such that node id:3 will always be the one to become leader

func getLeader(nd [5]Node) Node {
	var x, l int
	var ld int
	var arr [5]int
	arr = [5]int{0, 0, 0, 0, 0}
	for i := 0; i < 5; i++ {
		x = nd[i].LeaderId()
	}
	for i := 0; i < 5; i++ {
		x = nd[i].LeaderId()
		if x-1 >= 0 {
			(arr[x-1])++
		}
	}
	max := arr[0]
	for l = 0; l < 5; l++ {
		if arr[l] > max {
			max = arr[l]
			break
		}
	}

	ld = l + 1
	if max >= 3 {
		for i := 0; i < 5; i++ {
			if nd[i].(*RaftNode).config.Id == ld {
				return nd[i]
			}
		}
	}
	return nil
}

// This function tests whether our leader get successfully created and tests it

func TestLeaderCreation(t *testing.T) {
	fmt.Println("Testing Leader Creation")
	rafts = makeRafts()
	// array of []raft.Node
	ret := rafts[max-1].(*RaftNode).sm.ProcessEvent(TimeoutEv{rand.Intn(10)})
	rafts[max-1].(*RaftNode).doActions(ret)
	time.Sleep(10 * time.Second)
	//time.Sleep(10 * time.Second)
	p6 := [4]int{1, 2, 3, 5}
	//	var out1 = make([]interface{}, 1)
	//	var out2 = make([]interface{}, 1)

	ldr = getLeader(rafts)
    

	expectArr(t, ldr.(*RaftNode).sm.peers, p6, "Test:0")
	expect(t, ldr.(*RaftNode).sm.status, "leader", "Test:1")

	for i := 0; i < 5; i++ {

		expect(t, strconv.Itoa(rafts[i].(*RaftNode).sm.currentTerm), strconv.Itoa(ldr.(*RaftNode).sm.currentTerm), "Test:2")

	}
    time.Sleep(20*time.Second)
	for i := 0; i < 5; i++ {
		//time.Sleep(2 * time.Second)
		expect(t, strconv.Itoa(rafts[i].(*RaftNode).sm.currentLeader), strconv.Itoa(ldr.(*RaftNode).sm.id), "Test:3")

	}

}


// This function tests for the append() func of the leader and checks if it data being committed over all the peer node's channel

func TestAppend(t *testing.T) {
	fmt.Println("Testing Leader Append")

	ldr.(*RaftNode).Append([]byte("foo"))

	time.Sleep(20 * time.Second)

	//fmt.Println("checking")
	for _, node := range rafts {
		select {
		// to avoid blocking on channel.
		case ci := <-node.CommitChannel():
		//	fmt.Println(string(ci.Data))
			if ci.Err != nil {
				t.Fatal(ci.Err)
			}
			if string(ci.Data) != "foo" {
				t.Fatal("Got different data")
			}
		default:
			t.Fatal("Expected message on all nodes")
		}
	}

}


func TestMultipleAppend(t *testing.T) {
fmt.Println("Testing Multiple Append Request")
var lastlog []byte
var dat []byte
ldr.(*RaftNode).Append([]byte("print1"))
time.Sleep(10 * time.Second)
_,dat=ldr.Get(int(ldr.(*RaftNode).logRaft.GetLastIndex()))
fmt.Println(string(dat))
ldr.(*RaftNode).Append([]byte("print2"))
time.Sleep(10 * time.Second)
_,dat=ldr.Get(int(ldr.(*RaftNode).logRaft.GetLastIndex()))
fmt.Println(string(dat))
ldr.(*RaftNode).Append([]byte("print3"))
time.Sleep(10 * time.Second)
_,dat=ldr.Get(int(ldr.(*RaftNode).logRaft.GetLastIndex()))
fmt.Println(string(dat))


_, lastlog = rafts[1].Get(int(rafts[1].(*RaftNode).logRaft.GetLastIndex()))

if string(lastlog) != "print3" {
t.Fatal("log mismatch"); t.Fatal(rafts[1].(*RaftNode).logRaft.GetLastIndex())
}

//checking everyone's leaderid
for i := 0; i < len(rafts); i++ {
if ldr.(*RaftNode).Id() != int(rafts[i].(*RaftNode).LeaderId()) {
t.Fatal("unpropagated value for leader id")
}
}
//cleanup()
}



//This function checks for the get data function to get the data at a particular Index in the log

func TestGetData(t *testing.T) {
	fmt.Println("Testing Get Data")
	var m int
	err, dt := ldr.Get(2)
	if err != nil {
		fmt.Println(err)
	}

	expect(t, string(dt), "foo", "Test:4")
	time.Sleep(5 * time.Second)

	//Checking that all the nodes have data atleast upto leader's commit Index
	for i := 0; i < 4; i++ {
		m = ldr.(*RaftNode).sm.peers[i]
		err, dt1 := rafts[m-1].Get(ldr.(*RaftNode).sm.commitIndex)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println(string(dt1))
		expect(t, string(dt1), "foo", "Test:4")
	}
}

//Testing CommitIndex() function of RaftNode
func TestCommitIndex(t *testing.T) {
	fmt.Println("Testing Commit Index")
	//var l int
	//var cindex int
	n := ldr.(*RaftNode).sm.commitIndex
	k := ldr.CommittedIndex()
	expect(t, strconv.Itoa(k), strconv.Itoa(n), "Test:5")
	/*
		//Now Let's check if all the follower's have their commit Index to that of Leader's
		for i:=0;i<4;i++{
			l=ldr.(*RaftNode).sm.peers[i]
			cindex=rafts[l-1].(*RaftNode).CommittedIndex()
			expect(t, strconv.Itoa(cindex), strconv.Itoa(k), "Test:5")
		}
	*/

}

//Testing for the Leader Id of all the nodes :Taking no partition case
func TestLeaderId(t *testing.T) {
	fmt.Println("Testing Leader Id")
	k := ldr.(*RaftNode).LeaderId()
	//time.Sleep(2 * time.Second)
	//expect(t, strconv.Itoa(k), strconv.Itoa(n), "Test:5")
	//Now Let's check if all the follower's have their commit Index to that of Leader's
	for i := 0; i < 4; i++ {
		l := ldr.(*RaftNode).sm.peers[i]
		lid := rafts[l-1].(*RaftNode).LeaderId()
		expect(t, strconv.Itoa(k), strconv.Itoa(lid), "Test:5")
	}

}

// Testing for the Shut Down of the RaftNode

func TestShutDown(t *testing.T) {
	fmt.Println("Testing Shut Down")

	ldr.Shutdown()
	//fmt.Println("Before ShutDown calling Append")
	ldr.(*RaftNode).Append([]byte("hello"))
	//fmt.Println("Hello world")

}

//Generates error for two String Comparisions

func expect(t *testing.T, a string, b string, c string) {
	if a != b {
		t.Error(fmt.Sprintf("Expected %v, found %v  Error type:  %v", b, a, c)) // t.Error is visible when running `go test -verbose`
	}
}

//Generates error for two array Comparisions

func expectArr(t *testing.T, a [4]int, b [4]int, c string) {
	if a != b {
		t.Error(fmt.Sprintf("Expected %v, found %v  Error type:  %v", b, a, c)) // t.Error is visible when running `go test -verbose`
	}
}

package main

import (
	"github.com/cs733-iitb/cluster"
    "io/ioutil"
  	"os"
	"strings"
	"strconv"
)

// Create a raft object with the
func New(config Config) Node {

f, err := os.Create(config.LogDir)
    check(err)
defer f.Close()
var dat []byte

switch(config.Id){
	
case 1: 
        dat, err = ioutil.ReadFile("Persistentdata1")
        check(err)
case 2: 
        dat, err = ioutil.ReadFile("Persistentdata2")
        check(err)
case 3: 
        dat, err = ioutil.ReadFile("Persistentdata3")
        check(err)
case 4: 
		dat, err = ioutil.ReadFile("Persistentdata4")
        check(err)
case 5: 
		dat, err = ioutil.ReadFile("Persistentdata5")
        check(err)
 }   

stringtok := strings.Split(string(dat)," ")
term , _ := strconv.Atoi(strings.TrimSpace(stringtok[1]))
votedFor, _ :=strconv.Atoi(strings.TrimSpace(stringtok[2]))

sm = &StateMachine{id: config.Id, status: "follower", currentTerm: term,votedFor: votedFor}

}


func makeRafts(){
var raftnodes []Node
var myid int
myid=1
server1, err := cluster.New(myid, "config1.json")

	if err != nil {
		panic(err)
	}
myid++
server2, err := cluster.New(myid, "config1.json")

	if err != nil {
		panic(err)
	}
myid++
server3, err := cluster.New(myid, "config1.json")

	if err != nil {
		panic(err)
	}
myid++
server4, err := cluster.New(myid, "config1.json")

	if err != nil {
		panic(err)
	}
myid++
server5, err := cluster.New(myid, "config1.json")

	if err != nil {
		panic(err)
	}


raftnodes[0]= New(Config{cluster: server1, Id: 1,LogDir:"log1",ElectionTimeout: 100, HeartbeatTimeout:50})
raftnodes[1]= New(Config{cluster: server2, Id: 2,LogDir:"log2",ElectionTimeout: 100, HeartbeatTimeout:50 })
raftnodes[2]= New(Config{cluster: server3, Id: 3,LogDir:"log3",ElectionTimeout: 100, HeartbeatTimeout:50})
raftnodes[3]= New(Config{cluster: server4, Id: 4,LogDir:"log4",ElectionTimeout: 100, HeartbeatTimeout:50})
raftnodes[4]= New(Config{cluster: server5, Id: 5,LogDir:"log5",ElectionTimeout: 100, HeartbeatTimeout:50 })

return raftnodes

}
func TestBasic (t *testing.T) {
	
     
		rafts := makeRafts()
		// array of []raft.Node
	     ldr := getLeader(rafts)
		ldr.Append("foo")
		time.Sleep(1 time.Second)
		for _, node:= rafts {
			select{
					// to avoid blocking on channel.
					case ci := <- node.CommitChannel():
						if ci.err != nil{t.Fatal(ci.err)}
						if string(ci.data) !="foo"{
						t.Fatal(
						"Got different data"
						)
						}
				    default : t.Fatal("Expected message on all nodes")
			}
		}

}

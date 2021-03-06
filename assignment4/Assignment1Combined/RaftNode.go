package main

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	//"reflect"
	"math/rand"
	"strconv"
	"strings"
	"time"

	cluster "github.com/cs733-iitb/cluster"
	logf "github.com/cs733-iitb/log"
	"os"
)

type Time interface{}
//Message Struct to send, receive message into the Envelope
type Msg1 struct {
	From  int
	Event interface{}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Event interface{}

//Node Interface

type Node interface {
	Append([]byte)
	CommitChannel() <-chan CommitInfo
	CommittedIndex() int
	Get(index int) (error, interface{})
	Id() int
	LeaderId() int
	Shutdown()
}

type CommitInfo struct {
	Data  []byte
	Index int   // or int .. whatever you have in your code
	Err   string // Err can be errred
}

type Config struct {
	Cluster          []NetConfig
	Id               int
	LogDir           string
	ElectionTimeout  int
	HeartbeatTimeout int
}

func ToConfigure(configuration interface{}) (config *Config, err error) {
	var cfg1 Config
	var ok bool
	var configFile string
	if configFile, ok = configuration.(string); ok {
		var f *os.File
		if f, err = os.Open(configFile); err != nil {
			return nil, err
		}
		defer f.Close()
		dec := json.NewDecoder(f)
		if err = dec.Decode(&cfg1); err != nil {
			return nil, err
		}
	} else if cfg1, ok = configuration.(Config); !ok {
		return nil, errors.New("Expected a configuration.json file or a Config structure")
	}
	return &cfg1, nil
}

type NetConfig struct {
	Id   int
	Host string
	Port int
}

type RaftNode struct { // implements Node interface
	config    *Config
	sm        *StateMachine
	servl     cluster.Server
	logRaft   logf.Log
	Timer     *time.Timer
	eventCh   chan Event
	timeoutCh chan Time
	commitCh  chan CommitInfo
	HbeatTimer *time.Ticker
	Terminate chan int
}

func (rn *RaftNode) doActions(actions []interface{}) {

	gob.Register(AppendEntriesReqEv{})

	gob.Register(VoteRespEv{})


	for i := range actions {
		var msgid int64
		var m Msg1
		switch actions[i].(type) {

		case send:
			gob.Register(VoteRequestEv{})
			gob.Register(AppendEntriesRespEv{})
			m.From = actions[i].(send).peerId
			m.Event = actions[i].(send).event
			rn.servl.Outbox() <- &cluster.Envelope{Pid: m.From, MsgId: msgid, Msg: m}
			msgid = msgid + 1

		case Alarm:
			k:=actions[i].(Alarm)
			rn.Timer.Stop()
			rand.Seed(time.Now().UnixNano())
			RandNo := rand.Intn(500)
			rn.Timer = time.AfterFunc(time.Duration(k.time+RandNo)*time.Millisecond, func() { rn.timeoutCh <- TimeoutEv{} })
		
		case Commit:
			gob.Register(AppendEntriesReqEv{})
			rn.commitCh <- CommitInfo{Data: actions[i].(Commit).data, Err: actions[i].(Commit).err, Index: actions[i].(Commit).index}
			

		case LogStore:
			err := rn.logRaft.Append(log1{LogTerm:rn.sm.currentTerm,LogIndex:actions[i].(LogStore).Index,Command:actions[i].(LogStore).Data})
			if err != nil {
				fmt.Println(err)
			}


		}
	}
}

func (rn *RaftNode) Append(data []byte) {
	go func() { rn.eventCh <- AppendEv{data1: data} }()

}

func (rn *RaftNode) CommittedIndex() int {

	return rn.sm.commitIndex

}

func (rn *RaftNode) CommitChannel() <-chan CommitInfo {
	return rn.commitCh

}

func (rn *RaftNode) Get(index int) (error, interface{}) {
	var data interface{}
	data, _ = rn.logRaft.Get(int64(index))
	return nil, data
}

func (rn *RaftNode) Id() int {
	return rn.config.Id
}

func (rn *RaftNode) LeaderId() int {
	return rn.sm.currentLeader
}

func (rn *RaftNode) Shutdown() {
	rn.servl.Close()
	rn.Terminate <- 1

}

func Lineread() string {
	bio := bufio.NewReader(os.Stdin)
	line, err := bio.ReadString('\n')
	if err == nil && line != "" {
		line = strings.TrimSpace(line)
	}
	return line
}

// Create a raft object with the
func New(myid int, configuration interface{}) (nd Node, err error) {

	//var raft RaftNode
	raft := new(RaftNode)

	config := new(Config)

	if config, err = ToConfigure(configuration); err != nil {
		return nil, err
	}

	config.Id = myid
	
	config.LogDir=fmt.Sprintf("log_%d", myid)

	raft.config = config
  
	raft.servl, err = cluster.New(myid, "cluster_test_config.json")

	if err != nil {
		panic(err)
	}


	lg, err := logf.Open(config.LogDir)
	if err != nil {
		//t.Fatal(err)
		fmt.Println(err)
	}
	
	lg.RegisterSampleEntry(log1{})
	var dat []byte

	switch config.Id {

	case 1:
		dat, err = ioutil.ReadFile("persdata1")
		check(err)
	case 2:
		dat, err = ioutil.ReadFile("persdata2")
		check(err)
	case 3:
		dat, err = ioutil.ReadFile("persdata3")
		check(err)
	case 4:
		dat, err = ioutil.ReadFile("persdata4")
		check(err)
	case 5:
		dat, err = ioutil.ReadFile("persdata5")
		check(err)

	}

	stringtok := strings.Split(strings.TrimSpace(string(dat)), " ")
	term, _ := strconv.Atoi(strings.TrimSpace(stringtok[0]))
	votedFor, _ := strconv.Atoi(strings.TrimSpace(stringtok[1]))
	var prs [4]int
	var k int
	k = 0

	for i := 0; i < 5; i++ {
		if i != config.Id-1 {
			a := config.Cluster[i].Id
			prs[k] = a
			k++
		}
	}
	//fmt.Println(raft.config.ElectionTimeout,raft.config.HeartbeatTimeout)
	raft.sm = &StateMachine{id: config.Id, status: "follower", peers: prs, currentTerm: term, votenegCount: 0, votedFor: votedFor, prevLogIndex: -1, lastLogIndex: 0,nextIndex:[6]int{1,1,1,1,1,1}}
	raft.sm.ElectionTimeout=raft.config.ElectionTimeout
	raft.sm.HeartBeatTimeout=raft.config.HeartbeatTimeout
	go raft.processEvents()
	//fmt.Println(raft.sm)
	go raft.sm.ProcessEvent(make([]interface{}, 1))
	
	raft.eventCh = make(chan Event, 10)
	raft.timeoutCh = make(chan Time, 2)
	raft.commitCh = make(chan CommitInfo, 1)
	raft.Terminate = make(chan int, 1)
	raft.HbeatTimer = time.NewTicker(time.Duration(config.HeartbeatTimeout) * time.Millisecond)
	rand.Seed(time.Now().UnixNano())
	RandNo := rand.Intn(500)
	
	raft.Timer = time.AfterFunc(time.Duration(config.ElectionTimeout+RandNo)*time.Millisecond, func() { raft.timeoutCh <- TimeoutEv{} })
    
	raft.logRaft = *lg
	raft.sm.log=&raft.logRaft
	//raft.logRaft.TruncateToEnd(0)
	er:=raft.logRaft.Append(log1{LogIndex:0,LogTerm:0,Command:[]byte("init")})
    if(er!=nil){
    	fmt.Println(er)
    	}
    /*
		for i := 0; i < int(raft.logRaft.GetLastIndex()); i++ {
			raft.sm.log[i].logIndex = i
			_, raft.sm.log[i].command = raft.Get(i)
		}
    */
	nd = raft
	return nd, nil

}

func (rn *RaftNode) processEvents() {
	var shut int
	gob.Register(Msg1{})
	gob.Register(Alarm{})
	gob.Register(VoteRequestEv{})
	
	go func() {
		for {
			env := <-rn.servl.Inbox()
			//fmt.Println("Inbox",env)
			ret := rn.sm.ProcessEvent(env.Msg.(Msg1).Event)
			//fmt.Println("abc",ret)
			rn.doActions(ret)
			
		}
	}()

	for {
		var ev Event
		//fmt.Println("in raft node's process events")
	select {
		case ev = <-rn.eventCh:
			
		case ev = <-rn.timeoutCh:
			//rn.sm.ProcessEvent(ev)
		
		case <-rn.HbeatTimer.C:
			if rn.sm.status == "leader" {
			rn.eventCh <- TimeoutEv{}
			}
		
		case shut = <-rn.Terminate:
			if shut == 1 {
				return
			}
			fmt.Println("Checking ShutDown")
		}
		//fmt.Println(ev)
		actions := rn.sm.ProcessEvent(ev)
		//fmt.Println(actions)
		rn.doActions(actions)

	}

}

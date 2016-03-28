package main

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	//"reflect"
	"strconv"
	"strings"
	"time"

	cluster "github.com/cs733-iitb/cluster"
	logf "github.com/cs733-iitb/log"
	"os"
)

//Message Struct to send, receive message into the Envelope
type Msg struct {
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
	Get(index int) (error, []byte)
	Id() int
	LeaderId() int
	Shutdown()
}

type CommitInfo struct {
	Data  []byte
	Index int   // or int .. whatever you have in your code
	Err   error // Err can be errred
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
	eventCh   chan Event
	timeoutCh chan int
	commitCh  chan CommitInfo
	Terminate chan int
}

func (rn *RaftNode) doActions(actions []interface{}) {

	gob.Register(AppendEntriesReqEv{})

	gob.Register(VoteRespEv{})

	for i := range actions {
		var msgid int64
		var m Msg
		switch actions[i].(type) {
		case Alarm:
			rn.timeoutCh <- actions[i].(Alarm).time

		case send:
			gob.Register(VoteRequestEv{})
			gob.Register(AppendEntriesRespEv{})
			m.From = actions[i].(send).peerId
			m.Event = actions[i].(send).event
			//fmt.Println(actions[i], reflect.TypeOf(actions[i]))
			rn.servl.Outbox() <- &cluster.Envelope{Pid: m.From, MsgId: msgid, Msg: m}
			msgid = msgid + 1

		case Commit:
			gob.Register(AppendEntriesReqEv{})
			//err1 := errors.New(actions[i].(Commit).err)

			rn.commitCh <- CommitInfo{Data: actions[i].(Commit).data, Err: nil, Index: actions[i].(Commit).index}
			if rn.sm.status == "leader" {

				rn.sm.commitIndex++
				for i := 0; i < 4; i++ {
					m.From = rn.sm.peers[i]
					m.Event = AppendEntriesReqEv{Term: rn.sm.currentTerm, LeaderId: rn.sm.id, LeaderCommit: rn.sm.commitIndex, Entries: []byte{}, PrevLogIndex: rn.sm.prevLogIndex, PrevLogTerm: rn.sm.prevLogTerm}
					rn.servl.Outbox() <- &cluster.Envelope{Pid: m.From, MsgId: msgid, Msg: m}

				}

				time.Sleep(2 * time.Second)
			}

		case LogStore:
			err := rn.logRaft.Append([]byte(actions[i].(LogStore).Data))
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

func (rn *RaftNode) Get(index int) (error, []byte) {
	var data []byte
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
	switch config.Id {

	case 1:
		config.LogDir = "log1"
	case 2:
		config.LogDir = "log2"
	case 3:
		config.LogDir = "log3"
	case 4:
		config.LogDir = "log4"
	case 5:
		config.LogDir = "log5"

	}

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

	raft.sm = &StateMachine{id: config.Id, status: "follower", peers: prs, currentTerm: term, votenegCount: 0, votedFor: votedFor, prevLogIndex: -1, lastLogIndex: 0}
	go raft.processEvents()
	go raft.sm.ProcessEvent(make([]interface{}, 1))

	raft.eventCh = make(chan Event, 10)
	raft.timeoutCh = make(chan int, 2)
	raft.commitCh = make(chan CommitInfo, 1)
	raft.Terminate = make(chan int, 1)

	raft.logRaft = *lg
   
	for i := 0; i < int(raft.logRaft.GetLastIndex()); i++ {
		raft.sm.log[i].logIndex = i
		_, raft.sm.log[i].command = raft.Get(i)
	}

	nd = raft

	return nd, nil

}

func (rn *RaftNode) processEvents() {
	var shut int
	gob.Register(Msg{})
	go func() {
		for {
			env := <-rn.servl.Inbox()
			//fmt.Println(env)
			ret := rn.sm.ProcessEvent(env.Msg.(Msg).Event)
			//fmt.Println(ret)
			time.Sleep(1 * time.Second)
			rn.doActions(ret)
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		var ev Event
		//fmt.Println("in raft node's process events")
		select {
		case ev = <-rn.eventCh:

		case ev = <-rn.timeoutCh:
			//rn.sm.ProcessEvent(ev)
		case shut = <-rn.Terminate:
			if shut == 1 {
				return
			}
			fmt.Println("Checking ShutDown")
		}

		actions := rn.sm.ProcessEvent(ev)

		rn.doActions(actions)

	}

}

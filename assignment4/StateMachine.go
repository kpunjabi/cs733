package main

import (
	//"fmt"
	logf "github.com/cs733-iitb/log"
	"math"
	"math/rand"
	"strconv"
	"time"
)

var sm *StateMachine
var num_resp int

type log1 struct {
	LogIndex int
	LogTerm  int
	Command  []byte
}

type send struct {
	peerId int
	event  interface{}
}

type LogStore struct {
	Index int
	Data  []byte
}

//Timeout
type TimeoutEv struct {
	//time int
}

//global variable
var count int

type Alarm struct {
	time int
}

type StateMachine struct {
	id               int // server id
	status           string
	peers            [4]int // other server ids
	votedFor         int
	currentTerm      int
	commitIndex      int
	prevLogIndex     int
	prevLogTerm      int
	lastLogIndex     int
	lastLogTerm      int
	voteCount        int
	votenegCount     int
	voteReceived     [6]int
	nextIndex        [6]int
	matchIndex       [6]int
	currentLeader    int
	appendresCount   int
	appendIdresponse [6]int
	log              *logf.Log
	HeartBeatTimeout int
	ElectionTimeout  int
}

//AppendEntries Request
type AppendEntriesReqEv struct {
	Term         int
	LeaderId     int
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []byte
	LeaderCommit int
}

//AppendEntries Response
type AppendEntriesRespEv struct {
	Id      int
	Term    int
	Success bool
}

//var sm *StateMachine

//var ret =make([]interface{},1)
func StateChange(interface{}) {

}

func (a *AppendEntriesReqEv) AppendReqHandlerF() []interface{} {

	replyfalse := AppendEntriesRespEv{sm.id, sm.currentTerm, false}
	replytrue := AppendEntriesRespEv{sm.id, sm.currentTerm, true}
	var ret = make([]interface{}, 0)

	if a.Term < sm.currentTerm {

		ret = append(ret, send{a.LeaderId, replyfalse})

	} else {
		sm.currentLeader = a.LeaderId

		if len(a.Entries) == 0 {
			ret = append(ret, Alarm{sm.ElectionTimeout})
			sm.currentTerm = a.Term
			sm.votedFor = 0
			if a.LeaderCommit > sm.commitIndex {
				sm.commitIndex = a.LeaderCommit
				k, _ := sm.log.Get(int64(a.LeaderCommit))
				ret = append(ret, Commit{index: a.LeaderCommit, data: k.(log1).Command})
				
			}

		} else {
			m, _ := sm.log.Get(int64(a.PrevLogIndex))

			sm.lastLogIndex = m.(log1).LogIndex
			sm.lastLogTerm = m.(log1).LogTerm

			if sm.lastLogIndex == a.PrevLogIndex && sm.lastLogTerm == a.PrevLogTerm {

				sm.lastLogIndex++

				sm.currentTerm = a.Term
				sm.votedFor = 0

				ret = append(ret, LogStore{sm.lastLogIndex, a.Entries})

				if a.LeaderCommit > sm.commitIndex {

					sm.commitIndex = int(math.Min(float64(a.LeaderCommit), float64(sm.lastLogIndex)))
				}
				ret = append(ret, LogStore{sm.lastLogIndex, a.Entries})
				ret = append(ret, send{a.LeaderId, replytrue})

			} else if (sm.lastLogIndex == a.PrevLogIndex) && (sm.lastLogTerm != a.PrevLogTerm) {
				//Delete the existing entry and those that follows trying to make its log equivalent to that of leader
				ret = append(ret, send{a.LeaderId, replyfalse})
			}

		}

	}

	//fmt.Println(ret)
	return ret

}

func (ae *AppendEntriesReqEv) AppendReqHandlerC() []interface{} {

	replyfalse := AppendEntriesRespEv{sm.id, sm.currentTerm, false}

	var can = make([]interface{}, 0)
	if ae.Term < sm.currentTerm {

		can = append(can, send{ae.LeaderId, replyfalse})
	} else {

		sm.status = "follower"
		sm.currentTerm = ae.Term
		sm.votedFor = 0
		can = append(can, Alarm{sm.ElectionTimeout})

		if len(ae.Entries) == 0 {

			can = append(can, send{ae.LeaderId, AppendEntriesRespEv{sm.id, sm.currentTerm, true}})

		} else {

			can = append(can, send{ae.LeaderId, AppendEntriesRespEv{sm.id, sm.currentTerm, false}})

		}
	}

	return can
}

func (a *AppendEntriesReqEv) AppendReqHandlerL() []interface{} {
	var lead = make([]interface{}, 0)
	if sm.currentTerm <= a.Term {
		lead = append(lead, send{a.LeaderId, AppendEntriesRespEv{sm.id, sm.currentTerm, false}})
		return lead
	} else {
		sm.currentTerm = a.Term
		sm.status = "follower"
		sm.votedFor = 0
		lead = append(lead, Alarm{sm.HeartBeatTimeout})
		return lead
	}
	return (make([]interface{}, 0))

}

func (ar *AppendEntriesRespEv) AppendEntriesresF() []interface{} {
	var ret6 = make([]interface{}, 0)
	if ar.Term > sm.currentTerm {
		sm.currentTerm = ar.Term
		sm.votedFor = 0
	}
	ret6 = append(ret6, "ERROR_OCCURED")
	return (ret6)
	//return(make([]interface{},1))
}

func (ar *AppendEntriesRespEv) AppendEntriesresC() []interface{} {
	var ret7 = make([]interface{}, 0)
	if ar.Term > sm.currentTerm {
		sm.currentTerm = ar.Term
		sm.status = "follower"
		sm.votedFor = 0
		ret7 = append(ret7, Alarm{sm.ElectionTimeout})
	}
	ret7 = append(ret7, "ERROR_OCCURED")
	return (ret7)

}

func (ar *AppendEntriesRespEv) AppendEntriesresL() []interface{} {
	//fmt.Println(ar)
	var aresp = make([]interface{}, 0)
	if ar.Success == true {
		sm.appendIdresponse[ar.Id] = 1
		sm.appendresCount++

		sm.matchIndex[ar.Id]++
        
		g := sm.nextIndex[ar.Id]
		

		sm.nextIndex[ar.Id]++

		count = 0

		for i := 0; i < 6; i++ {
			if g <= sm.matchIndex[i] {
				count++
			}
		}

		if num_resp == 0 {
			count = 0
		}

		if count >= 3 && num_resp != 0 {
			num_resp = 0
			s, _ := sm.log.Get(int64(g))
			q := s.(log1).Command
			aresp = append(aresp, Commit{index: g, data: q})

			if sm.commitIndex < g {
				sm.commitIndex = g
			}

			count = 0
		}

	} else if ar.Success == false {
		//fmt.Println(ar)
		var prevx int
		if ar.Term > sm.currentTerm {
			sm.votedFor = 0
			sm.currentTerm = ar.Term
			sm.status = "follower"
			aresp = append(aresp, Alarm{sm.ElectionTimeout})

		} else {
			if sm.nextIndex[ar.Id] > 0 {
				sm.nextIndex[ar.Id]--
			}
			x := sm.nextIndex[ar.Id]

			if prevx > 0 {
				prevx = x - 1
			}
			k, _ := sm.log.Get(int64(prevx))
			t, _ := sm.log.Get(int64(x))
			aresp = append(aresp, send{ar.Id, AppendEntriesReqEv{Term: sm.currentTerm, PrevLogIndex: prevx, PrevLogTerm: k.(log1).LogTerm, Entries: t.(log1).Command, LeaderCommit: sm.commitIndex}})
		}
	}
	//fmt.Println(ar)
	return aresp
}

//AppendEv
type AppendEv struct {
	data1 []byte
}

type Commit struct {
	index int
	data  []byte
	err   string
}

func (ar *AppendEv) AppendF() []interface{} {
	var ret2 = make([]interface{}, 0)
	var k string

	k = "redirect to leader"
	//+ strconv.Itoa(sm.currentLeader)
	//fmt.Println(ar)
	ret2 = append(ret2, Commit{data: ar.data1, err: k})

	return (ret2)

}

func (ar *AppendEv) AppendC() []interface{} {
	var apc = make([]interface{}, 0)
	var k string
	k = "redirect to " + strconv.Itoa(sm.currentLeader)
	apc = append(apc, Commit{data: ar.data1, err: k})
	return (apc)

}

func (ar *AppendEv) AppendL() []interface{} {
	//fmt.Println(sm.nextIndex)
	var l int
	num_resp = 4
	var apc = make([]interface{}, 0)
	sm.lastLogIndex++
	sm.prevLogIndex++
	//u,_:=sm.log.Get(int64(sm.prevLogIndex))
	//sm.prevLogTerm = u.(log1).logTerm

	sm.matchIndex[sm.id]++

	sm.appendIdresponse[sm.id] = 1
	//fmt.Println(sm.prevLogIndex)
	//fmt.Println(sm.lastLogIndex)

	apc = append(apc, LogStore{sm.lastLogIndex, ar.data1})

	for i := 0; i < len(sm.peers); i++ {
		//fmt.Println(sm.nextIndex)
		//fmt.Println(sm.lastLogIndex)
		for l = sm.nextIndex[sm.peers[i]]; l < sm.lastLogIndex; l++ {
			//fmt.Println("Check",0,l,sm.peers[i])
			o, _ := sm.log.Get(int64(l))
			k := l
			if k > 0 {
				k = k - 1
			}
			h, _ := sm.log.Get(int64(k))

			apc = append(apc, send{sm.peers[i], AppendEntriesReqEv{Term: sm.currentTerm, LeaderId: sm.id, PrevLogIndex: h.(log1).LogIndex, Entries: o.(log1).Command, PrevLogTerm: h.(log1).LogTerm, LeaderCommit: sm.commitIndex}})
			//fmt.Println(apc)

		}

		e, _ := sm.log.Get(int64(l - 1))
		apc = append(apc, send{sm.peers[i], AppendEntriesReqEv{Term: sm.currentTerm, LeaderId: sm.id, PrevLogIndex: e.(log1).LogIndex, Entries: ar.data1, PrevLogTerm: e.(log1).LogTerm, LeaderCommit: sm.commitIndex}})

	}

	for i := 0; i < len(sm.peers); i++ {
		sm.nextIndex[sm.peers[i]] = sm.lastLogIndex
	}

	//fmt.Println(apc)
	return apc

}

func (t *TimeoutEv) timeoutHandlerF() []interface{} {
	//fmt.Println("Inside Follower Timeout",sm.id)
	sm.currentTerm++
	sm.status = "candidate"
	sm.voteReceived[sm.id] = 1
	sm.voteCount++
	sm.votedFor = sm.id
	var ret3 = make([]interface{}, 0)
	ret3 = append(ret3, Alarm{sm.ElectionTimeout})
	//fmt.Println(sm.currentTerm)
	for i := 0; i < len(sm.peers); i++ {
		//fmt.Println(i,sm.peers[i])
		ret3 = append(ret3, send{sm.peers[i], VoteRequestEv{Term: sm.currentTerm, LastLogIndex: sm.lastLogIndex, LastLogTerm: sm.lastLogTerm, CandidateId: sm.id}})
	}
	//fmt.Println(ret3)
	//ret3 = append(ret3, Alarm{t.time})

	return (ret3)

}

func random(min, max int) int {
	//fmt.Println("Hey")
	rand.Seed(time.Now().Unix())
	k := max - min
	return rand.Intn(k) + min
}

func (t *TimeoutEv) timeoutHandlerC() []interface{} {
	//fmt.Println("Inside Candidate Timeout",sm.id)
	sm.currentTerm++
	sm.votedFor = 0
	sm.status = "follower"
	empty := [6]int{0, 0, 0, 0, 0, 0}
	sm.voteReceived = empty
	sm.voteCount = 1
	var tout = make([]interface{}, 0)
	tout = append(tout, Alarm{random(sm.ElectionTimeout, 2*sm.ElectionTimeout)})
	for i := 0; i < len(sm.peers); i++ {
		tout = append(tout, send{sm.peers[i], VoteRequestEv{sm.currentTerm, sm.lastLogIndex, sm.lastLogTerm, sm.id}})
	}
	//tout = append(tout, Alarm{t.time})
	return (tout)

}

func (t *TimeoutEv) timeoutHandlerL() []interface{} {
	var tout = make([]interface{}, 0)
	for i := 0; i < 4; i++ {
		//fmt.Println(i)
		tout = append(tout, send{sm.peers[i], AppendEntriesReqEv{Term: sm.currentTerm, LeaderId: sm.id, PrevLogIndex: sm.prevLogIndex, Entries: []byte{}, PrevLogTerm: sm.prevLogTerm, LeaderCommit: sm.commitIndex}})

	}
	return tout

}

//VoteRequest
type VoteRequestEv struct {
	Term         int
	LastLogIndex int
	LastLogTerm  int
	CandidateId  int
}

func (rv *VoteRequestEv) VoterequestF() []interface{} {
	//fmt.Println("Inside VoteReqFollower",sm.id)
	var ret4 = make([]interface{}, 0)
	if rv.Term < sm.currentTerm {
		ret4 = append(ret4, send{rv.CandidateId, VoteRespEv{sm.id, sm.currentTerm, false}})
	} else if rv.Term >= sm.currentTerm && (sm.votedFor == 0 || sm.votedFor == rv.CandidateId) && rv.LastLogIndex >= sm.lastLogIndex {
		sm.currentTerm = rv.Term
		sm.votedFor = rv.CandidateId
		ret4 = append(ret4, send{rv.CandidateId, VoteRespEv{sm.id, sm.currentTerm, true}})
		ret4 = append(ret4, Alarm{sm.ElectionTimeout})

	} else {
		if rv.Term > sm.currentTerm {
			sm.currentTerm = rv.Term
		}
		ret4 = append(ret4, send{rv.CandidateId, VoteRespEv{sm.id, sm.currentTerm, false}})
	}

	return (ret4)
	return (make([]interface{}, 0))
}

func (rv *VoteRequestEv) VoterequestC() []interface{} {

	var ret5 = make([]interface{}, 0)
	if rv.Term < sm.currentTerm {

		ret5 = append(ret5, send{rv.CandidateId, VoteRespEv{sm.id, sm.currentTerm, false}})

		/*	} else if rv.Term >= sm.currentTerm && (sm.votedFor == 0 || sm.votedFor == rv.CandidateId) && rv.LastLogIndex >= sm.lastLogIndex {
			sm.currentTerm = rv.Term
			sm.votedFor = rv.CandidateId
			ret5 = append(ret5, send{rv.CandidateId, VoteRespEv{sm.id, sm.currentTerm, true}})
		*/
	} else {
		if rv.Term > sm.currentTerm {
			sm.currentTerm = rv.Term
		}
		ret5 = append(ret5, send{rv.CandidateId, VoteRespEv{sm.id, sm.currentTerm, false}})

	}
	return (ret5)
	//return(make([]interface{},1))

}

func (rv *VoteRequestEv) VoterequestL() []interface{} {
	var ret5 = make([]interface{}, 0)
	if rv.Term > sm.currentTerm {
		sm.currentTerm = rv.Term
		sm.status = "follower"
		sm.votedFor = 0
		empty1 := [6]int{0, 0, 0, 0, 0, 0}
		sm.appendIdresponse = empty1
		sm.appendresCount = 0
		ret5 = append(ret5, Alarm{sm.ElectionTimeout})
	}
	ret5 = append(ret5, "VOTE RESPONSE TOO LATE")
	return ret5

}

//VoteResponse
type VoteRespEv struct {
	Id          int
	Term        int
	VoteGranted bool
}

func (rvr *VoteRespEv) VoteResponseHandlerF() []interface{} {
	if rvr.Term > sm.currentTerm {
		sm.currentTerm = rvr.Term
		sm.votedFor = 0
	}
	var ret5 = make([]interface{}, 0)
	ret5 = append(ret5, "VOTE RESPONSE TOO LATE")
	return ret5

}

func (rvr *VoteRespEv) VoteResponseHandlerC() []interface{} {
	var vrhc = make([]interface{}, 0)
	if rvr.VoteGranted == true {
		//fmt.Println("voteResponseTrue",sm.id)
		sm.voteReceived[rvr.Id] = 1
		sm.voteCount++
		if sm.voteCount >= 3 {
			sm.status = "leader"
			vrhc = append(vrhc, Alarm{sm.HeartBeatTimeout})
			sm.currentLeader = sm.id
			for i := 0; i < len(sm.peers); i++ {
				vrhc = append(vrhc, send{sm.peers[i], AppendEntriesReqEv{Term: sm.currentTerm, LeaderId: sm.id, PrevLogIndex: sm.lastLogIndex, PrevLogTerm: sm.lastLogTerm, Entries: []byte{}}})
			}

		}

	}
	if rvr.VoteGranted == false {
		//fmt.Println("voteResponseFalse",sm.id)
		sm.voteReceived[rvr.Id] = 0
		sm.votenegCount++
		if sm.votenegCount >= 3 {
			sm.status = "follower"
			sm.votedFor = 0
			vrhc = append(vrhc, Alarm{sm.ElectionTimeout})
		}
		if sm.currentTerm < rvr.Term {
			sm.currentTerm = rvr.Term
		}

	}
	return vrhc
	return (make([]interface{}, 0))
}

func (rvr *VoteRespEv) VoteResponseHandlerL() []interface{} {

	var ret5 = make([]interface{}, 0)
	if rvr.Term > sm.currentTerm {
		sm.currentTerm = rvr.Term
		sm.status = "follower"
		sm.votedFor = 0
		ret5 = append(ret5, Alarm{sm.ElectionTimeout})
		empty2 := [6]int{0, 0, 0, 0, 0, 0}
		sm.appendIdresponse = empty2
		sm.appendresCount = 0
	}

	ret5 = append(ret5, "Error")
	return ret5

}

//Process Event
func (sm1 *StateMachine) ProcessEvent(ev interface{}) []interface{} {
	sm = sm1
	switch ev.(type) {

	case AppendEntriesReqEv:
		cmd := ev.(AppendEntriesReqEv)
		if sm1.status == "follower" {
			return (cmd.AppendReqHandlerF())
		} else if sm1.status == "candidate" {
			return (cmd.AppendReqHandlerC())
		} else if sm1.status == "leader" {
			return (cmd.AppendReqHandlerL())
		}

	case AppendEntriesRespEv:
		cmd := ev.(AppendEntriesRespEv)
		if sm1.status == "follower" {
			return (cmd.AppendEntriesresF())
		} else if sm1.status == "candidate" {
			return (cmd.AppendEntriesresC())
		} else if sm1.status == "leader" {
			//fmt.Println("hello",sm)
			return (cmd.AppendEntriesresL())
		}

	case VoteRequestEv:
		cmd := ev.(VoteRequestEv)
		if sm1.status == "follower" {
			return (cmd.VoterequestF())
		} else if sm1.status == "candidate" {
			return (cmd.VoterequestC())
		} else if sm1.status == "leader" {
			return (cmd.VoterequestL())
		}

		// do stuff with req
		//fmt.Printf("%v\n", cmd)
	case VoteRespEv:
		cmd := ev.(VoteRespEv)
		//fmt.Println("Inside Vote Response")
		//fmt.Println("Hello")
		if sm1.status == "follower" {
			return (cmd.VoteResponseHandlerF())
		} else if sm1.status == "candidate" {
			return (cmd.VoteResponseHandlerC())
		} else if sm1.status == "leader" {
			return (cmd.VoteResponseHandlerL())
		}

	case AppendEv:
		cmd := ev.(AppendEv)
		if sm1.status == "follower" {
			return (cmd.AppendF())
		} else if sm1.status == "candidate" {
			return (cmd.AppendC())
		} else if sm1.status == "leader" {
			return (cmd.AppendL())
		}

	case TimeoutEv:
		cmd := ev.(TimeoutEv)
		if sm1.status == "follower" {
			return (cmd.timeoutHandlerF())
		} else if sm1.status == "candidate" {
			return (cmd.timeoutHandlerC())
		} else if sm1.status == "leader" {
			//fmt.Println(cmd)
			return (cmd.timeoutHandlerL())
		}

		// other cases

	}
	/*go func() []interface{} {

		var hbmsg = make([]interface{}, 0)

		if sm1.status == "leader" {
			for i := 0; i < 4; i++ {
				hbmsg = append(hbmsg, send{sm.peers[i], AppendEntriesReqEv{Term: sm.id, LeaderId: sm.id, LeaderCommit: sm.commitIndex, Entries: []byte{}, PrevLogIndex: sm.prevLogIndex, PrevLogTerm: sm.prevLogTerm}})
			}

		}
		return hbmsg

	}()*/
	return (make([]interface{}, 1))
}

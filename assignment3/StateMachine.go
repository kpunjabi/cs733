package main

import (
	"math"
	"strconv"
)

type log1 struct {
	logIndex int
	logTerm  int
	command  []byte
}

type send struct {
	peerId int
	event  interface{}
}

type LogStore struct {
	index int
	data  []byte
}

//Timeout
type TimeoutEv struct {
	time int
}

//global variable
var count int

type Alarm struct {
	time int
}

type StateMachine struct {
	id               int // server id
	status           string
	peers            []int // other server ids
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
	nextIndex        []int
	matchIndex       []int
	currentLeader    int
	appendresCount   int
	appendIdresponse [6]int
}

//AppendEntries Request
type AppendEntriesReqEv struct {
	term         int
	leaderId     int
	prevLogIndex int
	prevLogTerm  int
	entries      []byte
	leaderCommit int
}

//AppendEntries Response
type AppendEntriesRespEv struct {
	id      int
	term    int
	success bool
}

//var sm *StateMachine
var log [101]log1

//var ret =make([]interface{},1)

func (a *AppendEntriesReqEv) AppendReqHandlerF() []interface{} {

	replyfalse := AppendEntriesRespEv{sm.id, sm.currentTerm, false}
	replytrue := AppendEntriesRespEv{sm.id, sm.currentTerm, true}
	var ret = make([]interface{}, 1)

	if a.term < sm.currentTerm {
		ret = append(ret, send{a.leaderId, replyfalse})

	} else {
		sm.currentLeader = a.leaderId
		if len(a.entries) == 0 {
			ret = append(ret, Alarm{150})
			sm.currentTerm = a.term
			sm.votedFor = 0

		} else {

			if sm.lastLogIndex == a.prevLogIndex && sm.lastLogTerm == a.prevLogTerm {
				// Delete all the log entries following sm.prevLogIndex

				if len(log[sm.lastLogIndex+1].command) != 0 {

					for i := sm.lastLogIndex + 1; log[i].logTerm != 0 && i < 100; i++ {
						//fmt.Println("check2")
						log[i].logTerm = 0
					}
				}

				sm.lastLogIndex++

				sm.currentTerm = a.term
				sm.votedFor = 0

				/*  m:=sm.lastLogIndex

				    //Generate output action : LogStore(index, entries[])
				          //log[m].logIndex=sm.lastLogIndex
				         log[m].logIndex=m
				          log[m].logTerm=a.term

				          log[m].command=a.entries
				*/
				ret = append(ret, LogStore{sm.lastLogIndex, a.entries})

				if a.leaderCommit > sm.commitIndex {

					sm.commitIndex = int(math.Min(float64(a.leaderCommit), float64(sm.lastLogIndex)))
				}

				ret = append(ret, send{a.leaderId, replytrue})

				//fmt.Println("check1")
			} else if len(log[a.prevLogIndex].command) == 0 {

				if a.prevLogIndex == -1 {
					// Generate the output action
					for i := 0; log[i].logTerm != 0 && i < 100; i++ {
						log[i].logTerm = 0
						//log[i].command={}

					}

					ret = append(ret, LogStore{0, a.entries})
					sm.lastLogIndex = 0
					ret = append(ret, send{a.leaderId, replytrue})

				} else {
					//fmt.Println("check0")
					ret = append(ret, send{a.leaderId, replyfalse})
				}

			} else if (sm.lastLogIndex == a.prevLogIndex) && (sm.lastLogTerm != a.prevLogTerm) {
				//Delete the existing entry and those that follows trying to make its log equivalent to that of leader
				ret = append(ret, send{a.leaderId, replyfalse})
			}

		}

	}
	return ret
	return (make([]interface{}, 1))
}

func (ae *AppendEntriesReqEv) AppendReqHandlerC() []interface{} {

	replyfalse := AppendEntriesRespEv{sm.id, sm.currentTerm, false}

	var can = make([]interface{}, 1)
	if ae.term < sm.currentTerm {

		can = append(can, send{ae.leaderId, replyfalse})
	} else {

		sm.status = "follower"
		sm.currentTerm = ae.term
		sm.votedFor = 0

		if len(ae.entries) == 0 {

			can = append(can, send{ae.leaderId, AppendEntriesRespEv{sm.id, sm.currentTerm, true}})

		} else {

			can = append(can, send{ae.leaderId, AppendEntriesRespEv{sm.id, sm.currentTerm, false}})

		}

	}

	return can
}

func (a *AppendEntriesReqEv) AppendReqHandlerL() []interface{} {
	var lead = make([]interface{}, 1)
	if sm.currentTerm <= a.term {
		lead = append(lead, send{a.leaderId, AppendEntriesRespEv{sm.id, sm.currentTerm, false}})
		return lead
	} else {
		sm.currentTerm = a.term
		sm.status = "follower"
		sm.votedFor = 0
	}
	return (make([]interface{}, 1))

}

func (ar *AppendEntriesRespEv) AppendEntriesresF() []interface{} {
	var ret6 = make([]interface{}, 1)
	if ar.term > sm.currentTerm {
		sm.currentTerm = ar.term
		sm.votedFor = 0
	}
	ret6 = append(ret6, "ERROR_OCCURED")
	return (ret6)
	//return(make([]interface{},1))
}

func (ar *AppendEntriesRespEv) AppendEntriesresC() []interface{} {
	var ret7 = make([]interface{}, 1)
	if ar.term > sm.currentTerm {
		sm.currentTerm = ar.term
		sm.status = "follower"
		sm.votedFor = 0
	}
	ret7 = append(ret7, "ERROR_OCCURED")
	return (ret7)

}

func (ar *AppendEntriesRespEv) AppendEntriesresL() []interface{} {
	var aresp = make([]interface{}, 1)
	if ar.success == true {
		sm.appendIdresponse[ar.id] = 1
		sm.appendresCount++

		x := sm.nextIndex[ar.id]
		q := log[x].command
		sm.nextIndex[ar.id]++

		for i := 0; i < 6; i++ {
			if x <= sm.matchIndex[i] {
				count++
			}
		}

		if count >= 2 {
			aresp = append(aresp, Commit{index: x, data: q})

			if sm.commitIndex < x {
				sm.commitIndex = x
			}
			sm.matchIndex[ar.id]++
			count = 0
		} else {

			return (make([]interface{}, 0))
		}
	} else if ar.success == false {
		sm.nextIndex[ar.id]--
		x := sm.nextIndex[ar.id]
		//fmt.Println(x)
		prevx := x - 1
		aresp = append(aresp, send{ar.id, AppendEntriesReqEv{term: sm.currentTerm, prevLogIndex: prevx, prevLogTerm: log[prevx].logTerm, entries: log[x].command, leaderCommit: sm.commitIndex}})

	}

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
	var ret2 = make([]interface{}, 1)
	var k string
	k = "redirect to " + strconv.Itoa(sm.currentLeader)
	ret2 = append(ret2, Commit{data: ar.data1, err: k})
	return (ret2)

}

func (ar *AppendEv) AppendC() []interface{} {
	var apc = make([]interface{}, 1)
	var k string
	k = "redirect to " + strconv.Itoa(sm.currentLeader)
	apc = append(apc, Commit{data: ar.data1, err: k})
	return (apc)

}

func (ar *AppendEv) AppendL() []interface{} {
	var apc = make([]interface{}, 1)
	sm.lastLogIndex++
	sm.prevLogIndex++
	sm.prevLogTerm = log[sm.prevLogIndex].logTerm
	sm.matchIndex[sm.id]++

	sm.appendIdresponse[sm.id] = 1

	apc = append(apc, LogStore{sm.lastLogIndex, ar.data1})
	for i := 0; i < len(sm.peers); i++ {

		sm.nextIndex[sm.peers[i]] = sm.lastLogIndex
		apc = append(apc, send{sm.peers[i], AppendEntriesReqEv{term: sm.currentTerm, leaderId: sm.id, prevLogIndex: sm.prevLogIndex, entries: ar.data1, prevLogTerm: sm.prevLogTerm, leaderCommit: sm.commitIndex}})

	}

	return apc

}

func (t *TimeoutEv) timeoutHandlerF() []interface{} {
	sm.currentTerm++
	sm.status = "candidate"
	sm.voteReceived[sm.id] = 1
	sm.voteCount++
	sm.votedFor = 0
	var ret3 = make([]interface{}, 1)
	for i := 0; i < len(sm.peers); i++ {
		ret3 = append(ret3, send{sm.peers[i], VoteRequestEv{sm.currentTerm, sm.lastLogIndex, sm.lastLogTerm, sm.id}})
	}
	ret3 = append(ret3, Alarm{t.time})
	return (ret3)

}

func (t *TimeoutEv) timeoutHandlerC() []interface{} {
	sm.currentTerm++
	sm.votedFor = 0
	sm.status = "candidate"
	empty := [6]int{0, 0, 1, 0, 0, 0}
	sm.voteReceived = empty
	sm.voteCount = 1
	var tout = make([]interface{}, 1)
	for i := 0; i < len(sm.peers); i++ {
		tout = append(tout, send{sm.peers[i], VoteRequestEv{sm.currentTerm, sm.lastLogIndex, sm.lastLogTerm, sm.id}})
	}
	tout = append(tout, Alarm{t.time})
	return (tout)

}

func (t *TimeoutEv) timeoutHandlerL() []interface{} {
	var tout = make([]interface{}, 1)
	for i := 0; i < 4; i++ {
		//fmt.Println(i)
		tout = append(tout, send{sm.peers[i], AppendEntriesReqEv{term: sm.currentTerm, leaderId: sm.id, prevLogIndex: sm.prevLogIndex, entries: []byte{}, prevLogTerm: sm.prevLogTerm, leaderCommit: sm.commitIndex}})

	}
	return tout

}

//VoteRequest
type VoteRequestEv struct {
	term         int
	lastLogIndex int
	lastLogTerm  int
	candidateId  int
}

func (rv *VoteRequestEv) VoterequestF() []interface{} {

	var ret4 = make([]interface{}, 1)
	if rv.term < sm.currentTerm {
		ret4 = append(ret4, send{rv.candidateId, VoteRespEv{sm.id, sm.currentTerm, false}})
	} else if rv.term >= sm.currentTerm && (sm.votedFor == 0 || sm.votedFor == rv.candidateId) && rv.lastLogIndex >= sm.lastLogIndex {
		sm.currentTerm = rv.term
		sm.votedFor = rv.candidateId
		ret4 = append(ret4, send{rv.candidateId, VoteRespEv{sm.id, sm.currentTerm, true}})

	} else {
		ret4 = append(ret4, send{rv.candidateId, VoteRespEv{sm.id, sm.currentTerm, false}})
	}
	return (ret4)
	return (make([]interface{}, 1))
}

func (rv *VoteRequestEv) VoterequestC() []interface{} {

	var ret5 = make([]interface{}, 1)
	if rv.term < sm.currentTerm {

		ret5 = append(ret5, send{rv.candidateId, VoteRespEv{sm.id, sm.currentTerm, false}})

	} else if rv.term >= sm.currentTerm && (sm.votedFor == 0 || sm.votedFor == rv.candidateId) && rv.lastLogIndex >= sm.lastLogIndex {
		sm.currentTerm = rv.term
		sm.votedFor = rv.candidateId
		ret5 = append(ret5, send{rv.candidateId, VoteRespEv{sm.id, sm.currentTerm, true}})

	} else {
		ret5 = append(ret5, send{rv.candidateId, VoteRespEv{sm.id, sm.currentTerm, false}})

	}
	return (ret5)
	//return(make([]interface{},1))

}

func (rv *VoteRequestEv) VoterequestL() []interface{} {
	var ret5 = make([]interface{}, 1)
	if rv.term > sm.currentTerm {
		sm.currentTerm = rv.term
		sm.status = "follower"
		sm.votedFor = 0
		empty1 := [6]int{0, 0, 1, 0, 0, 0}
		sm.appendIdresponse = empty1
		sm.appendresCount = 0
	}
	ret5 = append(ret5, "VOTE RESPONSE TOO LATE")
	return ret5

}

//VoteResponse
type VoteRespEv struct {
	id          int
	term        int
	voteGranted bool
}

func (rvr *VoteRespEv) VoteResponseHandlerF() []interface{} {
	if rvr.term > sm.currentTerm {
		sm.currentTerm = rvr.term
		sm.votedFor = 0
	}
	var ret5 = make([]interface{}, 1)
	ret5 = append(ret5, "VOTE RESPONSE TOO LATE")
	return ret5

}

func (rvr *VoteRespEv) VoteResponseHandlerC() []interface{} {
	var vrhc = make([]interface{}, 1)
	if rvr.voteGranted == true {
		sm.voteReceived[rvr.id] = 1
		sm.voteCount++
		if sm.voteCount == 3 {
			sm.status = "leader"
			for i := 0; i < len(sm.peers); i++ {
				vrhc = append(vrhc, send{sm.peers[i], AppendEntriesReqEv{term: sm.currentTerm, leaderId: sm.id, prevLogIndex: sm.lastLogIndex, prevLogTerm: sm.lastLogTerm, entries: []byte{}}})
			}

		}
	}
	if rvr.voteGranted == false {
		sm.voteReceived[rvr.id] = 0
		sm.votenegCount++
		if sm.votenegCount == 3 {
			sm.status = "follower"
			sm.votedFor = 0
		}
		if sm.currentTerm < rvr.term {
			sm.currentTerm = rvr.term
		}

	}
	return vrhc
	return (make([]interface{}, 0))
}

func (rvr *VoteRespEv) VoteResponseHandlerL() []interface{} {
	if rvr.term > sm.currentTerm {
		sm.currentTerm = rvr.term
		sm.status = "follower"
		sm.votedFor = 0
		empty2 := [6]int{0, 0, 1, 0, 0, 0}
		sm.appendIdresponse = empty2
		sm.appendresCount = 0
	}
	var ret5 = make([]interface{}, 1)
	ret5 = append(ret5, "Error")
	return ret5

}

//Process Event
func (sm *StateMachine) ProcessEvent(ev interface{}) []interface{} {
	switch ev.(type) {

	case AppendEntriesReqEv:
		cmd := ev.(AppendEntriesReqEv)
		if sm.status == "follower" {
			return (cmd.AppendReqHandlerF())
		} else if sm.status == "candidate" {
			return (cmd.AppendReqHandlerC())
		} else if sm.status == "leader" {
			return (cmd.AppendReqHandlerL())
		}

	case AppendEntriesRespEv:
		cmd := ev.(AppendEntriesRespEv)
		if sm.status == "follower" {
			return (cmd.AppendEntriesresF())
		} else if sm.status == "candidate" {
			return (cmd.AppendEntriesresC())
		} else if sm.status == "leader" {
			return (cmd.AppendEntriesresL())
		}

	case VoteRequestEv:
		cmd := ev.(VoteRequestEv)
		if sm.status == "follower" {
			return (cmd.VoterequestF())
		} else if sm.status == "candidate" {
			return (cmd.VoterequestC())
		} else if sm.status == "leader" {
			return (cmd.VoterequestL())
		}

		// do stuff with req
		//fmt.Printf("%v\n", cmd)
	case VoteRespEv:
		cmd := ev.(VoteRespEv)
		if sm.status == "follower" {
			return (cmd.VoteResponseHandlerF())
		} else if sm.status == "candidate" {
			return (cmd.VoteResponseHandlerC())
		} else if sm.status == "leader" {
			return (cmd.VoteResponseHandlerL())
		}

	case AppendEv:
		cmd := ev.(AppendEv)
		if sm.status == "follower" {
			return (cmd.AppendF())
		} else if sm.status == "candidate" {
			return (cmd.AppendC())
		} else if sm.status == "leader" {
			return (cmd.AppendL())
		}

	case TimeoutEv:
		cmd := ev.(TimeoutEv)
		if sm.status == "follower" {
			return (cmd.timeoutHandlerF())
		} else if sm.status == "candidate" {
			return (cmd.timeoutHandlerC())
		} else if sm.status == "leader" {
			//fmt.Println(cmd)
			return (cmd.timeoutHandlerL())
		}

		// other cases
		//default: println ("Unrecognized")
	}
	return (make([]interface{}, 1))
}

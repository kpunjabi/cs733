package main

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

//type actions []interface{}
var sm *StateMachine

//var log *log

func TestFollowerAppendEntriesReqEv(t *testing.T) {

	fmt.Println("\nTestFollowerAppendEntriesReqEv:Tested\n")

	en := []byte("read")
	p1 := []int{2, 3, 4, 5}
	var out1 = make([]interface{}, 1)
	var out2 = make([]interface{}, 1)
	var out3 = make([]interface{}, 1)
	var out4 = make([]interface{}, 1)

	//checking when leader's term>follower's current term
	sm = &StateMachine{id: 1, status: "follower", peers: p1, currentTerm: 3, lastLogIndex: 1, lastLogTerm: 3}
	actions := sm.ProcessEvent(AppendEntriesReqEv{term: 10, leaderId: 4, prevLogIndex: 1, prevLogTerm: 3, entries: en, leaderCommit: 6})
	out1 = append(out1, LogStore{2, en})
	out1 = append(out1, send{4, AppendEntriesRespEv{1, 3, true}})
	expectedStr(t, out1, actions, "Test1: 1")

	//checking for reciever's current term > leader's current term
	sm = &StateMachine{id: 2, status: "follower", peers: p1, currentTerm: 5, lastLogIndex: 10, lastLogTerm: 5}
	actions = sm.ProcessEvent(AppendEntriesReqEv{term: 4, leaderId: 4, prevLogIndex: 10, prevLogTerm: 3, entries: en, leaderCommit: 9})
	out2 = append(out2, send{4, AppendEntriesRespEv{2, 5, false}})
	expectedStr(t, out2, actions, "Test1: 2")

	//checking for prevIndex=-1
	sm = &StateMachine{id: 2, status: "follower", peers: p1, currentTerm: 1, lastLogIndex: 0, lastLogTerm: 0}
	actions = sm.ProcessEvent(AppendEntriesReqEv{term: 2, leaderId: 4, prevLogIndex: 5, prevLogTerm: 1, entries: en, leaderCommit: 4})
	out3 = append(out3, send{4, AppendEntriesRespEv{2, 1, false}})
	expectedStr(t, out3, actions, "Test1: 3")

	//checking for heart beat message
	sm = &StateMachine{id: 2, status: "follower", peers: p1, currentTerm: 1, lastLogIndex: 2, lastLogTerm: 1}
	actions = sm.ProcessEvent(AppendEntriesReqEv{term: 4, leaderId: 4, prevLogIndex: 5, prevLogTerm: 3, entries: []byte{}, leaderCommit: 4})
	out4 = append(out4, Alarm{150})
	expectedStr(t, out4, actions, "Test1: 4")

}

func TestFollowerVoteRequest(t *testing.T) {
	fmt.Println("\nTestFollowerVoteRequest:Tested\n")
	p2 := []int{2, 3, 4, 5}
	en := []byte("read")
	var out1 = make([]interface{}, 1)
	var out2 = make([]interface{}, 1)
	var out3 = make([]interface{}, 1)
	var out4 = make([]interface{}, 1)

	sm = &StateMachine{id: 2, status: "follower", votedFor: 0, peers: p2, currentTerm: 2, lastLogIndex: 3, lastLogTerm: 2}
	actions1 := sm.ProcessEvent(VoteRequestEv{term: 5, lastLogIndex: 6, lastLogTerm: 4, candidateId: 4})
	out1 = append(out1, send{4, VoteRespEv{2, 5, true}})
	expectedStr(t, out1, actions1, "Test2: 1")

	en = []byte("jump")
	p2 = []int{2, 3, 1, 5}
	sm = &StateMachine{id: 4, status: "follower", votedFor: 0, peers: p2, currentTerm: 4, lastLogIndex: 5, lastLogTerm: 4}
	actions1 = sm.ProcessEvent(VoteRequestEv{term: 3, lastLogIndex: 4, lastLogTerm: 2, candidateId: 3})
	out2 = append(out2, send{3, VoteRespEv{4, 4, false}})
	expectedStr(t, out2, actions1, "Test2: 2")

	en = []byte("add")
	p2 = []int{1, 3, 4, 5}
	sm = &StateMachine{id: 2, status: "follower", votedFor: 5, peers: p2, currentTerm: 2, lastLogIndex: 4, lastLogTerm: 2}
	actions1 = sm.ProcessEvent(VoteRequestEv{term: 3, lastLogIndex: 6, lastLogTerm: 2, candidateId: 5})
	out3 = append(out3, send{5, VoteRespEv{2, 3, true}})
	expectedStr(t, out3, actions1, "Test2: 3")

	en = []byte("write")
	p2 = []int{2, 1, 4, 5}
	sm = &StateMachine{id: 3, status: "follower", votedFor: 3, peers: p2, currentTerm: 3, lastLogIndex: 5, lastLogTerm: 3}
	actions1 = sm.ProcessEvent(VoteRequestEv{term: 5, lastLogIndex: 8, lastLogTerm: 4, candidateId: 4})
	out4 = append(out4, send{4, VoteRespEv{3, 3, false}})
	expectedStr(t, out4, actions1, "Test2: 4")
}

func TestFollowerAppendEntriesResp(t *testing.T) {
	fmt.Println("\nTestFollowerAppendEntriesResp:Tested\n")
	p3 := []int{2, 3, 4, 5}
	var out1 = make([]interface{}, 1)
	var out2 = make([]interface{}, 1)

	sm = &StateMachine{id: 1, status: "follower", votedFor: 0, peers: p3, currentTerm: 3, lastLogIndex: 5, lastLogTerm: 2}
	actions2 := sm.ProcessEvent(AppendEntriesRespEv{id: 4, term: 5, success: true})
	out1 = append(out1, "ERROR_OCCURED")
	expectedStr(t, out1, actions2, "Test3: 1")
	expect(t, strconv.Itoa(sm.currentTerm), strconv.Itoa(5), "AppendEntriesResptermincrement")

	sm = &StateMachine{id: 1, status: "follower", votedFor: 0, peers: p3, currentTerm: 3, lastLogIndex: 5, lastLogTerm: 2}
	actions2 = sm.ProcessEvent(AppendEntriesRespEv{id: 3, term: 2, success: false})
	out2 = append(out2, "ERROR_OCCURED")
	expectedStr(t, out2, actions2, "Test3: 2")
	expect(t, strconv.Itoa(sm.currentTerm), strconv.Itoa(3), "AppendEntriesResptermincrement")

}

func TestFollowerVoteResp(t *testing.T) {

	fmt.Println("\nTestFollowerVoteResp:Tested\n")

	p4 := []int{2, 3, 4, 5}
	var out1 = make([]interface{}, 1)
	var out2 = make([]interface{}, 1)

	sm = &StateMachine{id: 1, status: "follower", votedFor: 0, peers: p4, currentTerm: 3, lastLogIndex: 3, lastLogTerm: 2}
	actions5 := sm.ProcessEvent(VoteRespEv{id: 1, term: 2, voteGranted: false})
	out1 = append(out1, "VOTE RESPONSE TOO LATE")
	expectedStr(t, out1, actions5, "Test4: 1")
	expect(t, strconv.Itoa(sm.currentTerm), strconv.Itoa(3), "voteResptermincrement")

	sm = &StateMachine{id: 1, status: "follower", votedFor: 0, peers: p4, currentTerm: 3, lastLogIndex: 3, lastLogTerm: 2}
	actions5 = sm.ProcessEvent(VoteRespEv{id: 1, term: 5, voteGranted: true})
	out2 = append(out2, "VOTE RESPONSE TOO LATE")
	expectedStr(t, out2, actions5, "Test4: 2")
	expect(t, strconv.Itoa(sm.currentTerm), strconv.Itoa(5), "voteResptermincrement")

}

func TestFollowerTimeout(t *testing.T) {

	fmt.Println("\nTestFollowerTimeout:Tested\n")
	p5 := []int{2, 3, 4, 5}
	var out1 = make([]interface{}, 1)

	sm = &StateMachine{id: 1, status: "follower", votedFor: 0, peers: p5, currentTerm: 2, lastLogIndex: 3, lastLogTerm: 2}
	actions4 := sm.ProcessEvent(TimeoutEv{time: 100})
	for i := 0; i < len(sm.peers); i++ {
		out1 = append(out1, send{sm.peers[i], VoteRequestEv{3, 3, 2, 1}})
	}
	out1 = append(out1, Alarm{100})
	expectedStr(t, out1, actions4, "Test5: 1")
	expect(t, strconv.Itoa(sm.currentTerm), strconv.Itoa(3), "Timertermincrement")
	expect(t, sm.status, "candidate", "Timerstatechange")
	expect(t, strconv.Itoa(sm.voteCount), strconv.Itoa(1), "Timervotecountincrement")

}

func TestFollowerAppend(t *testing.T) {
	fmt.Println("\nTestFollowerAppend:Tested\n")
	p6 := []int{2, 3, 4, 5}
	var k string
	d := []byte("hello")
	var out1 = make([]interface{}, 1)
	sm = &StateMachine{id: 1, status: "follower", votedFor: 0, peers: p6, currentTerm: 2, lastLogIndex: 3, lastLogTerm: 2}
	actions3 := sm.ProcessEvent(AppendEv{data1: d})
	k = "redirect to " + strconv.Itoa(sm.currentLeader)
	out1 = append(out1, Commit{data: d, err: k})
	expectedStr(t, out1, actions3, "Test6: 1")

}

//Candidate

func TestCandidateAppendEntriesReqEv(t *testing.T) {
	var out1 = make([]interface{}, 1)
	var out2 = make([]interface{}, 1)
	var out3 = make([]interface{}, 1)
	fmt.Println("\nTestCandidateAppendEntriesReqEv:Tested\n")
	peer := []int{2, 3, 4, 5}
	sm = &StateMachine{id: 1, status: "candidate", votedFor: 0, peers: peer, currentTerm: 3, lastLogIndex: 6, lastLogTerm: 2}
	actions3 := sm.ProcessEvent(AppendEntriesReqEv{term: 2, leaderId: 4, prevLogIndex: 5, prevLogTerm: 2})
	out1 = append(out1, send{4, AppendEntriesRespEv{1, 3, false}})
	expectedStr(t, out1, actions3, "Test7: 1")

	sm = &StateMachine{id: 1, status: "candidate", votedFor: 0, peers: peer, currentTerm: 3, lastLogIndex: 6, lastLogTerm: 2}
	actions3 = sm.ProcessEvent(AppendEntriesReqEv{term: 5, leaderId: 3, prevLogIndex: 7, entries: []byte{}, prevLogTerm: 4})
	out2 = append(out2, send{3, AppendEntriesRespEv{1, 5, true}})
	expectedStr(t, out2, actions3, "Test7: 2")

	sm = &StateMachine{id: 1, status: "candidate", votedFor: 0, peers: peer, currentTerm: 3, lastLogIndex: 6, lastLogTerm: 2}
	actions3 = sm.ProcessEvent(AppendEntriesReqEv{term: 5, leaderId: 4, prevLogIndex: 7, entries: []byte("read"), prevLogTerm: 4})
	out3 = append(out3, send{4, AppendEntriesRespEv{1, 5, false}})
	expectedStr(t, out3, actions3, "Test7: 3")
	expect(t, strconv.Itoa(sm.currentTerm), strconv.Itoa(5), "AppendEntriesReqtermincrement")
	expect(t, sm.status, "follower", "AppendEntriesReqstatechange")

}

func TestCandidateAppendEntriesResp(t *testing.T) {
	fmt.Println("\nTestCandidateAppendEntriesResp:Tested\n")
	p3 := []int{2, 3, 4, 5}
	var out1 = make([]interface{}, 1)
	var out2 = make([]interface{}, 1)

	sm = &StateMachine{id: 1, status: "candidate", votedFor: 0, peers: p3, currentTerm: 3, lastLogIndex: 5, lastLogTerm: 2}
	actions2 := sm.ProcessEvent(AppendEntriesRespEv{id: 2, term: 5, success: true})
	out1 = append(out1, "ERROR_OCCURED")
	expectedStr(t, out1, actions2, "Test8: 1")
	expect(t, strconv.Itoa(sm.currentTerm), strconv.Itoa(5), "AppendEntriesResptermincrement")

	sm = &StateMachine{id: 1, status: "candidate", votedFor: 0, peers: p3, currentTerm: 3, lastLogIndex: 5, lastLogTerm: 2}
	actions2 = sm.ProcessEvent(AppendEntriesRespEv{id: 3, term: 2, success: false})
	out2 = append(out2, "ERROR_OCCURED")
	expectedStr(t, out2, actions2, "Test8: 2")
	expect(t, strconv.Itoa(sm.currentTerm), strconv.Itoa(3), "AppendEntriesResptermincrement")

}

func TestCandidateTimeout(t *testing.T) {

	fmt.Println("\nTestCandidateTimeout:Tested\n")
	p5 := []int{2, 1, 4, 5}
	var out1 = make([]interface{}, 1)
	arr1 := [6]int{0, 0, 1, 0, 0, 0}
	sm = &StateMachine{id: 3, status: "candidate", votedFor: 0, peers: p5, currentTerm: 5, lastLogIndex: 9, lastLogTerm: 4}
	actions4 := sm.ProcessEvent(TimeoutEv{time: 100})
	for i := 0; i < len(sm.peers); i++ {
		out1 = append(out1, send{sm.peers[i], VoteRequestEv{sm.currentTerm, sm.lastLogIndex, sm.lastLogTerm, sm.id}})
	}
	out1 = append(out1, Alarm{100})
	expectedStr(t, out1, actions4, "Test9: 1")
	expect(t, strconv.Itoa(sm.currentTerm), strconv.Itoa(6), "Timertermincrement")
	expect(t, sm.status, "candidate", "Timerstatechange")
	expect(t, strconv.Itoa(sm.voteCount), strconv.Itoa(1), "Timervotecountincrement")
	//fmt.Println(sm.voteReceived)
	expectArr(t, arr1, sm.voteReceived, "voteRecievedarray")

}

func TestCandidateVoteRequest(t *testing.T) {
	fmt.Println("\nTestCandidateVoteRequest:Tested\n")
	p2 := []int{2, 3, 4, 5}
	var out1 = make([]interface{}, 1)
	var out2 = make([]interface{}, 1)
	var out3 = make([]interface{}, 1)
	var out4 = make([]interface{}, 1)

	sm = &StateMachine{id: 2, status: "candidate", votedFor: 0, peers: p2, currentTerm: 2, lastLogIndex: 3, lastLogTerm: 2}
	actions1 := sm.ProcessEvent(VoteRequestEv{term: 5, lastLogIndex: 6, lastLogTerm: 4, candidateId: 4})
	out1 = append(out1, send{4, VoteRespEv{2, 5, true}})
	expectedStr(t, out1, actions1, "Test10: 1")

	p2 = []int{2, 3, 1, 5}
	sm = &StateMachine{id: 4, status: "candidate", votedFor: 0, peers: p2, currentTerm: 4, lastLogIndex: 5, lastLogTerm: 4}
	actions1 = sm.ProcessEvent(VoteRequestEv{term: 3, lastLogIndex: 4, lastLogTerm: 2, candidateId: 3})
	out2 = append(out2, send{3, VoteRespEv{4, 4, false}})
	expectedStr(t, out2, actions1, "Test10: 2")

	p2 = []int{1, 3, 4, 5}
	sm = &StateMachine{id: 2, status: "candidate", votedFor: 5, peers: p2, currentTerm: 2, lastLogIndex: 4, lastLogTerm: 2}
	actions1 = sm.ProcessEvent(VoteRequestEv{term: 3, lastLogIndex: 6, lastLogTerm: 2, candidateId: 5})
	out3 = append(out3, send{5, VoteRespEv{2, 3, true}})
	expectedStr(t, out3, actions1, "Test10: 3")

	p2 = []int{2, 1, 4, 5}
	sm = &StateMachine{id: 3, status: "candidate", votedFor: 3, peers: p2, currentTerm: 3, lastLogIndex: 5, lastLogTerm: 3}
	actions1 = sm.ProcessEvent(VoteRequestEv{term: 3, lastLogIndex: 8, lastLogTerm: 4, candidateId: 4})
	out4 = append(out4, send{4, VoteRespEv{3, 3, false}})
	expectedStr(t, out4, actions1, "Test10: 4")

}

func TestCandidateVoteResp(t *testing.T) {
	fmt.Println("\nTestCandidateVoteResp:Tested\n")
	p2 := []int{1, 3, 4, 5}
	var out1 = make([]interface{}, 1)
	var out2 = make([]interface{}, 1)
	var out3 = make([]interface{}, 1)
	var out4 = make([]interface{}, 1)
	var out5 = make([]interface{}, 1)
	var out6 = make([]interface{}, 1)
	var out7 = make([]interface{}, 1)
	var out8 = make([]interface{}, 1)

	arr1 := [6]int{0, 1, 1, 0, 0, 0}

	sm = &StateMachine{id: 2, status: "candidate", votedFor: 0, peers: p2, currentTerm: 3, lastLogIndex: 3, lastLogTerm: 2}
	sm.voteReceived[sm.id] = 1
	sm.voteCount = 1

	actions5 := sm.ProcessEvent(VoteRespEv{id: 1, term: 1, voteGranted: true})
	arr1 = [6]int{0, 1, 1, 0, 0, 0}
	expectArr(t, arr1, sm.voteReceived, "arrtest11:1")
	expectedStr(t, out1, actions5, "Test11: 1")
	expect(t, strconv.Itoa(sm.currentTerm), strconv.Itoa(3), "voteResptermincrement")
	expect(t, sm.status, "candidate", "AppendEntriesReqstatechange")

	actions5 = sm.ProcessEvent(VoteRespEv{id: 3, term: 2, voteGranted: false})
	arr1 = [6]int{0, 1, 1, 0, 0, 0}
	expectArr(t, arr1, sm.voteReceived, "arrtest11:2")
	expectedStr(t, out2, actions5, "Test11: 2")
	expect(t, strconv.Itoa(sm.currentTerm), strconv.Itoa(3), "voteResptermincrement")
	expect(t, sm.status, "candidate", "AppendEntriesReqstatechange")

	actions5 = sm.ProcessEvent(VoteRespEv{id: 4, term: 2, voteGranted: true})
	arr1 = [6]int{0, 1, 1, 0, 1, 0}
	expectArr(t, arr1, sm.voteReceived, "arrtest11:3")
	for i := 0; i < len(sm.peers); i++ {
		out3 = append(out3, send{sm.peers[i], AppendEntriesReqEv{term: sm.currentTerm, leaderId: sm.id, prevLogIndex: sm.lastLogIndex, prevLogTerm: sm.lastLogTerm, entries: []byte{}}})
	}
	expectedStr(t, out3, actions5, "Test11: 3")
	expect(t, strconv.Itoa(sm.currentTerm), strconv.Itoa(3), "voteResptermincrement")
	expect(t, sm.status, "leader", "VoteReqEvstatechange")

	actions5 = sm.ProcessEvent(VoteRespEv{id: 5, term: 1, voteGranted: true})
	//fmt.Println(sm.voteReceived)
	arr1 = [6]int{0, 1, 1, 0, 1, 0}
	expectArr(t, arr1, sm.voteReceived, "arrtest11:4")
	out4 = append(out4, "Error")
	expectedStr(t, out4, actions5, "Test11: 4")
	expect(t, strconv.Itoa(sm.currentTerm), strconv.Itoa(3), "voteResptermincrement")
	expect(t, sm.status, "leader", "AppendEntriesReqstatechange")

	//Testing majority false case
	p1 := []int{1, 3, 4, 5}
	sm = &StateMachine{id: 3, status: "candidate", votedFor: 0, peers: p1, currentTerm: 5, lastLogIndex: 9, lastLogTerm: 4}
	sm.voteReceived[sm.id] = 1
	sm.voteCount = 1

	expect(t, strconv.Itoa(sm.votenegCount), strconv.Itoa(0), "negative vote count increment")
	actions6 := sm.ProcessEvent(VoteRespEv{id: 1, term: 3, voteGranted: true})
	arr1 = [6]int{0, 1, 0, 1, 0, 0}
	expectArr(t, arr1, sm.voteReceived, "arrtest11:5")

	expectedStr(t, out5, actions6, "Test11: 5")
	expect(t, sm.status, "candidate", "AppendEntriesReqstatechange")

	actions6 = sm.ProcessEvent(VoteRespEv{id: 1, term: 4, voteGranted: false})
	arr1 = [6]int{0, 0, 0, 1, 0, 0}
	expectArr(t, arr1, sm.voteReceived, "arrtest11:6")
	expectedStr(t, out6, actions6, "Test11: 6")
	expect(t, strconv.Itoa(sm.currentTerm), strconv.Itoa(5), "voteResptermincrement")
	expect(t, sm.status, "candidate", "AppendEntriesReqstatechange")
	expect(t, strconv.Itoa(sm.votenegCount), strconv.Itoa(1), "negative vote count increment")

	actions6 = sm.ProcessEvent(VoteRespEv{id: 1, term: 1, voteGranted: false})
	arr1 = [6]int{0, 0, 0, 1, 0, 0}
	expectArr(t, arr1, sm.voteReceived, "arrtest11:7")

	expectedStr(t, out7, actions6, "Test11: 7")
	expect(t, strconv.Itoa(sm.currentTerm), strconv.Itoa(5), "voteResptermincrement")
	expect(t, sm.status, "candidate", "AppendEntriesReqstatechange")
	expect(t, strconv.Itoa(sm.votenegCount), strconv.Itoa(2), "negative vote count increment")
	actions6 = sm.ProcessEvent(VoteRespEv{id: 1, term: 6, voteGranted: false})
	arr1 = [6]int{0, 0, 0, 1, 0, 0}
	expectArr(t, arr1, sm.voteReceived, "arrtest11:8")

	expectedStr(t, out8, actions6, "Test11: 8")
	expect(t, strconv.Itoa(sm.currentTerm), strconv.Itoa(6), "voteResptermincrement")
	expect(t, sm.status, "follower", "AppendEntriesReqstatechange")
	expect(t, strconv.Itoa(sm.votenegCount), strconv.Itoa(3), "negative vote count increment")

}

func TestCandidateAppend(t *testing.T) {

	fmt.Println("\nTestCandidateAppend:Tested\n")
	p6 := []int{2, 3, 4, 5}
	d := []byte("read")
	var out1 = make([]interface{}, 1)
	var k string

	sm = &StateMachine{id: 1, status: "candidate", votedFor: 0, peers: p6, currentTerm: 2, lastLogIndex: 3, lastLogTerm: 2}
	actions3 := sm.ProcessEvent(AppendEv{data1: d})
	k = "redirect to " + strconv.Itoa(sm.currentLeader)
	out1 = append(out1, Commit{data: d, err: k})
	expectedStr(t, out1, actions3, "Test12: 1")
}

//Leader

func TestLeaderAppendEntriesReqEv(t *testing.T) {
	fmt.Println("\nTestLeaderAppendEntriesReqEv:Tested\n")
	p6 := []int{1, 3, 4, 5}
	d := []byte("read")
	var out1 = make([]interface{}, 1)
	var out2 = make([]interface{}, 1)

	sm = &StateMachine{id: 2, status: "leader", votedFor: 0, peers: p6, currentTerm: 5, lastLogIndex: 7, lastLogTerm: 5, commitIndex: 5}
	actions := sm.ProcessEvent(AppendEntriesReqEv{term: 4, leaderId: 4, prevLogIndex: 5, prevLogTerm: 3, entries: d, leaderCommit: 3})
	expectedStr(t, out1, actions, "Test13: 1")
	expect(t, strconv.Itoa(sm.currentTerm), strconv.Itoa(4), "appendEntryReqleaderterm")
	expect(t, sm.status, "follower", "AppendEntriesReqstatechange")

	sm = &StateMachine{id: 2, status: "leader", votedFor: 0, peers: p6, currentTerm: 5, lastLogIndex: 7, lastLogTerm: 5, commitIndex: 5}
	actions = sm.ProcessEvent(AppendEntriesReqEv{term: 6, leaderId: 4, prevLogIndex: 5, prevLogTerm: 3, entries: d, leaderCommit: 3})
	expect(t, strconv.Itoa(sm.currentTerm), strconv.Itoa(5), "appendEntryReqleaderterm")
	expect(t, sm.status, "leader", "AppendEntriesReqstatechange")
	out2 = append(out2, send{4, AppendEntriesRespEv{sm.id, sm.currentTerm, false}})
	expectedStr(t, out2, actions, "Test13: 2")
}

func TestLeaderAppendEntriesResp(t *testing.T) {
	fmt.Println("\nTestLeaderAppendEntriesRespEv:Tested\n")
	p6 := []int{1, 3, 4, 5}
	mI := []int{0, 3, 6, 4, 3, 5}
	nI := []int{0, 4, 7, 5, 4, 6}
	//log[0]=struct{0,1,[]byte("read")}
	cmd0 := []byte("read")
	cmd1 := []byte("add")
	cmd2 := []byte("write")
	cmd3 := []byte("delete")
	log[0] = log1{logIndex: 0, logTerm: 1, command: cmd1}
	log[1] = log1{logIndex: 1, logTerm: 2, command: cmd0}
	log[2] = log1{logIndex: 2, logTerm: 2, command: cmd1}
	log[3] = log1{logIndex: 3, logTerm: 3, command: cmd0}
	log[4] = log1{logIndex: 4, logTerm: 4, command: cmd2}
	log[5] = log1{logIndex: 5, logTerm: 5, command: cmd3}
	log[6] = log1{logIndex: 6, logTerm: 5, command: cmd1}
	var out1 = make([]interface{}, 1)
	var out2 = make([]interface{}, 1)
	var out3 = make([]interface{}, 1)
	var out4 = make([]interface{}, 0)

	sm = &StateMachine{id: 2, status: "leader", votedFor: 0, peers: p6, currentTerm: 6, matchIndex: mI, nextIndex: nI, lastLogIndex: 7, lastLogTerm: 5, commitIndex: 3}
	actions2 := sm.ProcessEvent(AppendEntriesRespEv{id: 3, term: 5, success: false})
	//fmt.Println(actions2[1])
	prevx := sm.nextIndex[3] - 1
	x := prevx + 1
	out1 = append(out1, send{3, AppendEntriesReqEv{term: sm.currentTerm, prevLogIndex: prevx, prevLogTerm: log[prevx].logTerm, entries: log[x].command, leaderCommit: sm.commitIndex}})
	expectedStr(t, out1, actions2, "Test14: 1")
	expect(t, strconv.Itoa(sm.nextIndex[3]), strconv.Itoa(4), "nextIndex[3]")
	expect(t, strconv.Itoa(sm.commitIndex), strconv.Itoa(3), "leaderCommitIndex")
	expect(t, strconv.Itoa(sm.matchIndex[3]), strconv.Itoa(4), "matchIndex[3]")

	actions2 = sm.ProcessEvent(AppendEntriesRespEv{id: 1, term: 3, success: false})
	//fmt.Println(actions2[1])
	prevx = sm.nextIndex[1] - 1
	x = prevx + 1
	out2 = append(out2, send{1, AppendEntriesReqEv{term: sm.currentTerm, prevLogIndex: prevx, prevLogTerm: log[prevx].logTerm, entries: log[x].command, leaderCommit: sm.commitIndex}})
	expectedStr(t, out2, actions2, "Test14: 2")
	expect(t, strconv.Itoa(sm.nextIndex[1]), strconv.Itoa(3), "nextIndex[1]")
	expect(t, strconv.Itoa(sm.commitIndex), strconv.Itoa(3), "leaderCommitIndex")
	expect(t, strconv.Itoa(sm.matchIndex[1]), strconv.Itoa(3), "matchIndex[1]")

	actions2 = sm.ProcessEvent(AppendEntriesRespEv{id: 4, term: 3, success: true})
	x = sm.nextIndex[4] - 1
	q := log[x].command
	out3 = append(out3, Commit{index: x, data: q})
	expectedStr(t, out3, actions2, "Test14: 3")
	expect(t, strconv.Itoa(sm.nextIndex[4]), strconv.Itoa(5), "nextIndex[4]")
	expect(t, strconv.Itoa(sm.commitIndex), strconv.Itoa(4), "leaderCommitIndex")
	expect(t, strconv.Itoa(sm.matchIndex[4]), strconv.Itoa(4), "matchIndex[4]")

	actions2 = sm.ProcessEvent(AppendEntriesRespEv{id: 5, term: 5, success: true})
	expectedStr(t, out4, actions2, "Test14: 4")
	expect(t, strconv.Itoa(sm.nextIndex[5]), strconv.Itoa(7), "nextIndex[5]")
	expect(t, strconv.Itoa(sm.commitIndex), strconv.Itoa(4), "leaderCommitIndex")
	expect(t, strconv.Itoa(sm.matchIndex[5]), strconv.Itoa(5), "matchIndex[5]")

}

func TestLeaderTimeout(t *testing.T) {

	fmt.Println("\nTestLeaderTimeout:Tested\n")
	p5 := []int{2, 1, 4, 5}
	var out1 = make([]interface{}, 1)

	sm = &StateMachine{id: 3, status: "leader", votedFor: 0, peers: p5, currentTerm: 5, lastLogIndex: 9, lastLogTerm: 4}
	actions4 := sm.ProcessEvent(TimeoutEv{time: 100})
	for i := 0; i < 4; i++ {

		out1 = append(out1, send{sm.peers[i], AppendEntriesReqEv{term: sm.currentTerm, leaderId: sm.id, prevLogIndex: sm.prevLogIndex, entries: []byte{}, prevLogTerm: sm.prevLogTerm, leaderCommit: sm.commitIndex}})

	}

	expectedStr(t, out1, actions4, "Test15: 1")
}

func TestLeaderVoteRequest(t *testing.T) {
	fmt.Println("\nTestLeaderVoteRequest:Tested\n")
	p6 := []int{1, 2, 3, 5}
	//d :=[]byte("read")
	var out1 = make([]interface{}, 1)
	var out2 = make([]interface{}, 1)

	sm = &StateMachine{id: 4, status: "leader", votedFor: 0, peers: p6, currentTerm: 5, lastLogIndex: 7, lastLogTerm: 5, commitIndex: 5}
	actions := sm.ProcessEvent(VoteRequestEv{term: 4, lastLogIndex: 5, lastLogTerm: 3, candidateId: 5})
	out1 = append(out1, "VOTE RESPONSE TOO LATE")
	expectedStr(t, out1, actions, "Test16: 1")

	p6 = []int{4, 2, 3, 5}
	sm = &StateMachine{id: 1, status: "leader", votedFor: 0, peers: p6, currentTerm: 5, lastLogIndex: 7, lastLogTerm: 5, commitIndex: 5}
	actions = sm.ProcessEvent(VoteRequestEv{term: 7, lastLogIndex: 10, lastLogTerm: 6, candidateId: 3})
	out2 = append(out2, "VOTE RESPONSE TOO LATE")
	expectedStr(t, out2, actions, "Test16: 2")
}

func TestLeaderVoteResp(t *testing.T) {
	fmt.Println("\nTestLeaderVoteResp:Tested\n")
	p6 := []int{4, 2, 1, 5}
	var out1 = make([]interface{}, 1)
	var out2 = make([]interface{}, 1)

	sm = &StateMachine{id: 3, status: "leader", votedFor: 0, peers: p6, currentTerm: 6, lastLogIndex: 11, lastLogTerm: 5, commitIndex: 8}
	actions := sm.ProcessEvent(VoteRespEv{id: 5, term: 4, voteGranted: true})
	out1 = append(out1, "Error")
	expectedStr(t, out1, actions, "Test17: 1")

	sm = &StateMachine{id: 2, status: "leader", votedFor: 0, peers: p6, currentTerm: 5, lastLogIndex: 7, lastLogTerm: 5, commitIndex: 5}
	actions = sm.ProcessEvent(VoteRespEv{id: 2, term: 4, voteGranted: false})
	out2 = append(out2, "Error")
	expectedStr(t, out2, actions, "Test17: 2")

}

func TestLeaderAppend(t *testing.T) {
	fmt.Println("\nTestLeaderAppend:Tested\n")
	p6 := []int{2, 3, 4, 5}
	mI := []int{0, 3, 6, 4, 3, 5}
	nI := []int{0, 0, 7, 5, 4, 6}
	d := []byte("hello")
	var out1 = make([]interface{}, 1)

	sm = &StateMachine{id: 1, status: "leader", votedFor: 0, peers: p6, matchIndex: mI, nextIndex: nI, currentTerm: 5, lastLogIndex: 10, lastLogTerm: 5}
	actions3 := sm.ProcessEvent(AppendEv{data1: d})
	out1 = append(out1, LogStore{index: sm.lastLogIndex, data: d})
	for i := 0; i < len(sm.peers); i++ {

		expect(t, strconv.Itoa(sm.nextIndex[sm.peers[i]]), strconv.Itoa(11), "nextIndexincrement")

		out1 = append(out1, send{sm.peers[i], AppendEntriesReqEv{term: sm.currentTerm, leaderId: sm.id, prevLogIndex: sm.prevLogIndex, entries: d, prevLogTerm: sm.prevLogTerm, leaderCommit: sm.commitIndex}})

	}
	expectedStr(t, out1, actions3, "Test18: 1")

}

func expect(t *testing.T, a string, b string, c string) {
	if a != b {
		t.Error(fmt.Sprintf("Expected %v, found %v  Error type:  %v", b, a, c)) // t.Error is visible when running `go test -verbose`
	}
}

func expectedStr(t *testing.T, a []interface{}, b []interface{}, c string) {
	var flag = reflect.DeepEqual(a, b)
	if flag == false {
		t.Error(fmt.Sprintf("Expected %v, found %v  Error type:   %v", b, a, c)) // t.Error is visible when running `go test -verbose`
	}
}

func expectArr(t *testing.T, a [6]int, b [6]int, c string) {
	if a != b {
		t.Error(fmt.Sprintf("Expected %v, found %v  Error type:  %v", b, a, c)) // t.Error is visible when running `go test -verbose`
	}
}

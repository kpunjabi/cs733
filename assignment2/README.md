

Problem Statement :-

We have to develop a Raft StateMachine in this assignment and then check the behaviour of the Raft statemachine on different input events.Raftstate machine will accept the input events from other Raft nodes and also from the file system layer and depending upon the state(Candidate/Follower/Leader) it will act accordingly.We have to test the Machine by giving different events and observing the set of actions returned by the Statemachine.

Code Description:-

There are two files-> (1) Machine_test.go (2) Raftmachine.go

In the Raftmachine.go the code of the Raftmachine is written and in the Machine_test.go the code to test the statemachine is present.

Raftmachine.go description:-

((For every event(input/output) there is a struct that will contain the parameters of the functions defined for the corresponding event))

Input Events:- (1) AppendEntriesRequestEv:- Send by the Leader to the followers to replicate log entries.And also used as a heartbeat.

(2) AppendEntriesResponseEv:- Used to send the response of AppendEntriesRequestEv.

(3) VoteRequestEv:- Used to send the vote request when the machine is in candidate state.

(4) VoteResponseEv:- Used to send the response to vote request send by other candidate.

(5) AppendEv :- Send by the client to the statemachine and if the statemachine is a leader then it will append the data to its log and then send the append entries request to all the other statemachine in order to replicate that log entry.

(6) TimeoutEv:- It has different meaning(if the machine is Follower then on timeout it will become candidate and start election, if the machine is Candidate then on timeout it will again become the candidate and start elections,if the machine is in Leader state then on timeout it will send heartbeat message to all the followers)

Output Events:-

(1) send(peerid,event):- This is used by the statemachine to send any event to other peer.

(2) Commit(index,data,error):- Used by the statemachine when the entry is successfully committed in the log.

(3) Alarm(t):- Used to send the alarm on time out event.

(4) LogStore(index i,data):- Used by the statemachine to store the new data into the log

Machine_test.go description:-

In this file different functions are present that are used to test the statemachine in different states. For example the function TestFollowerAppendEntryreq will set the machine into the follower state and initialize the variables of the statemachine and then send the AppendEntriesRequestEv in order see the behaviour of the statemachine. The statemachine will return an array of actions that are compared with the expected actions in order to see the behaviour of Statemachine.

How to run the code :- (1) First install Go. (2) Then, put the two files inside the Go directory and then type go test (You will see all the testcases pass message in different inputevents in different states.)

package main


func check(e error) {
    if e != nil {
        panic(e)
    }
}


var sm *StateMachine

type Node interface{

	Append([]byte)

	CommitChannel() <-chan CommitInfo

	CommittedIndex() int

	Get(index int) (err, []byte)

	Id()

	LeaderId() int

	Shutdown()

}

type  CommitInfo struct{
	Data []byte
  	Index int64 // or int .. whatever you have in your code
  	Err error  // Err can be errred
}


type  Config struct{
	cluster []NetConfig  
	Id int 
	LogDir string 

	ElectionTimeout int
	HeartbeatTimeout int
}


type  NetConfig struct{
	Id int
	Host string
	Port int
}



type RaftNode struct { // implements Node interface
		eventCh chan Event
		timeoutCh <-chan Time
}


func (rn *RaftNode) Append(data) {
		rn.eventCh <- Append{data: data}
}

func (rn *RaftNode) CommittedIndex() int{


}

func (rn *RaftNode) Get(index int) (err, []byte){

}

func (rn *RaftNode) Id() int{
   return rn.Id
}

func (rn *RaftNode) LeaderId() int{
   
}

func (rn *RaftNode) Shutdown(){


}

func (rn *RaftNode) processEvents {
for { 
	var ev Event
	select {
	case ev = <- eventCh {ev = msg}
	<- timeoutCh {ev = Timeout{}}
	}
	actions := sm.ProcessEvent(ev)
    doActions(actions)
    }
}
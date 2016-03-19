												Assignment-3
											RAFT Node Implementation
Name : Kuldeep Punjabi
Roll no : 153050006

In this assignment, I have implemented RAFT Node, and integrated it together with the State Machine that I had previously designed.
The implemenation consists of following steps:
- Implemented the Node Interface by defining a new Struct called RaftNode that implements all its functions as follows:
    1. Append([]byte)
	2. CommitChannel() <-chan CommitInfo
	3. CommittedIndex() int
	4. Get(index int) (error, []byte)
	5. Id() int
	6. LeaderId() int
	7. Shutdown()
- For the persistent storage of votedFor and term I have made 5 files corresponding to each node
- First I made a makerafts function which initializes an array of Nodes also start ther whole processing including state machine start etc- Then I called timeout on a particular node (numbered 3) so that it will start election and send vote request to all , then after getting response true from majority of servers get elected as the leader
- after getting the leader in the get leader function 
- I have tested the functionalities like Append(), Get(), leaderId(), CommitedIndex(), ShutDown()

How to execute my program?
- Just call for 
     go test 
   in this folder

Although I have run various  testcases for checking the implemented functionality but failed to implement the mock cluster and check for partition test 
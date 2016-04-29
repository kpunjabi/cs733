												Assignment-4
											(RAFT Implementation)
Name : Kuldeep Punjabi
Roll no : 153050006

In this assignment I have implemented RAFT which is a consensus protocol to manage log replication and maintaining linearaziblity.
The whole implemenataion of this assignment can be recapitulated as follows:

The main code consists of following files :

>>> serve.go : It is creating a file server together with the raft node on  which that file server resides, to configure it we need to provide arguements from the rpc_test.go which are 1) id for the RAFT node 2) port on which this file server runs and it also generates a client handler and a channel corresponding to that client on every file request that arrives , this client handler then accepts requests from the client and the client is also provided a client-id then it puts the clientid to channel mapping in the FSMap maintained in the Fileserver mapper i.e. fs.go .There is a serve thread which is generated as a client handler  which on accepting request calls raftnode.append() also on recieving the response back from the FSmap it gives back the reply to the client. Also it adds the id of client to each message that it send to raft machine and also performs marshalling of message before sending

>>> fs.go :implements the file system part of  the whole setup, contains a FS struct that consists of map for fileinfo and a FSMap which contains maping for client id to channel. It runs in a go routine continuosly monitoring the commit channel of raft node upon recieving a response on commit channel takes out client id from the message then search for the corresponding client channel in the FSmap and puts the response on the channel. It handles the whole filesystem read, write cas and delete functionality.

>>> RaftNode.go : this is the raft node code which handles the request arriving from  the Filesytem also it communicates with the other raft node in the cluster to perform leader based log replication. This was covered as a part of our assignment03 , this time didn't require to make any changes to it. Client contacts first directs its message to any file server(this I have made to port 8090) in the cluster, at the raft node level we find who is the current leader and sends an errror redirect message together with the leaderid if it is not the current leader otherwise appends the entry in the log and performs the necessary replication task.

>> Statemachine.go : This the statemachine code covered as a part of our assignment2 which stores our state machine state machine and contains the code for whole log replication mechanisms, all the message transfers between  raftnode and also code the leader election

Testing of Code:

>>rpc_test.go : It is the test code it generates 5 fileserver first together with their raftnode It passes as argumenst the corresponding cinfigurations then it performs many tests for all situations of file operations. It also defines a message structure that consists of 
	Kind     byte
	Filename string
	Contents []byte
	Numbytes int
	Exptime  int // expiry time in seconds
	Version  int
which is the basic format for all the message transfers that takesplace
In this I have started five servers using os.StartProcess and then run all the testcases and finally to cleanup I have killed these.

>>  First you need to perform go install in the folder Assignmnet1Combined then To run this code just run go test in the folder Assignment1Combined inside assignment4





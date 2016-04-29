												Assignment-4
											(First Check Point)
Name : Kuldeep Punjabi
Roll no : 153050006

To do: 
We need to combine our assignment-3 in which we implemented raft node (integrated with the corresponding state machine) with assignment-1 which is our file system 


What is done uptill now :
As we need come up with a whole new integrated raft protocol implementation that's why for this check point, I tried correcting my code for the assignment-3 so that it becomes suitable to integrate with assignment-1.

I have done many changes to assignment-3 :

-- Earlier the concept of Heartbeat Timeout and Election timeout was not implemented properly so I tried correcting them
-- Also the Leader creation has now turned into random as any of the server acn become the leader now in start earlier I had a generated 	a random no to make a particular node timeout now the concept of leader election has turned into what was desired
-- Also there were were a lot of mistakes in the functions like AppendEntries response for leader, append of leader so corrected them
-- Also earlier the log in the state machine was an array of structure but now I have pointed it out to actual log in raftnode


How to execute my program?
- Just call for 
     go test 
   in this folder

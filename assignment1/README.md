
Name    : Kuldeep Punjabi
Roll no : 153050006

Approach that I followed:

* Data Structure that I used to store information related to files:
  I created a Structure named fileinfo consisting of 5 fields as follows:
   - filename string 
   - version  int64
   - numbytes int
   - exptime  int64
   - createtime int64 - it signifies the time at which file is created or last updated; it is in seconds; its value is assigned using the command time.Now().Unix(); 
   
* I implemented four functions for the file server:
   - write
   - read 
   - cas (compare and swap)
   - delete

* I also designed the file server to be dynamic in the sense that once a file's lifetime has been expired, it get deleted automatically.
  For this, I executed a go routine, that compares the current time and the sum of file's creation time + its expiry time, if the current time is more, it deletes all such file. If the expiry time of the file is 0 then I donot delete as I consider it to stay permanent or with no expiry. 
     

* If I execute write command for a file there are 2 cases:
   
    Firstly, either the file already exists :
      - if I am passing expiry time, then file's expiry time gets updated to the new expiry time passed, and the createtime variable is also updated to the current time, also the file's content gets updated to the passed ones
      - if I am not passing the expiry time, nothing updated
   
    Or the case when file does not already exist:
      - if I am passing the expiry time, then file's expiry time gets updated to the new expiry time passed, and the createtime variable is also updated to the current time, also the file's content gets updated to the passed ones 
      - if I am not passing the expiry time, expiry time is set to 0 which considers that the file has no expiry and is meant to be permanent
  
      
* If I execute cas command for a file there are 2 cases:
 
    Firstly, either the file already exists :
      - First I check for the version that file's current version is same as the version passed otherwise an error is thrown as "ERR_VERSION"
      - if I am passing expiry time, then file's expiry time gets updated to the new expiry time passed, and the createtime variable is also updated to the current time, also the file's content gets updated to the passed ones, and its version is incremented by 1
      - if I am not passing the expiry time, no updation to expiry time rest other updates are same
      
    Or the case when file does not already exist:
      - Return error as "ERR_FILE_NOT_FOUND"

* I even checked for the format of the command that its length is correct :
       - the length of the write command must be either 3 or 4 depending upon whether expiry time is passed or not
       - the length of the cas command must be either 4 or 5 depending upon whether expiry time is passed or not
       - the length of the read command must be 2
       - the length of the delete command must be 2 
  if there is an error in command format, server passes an error "ERR_CMD_ERR"
  
* I used mutex Lock-Unlock mechanism to provide concurrency in the code wherever there is an access to the structure fileinfo I have put that in the mutex as a critical section, so that it provide consistency in the operations to every clients that run.  

* I have even considered the case when the data itself contains "\r\n" as its part, then I have first split it into two parts based on the first "\r\n" and then used trim to remove the occurence of "\r\n" from the end of the content


* Assumptions considered:
 - I consider the server to be up and running as the structure used to store information about the files stays in the main memory gets deleted after every execution of the server.
 - I have even created two files statically in the structure and stored their info already.
 - Also as I have used struct as the storage structure, I have considered its length to be static as 500 files
 
* Problem encountering :
  - As the structure needs to always be in memory for the server to work properly, so the code will work if the file doesn't exist already, in case if it already exists I am encountering some problem that it is not able to get the file info that existed earlier

* Testing of the server:
 - I have included the test file basic_test.go on which I have tested my file server and it ran successfully PASS ok
 - I have tested it using the command go test -race
 - I tested the server for its concurrency by creating 10 clients, made them to write, read ,cas and delete the same file 
    ~ for the write operation, the file end up with the content written by the last client
    ~ for the read operation, they read the latest version of the content
    ~ for the cas operation, only the first client was able to execute the operation and the other nine end up with an error as ERROR_VERSION
    ~ for the delete operation, only the first client to execute the operation, was able to delete the file, others could not.
 
 
* References that I took help from:
 - regarding tcp socket programming in go from the site : http://loige.co/simple-echo-server-written-in-go-dockerized/
 - for the string functions in go : https://golang.org/pkg/strings
 - for the struct functionality : https://gobyexample.com/structs
 - https://golang.org/src/net/http/fs.go?h=FileServer#L473     
 - http://stackoverflow.com/questions/19208725/example-for-sync-waitgroup-correct
 - https://gobyexample.com/mutexes

package mr

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)


type Coordinator struct {
	// Your definitions here.
	intermediate   []KeyValue
	files          []string
	fileIndex      int
	mu             sync.Mutex
	wg             sync.WaitGroup
	updateMapping  sync.Mutex
	updateReducing sync.Mutex
}

// Your code here -- RPC handlers for the worker to call.
func (c *Coordinator) assignMappingTask(args *AssignTaskArgs, reply *AssignTaskReply) error {
	defer c.mu.Unlock()
	reply.filename = c.files[c.fileIndex]
	file, err := os.Open(reply.filename)
	if err != nil {
		log.Fatalf("cannot open %v", reply.filename)
		return err
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", reply.filename)
		return err
	}
	file.Close()
	reply.content = string(content)
	c.fileIndex++
}

func (c *Coordinator) assignReducingTask(args *AssignTaskArgs, reply *AssignTaskReply) error {
	defer c.mu.Unlock()

}

func (c *Coordinator) UpdateMappingResult(args *UpdateMappingResultArgs, reply *UpdateMappingResultReply) error {
	c.updateMapping.Lock()
	c.intermediate = append(c.intermediate, args.intermediate...)
	c.updateMapping.Unlock()
}

func (c *Coordinator) UpdateReducingResult(args *UpdateReducingResultArgs, reply *UpdateReducingResultReply) error {
	c.updateMapping.Lock()
	c.updateMapping.Unlock()
}

// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
func (c *Coordinator) AssignTask(args *AssignTaskArgs, reply *AssignTaskReply) error {
	c.mu.Lock()
	if c.fileIndex < len(c.files) {
		c.assignMappingTask(args, reply)
	} else {
		c.assignReducingTask(args, reply)
	}
}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.

	return ret
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}
	c.fileIndex = 0

	// Your code here.

	c.server()
	return &c
}

//1st attempt
//files
//worker?
//files?
//reading a file doesnt change, no copy needed send as many ?
//split the map task in x workers and jump in the loop by x for the threads.
//channel for shared memory? we need a thread safe queue for sure.
//we need a view of these files for the x workers, thats also fine.
//is it important? everyone just writes into it? (noone needs to consume it atleast yet.)( just writing
// needing to read is atomic I believe) ( its not) ( need the channel here) ( reducers will consume it )
//when all threads are done? how do we know? when the last element is processed by the thread, call it done
//maybe call a thread  for each file?
//and call it done after
// sort a channel by key?
//same keys go to the same reducer
//send the reducer keys
//first thought is to go through the map key values we got sequentially, to give it to reducers which expect a key
//and a list of values
//I am thinking a channel for each reducer/key is probably better
//second I am thinking sort itself might be nlogn and that too sequential, do I really wanna use sort here?
//at any point I can have as many unique channels as the reducers?
// no that doesnt work why? because we might not get all the keys?
// channels equal to reducers or equal to keys, not sure ( equal to nreducers is my idea down expressed)
//if we have a channel for a unique key maybe, but do we want that?
// we dont need it, n reducers are enough, let them handle multiple keys all at once and maintain the state ( counter for all fo them?)
// but that is no the reducer I got, the reducer I got expect a key and list of all the values
//my reducer may have been far less time complex though, far faster.
//but fine
//no, my algo feels like reducer is not the time consuming/ cpu consuming task, but rather the breaking the map
//to give these reducer values is, which should not be the point, though it should not,
// and my algo, would essentially be a sequential caller for the reduce function
// they limited the nreduce, for this, so that I cant have as many threads I want?
// and I am forced to sort these keys and then send the batches to reduce functions.
// how do we construct this same key and list of values from the map quickly?
// so If I sort it then I can just send it as a view to my reducer which will call these guys reducer
//If I dont sort it then I have to iterate upon it and make a big list of lists and hand over the lists to my reducer
// as it is
///key string of 1s
//map key, array of strings -> key, array of strings of 1
//dont sort it, send it to the workers, let them start multiple reducers themselves.? our threads will be more, but we are te
//technically only using n reducers?
// yeah the rest of the thought is seeming correct, we are even sending the next part of it wiht the threads ( the next reduce
// is with threads and thus completely correct and no sorting needed)
// the only question now is , what does it mean by number of nReduce? ( is it the total number of threads or not?)

//2nd attempt
//now then taking another hit at it
//maincoordinator is passed nreduce and files.
//worker is not passed anything and asks maincoordinator for the files and so on.
//or for the task and it c
//either mapper or reducer later
// we have a channel key value map -> key value
//we have a channel filename : string of words
// word1:1, word2:1 etc. ( just for its own files i think)
//sort this one mapper results and send them?
//no, we needed this to be stored in separate files for worker 1, worker 2 etc based on key.
//so all unique keys go to same worker. ( and then the worker sorts(to group) and sends combined result for the keys of one type)
//
// give him a bunch of file chunks, thats a mapper? how many? no idea, assume equal to nreduce
//divide the files and send, no special thing for the mappers.
//get a shared channel I think for each mapper ( thats like their files locaation?)
// we have to get a range of channels based on the keys ( unsorted sorted doesnt matter for now)
//send all the channels to each mapper. and the size of channels? buffered, Max size of a channel needs to be the total file size* number of files
//no other thread safe way of putting things in? we dont need to keep every operation atomic ( we dont to save in one by one)
// in each map we can maintain separate lists for each unique key? no each unique reducer key it encounters and put them in ,
// then at last put them in the channel. this channel is a complete data to be sent for each reducer.
//we cant declare sizes, but we can say that it will be m total lists?
//each mapper is supposed to put their lists into this, even if its empty. ( thats when we can ensure to create a buffer of appro
//priate size, if we cant make it infinite) ( we cant so, this is the only way and nicely, it doesnt need to finish processing from all mappers?)
//no mappers can only put an empty list when they have done a full processing on their own data.
//one thinks there are nreduce mappers or some number mappers but we dont know that( I am not even calling the worker) ( maybe
//there will be just one worker, so we can only call these tasks and think in that term, not mappers or reducers or workers or
//machines ( its agnostic to that))
//nreduce mappers * unique keys, that will be the buffer size, unique keys will be number of lists
// but I dont know the number of keys before hand ( I am not allowed to count in coordinator)
//yeah, channels dont work in this, not like this(can be made to work if I process the data to find unique
// keys in coordinator, but thats just hacky nonsense), this is state management and synchronization, rather than communication with
//each other
//mutex or files are my only options
//the threads dont really communicate with each other for map or reduce
//first all individual mappers need to be done and then all the reducers can go individually on their own data,
//not really communication, rather synchronization.
//use mutex locks or files
//how do I apply mutex lock? ( cant share data with mutex lock, but make sure that I can use a regular object and put
// rather than it being a channel)
//I can make sure that a regular object becomes thread safe ( atomic visible etc)
//once all workers done, we can send it
//.done(), add and wait on sync group can also be a pain here, because they said that some worker might not
//perform, so we add one more number, no we dont, we dont add for any extra workers and wait for done(), only either
// one of them can call done and the thread will go ahead, anyone calling done() after should have no effect ( negative done
// should have no effect)?
//and then I process this all set of keys? and call the reducer in coordinator?
//maybe not, thats just a bottleneck, probably what I need to do
//have an unordered map and key and list of values in it, so each key is unique and we can just have it by the key.
//this way main wont process anything? generally a reducer is supposed to group the data, but right now I am grouping it and giving
//him 4( or however many keys he gets) separate lists to process?
//looks good to me.

//3rd attempt
//coords have the files and nreduce
//workers have the functions map and reduce
//the worker will ask for a task from coordinator and will receive one
//either it receive the map and go, or reducer and wait for maps to finish first ,
//the communication cant happen throgh anything else other than files and rpc.
//the shared memory is only accessed by threads within a process.
//so we have files and message queues and other such things and I think the RPC too
//use rpc I guess, or use files as in the real systems do, for both mappers and reducers maybe.
//mappers will return the files or file descriptors or what??
//reducers dont need to return anything they just output their files.
//the mappers can just return the name of the file too, if that lets reducer open up later.
// there is no need to decide on number of mappers, let workers ask, and give each worker a file,
//whoever gets the job done and shares with me the filename back is good,
//when all files are done we are done.
//nreduce will be the number of files generated finally. ( but cant be more than the numbeer of unique keys)
// a worker atleast gets one key.
// The multithreading aspect of the question comes from the coordinator running a hundred thread for each of the connection
//request it gets from the workers which is the main part of handling concurrency here.
//

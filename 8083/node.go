package main
import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"net"
	"time"
	"strconv"
	"math/rand"
)

type Node struct {
	ip   string
	port string
	receive bool 
	send bool
	size int
}
var sendBackFlag bool
var iniFlag bool
var sendIniMessFlag bool
var Iam Node
var Parent Node
var neighbors []Node
var leader string
var status bool
var networkSize int
var myRandId int
var roundNumber int

func main() {
	filename := `configuration.conf`
	fmt.Println("Start..."  )
	// Find Iam,Initiator,Neighbors
	readFile(filename)
	fmt.Printf("I am %s:%s \niniFlag is: %t \nAll my neighbors are: %v \n" , Iam.ip , Iam.port , iniFlag , neighbors)
	go server(Iam)

	if checkNeighborServer(neighbors) {
		done := false
		for {
			if iniFlag && !sendIniMessFlag {
				fmt.Println("Start to send message from initiator: "  )
				
				roundNumber = 0
				
				leader := strconv.Itoa(myRandId) 
				rn := strconv.Itoa(roundNumber) 
				iniMss := "&Iam="+Iam.ip+":"+Iam.port+"&id="+leader+"&round="+rn+"&size=0"+ "&back=false"
				sendMssToAllNeighbors(iniMss)
				sendIniMessFlag = true
				fmt.Printf("All my neighbors are: %v \n" , neighbors)
			}else{
				time.Sleep(3000 * time.Millisecond)
				strSize , intSize := findSize()
				if status {
					if intSize == networkSize {
						fmt.Println("\nDone "  )
						done = true	
					}else{
						
						if checkReceiveFromAll(){	
							roundNumber = roundNumber+1
							rn := strconv.Itoa(roundNumber) 
							myRandId = selectRandomId()
							leader := strconv.Itoa(myRandId) 
							fmt.Println(" NEW ID : "  , leader  )
							iniMss := "&Iam="+Iam.ip+":"+Iam.port+"&id="+leader+"&round="+rn+"&size=0"+ "&back=false"
							deleteAllActivities()
							sendMssToAllNeighbors(iniMss)
						}
					}
				}else{
					if checkReceiveFromAll(){
						m := "&Iam="+Iam.ip+":"+Iam.port+"&id="+leader+"&round="+strconv.Itoa(roundNumber)+"&size=" + strSize+ "&back=true"
						sendMessage(m,Parent)
						fmt.Printf( "Neighbors: %v \nNetwork size is: %v\nI find %v\n" , neighbors,networkSize,intSize )
						deleteAllActivities()
						Parent.port = "0"
						//done = true	
					}
				}
				
			}
			if done {
				_ , intSize := findSize()
				fmt.Printf( "Neighbors: %v \nNetwork size is: %v\nI find %v\n" , neighbors,networkSize,intSize )
				fmt.Println( "I am the Leader"  )
				break
			}
		}
	}
	
	
}

func selectRandomId() int {
	rand.Seed(time.Now().Unix() * 1234)
	return rand.Intn(networkSize-1)+1
}

func readFile(fileName string){
	f, _ := os.Open(fileName)
	defer f.Close()
	r := bufio.NewReaderSize(f, 2*1024)
	line, isPrefix, err := r.ReadLine()
	i := 1
	for err == nil && !isPrefix {
		s := string(line)
		if i == 1 {
				// Find Iam
				t :=strings.Split(s, ":")
				Iam = Node{t[0],t[1],false,false,0}
		}else{
			k :=strings.Split(s, ":")
			if k[0] == "initiator" {
				// Find if initiator
				iniFlag = true
				//ser active
				status = true
			}else{
				if k[0] == "size" {
				networkSize,_ = strconv.Atoi(k[1]) 
				fmt.Println("Network Size is: " , networkSize)
				}else{
					// Find neighbors
					neighbors = append(neighbors, Node{k[0],k[1],false,false,0})
					//set passive
					status = false					
				}
				
			}
			
		}
		i++
		line, isPrefix, err = r.ReadLine()		
	}
	
	//find rand number
	myRandId = selectRandomId()
	fmt.Println("My ID: ", myRandId  )

}

func analizMessage(message string) map[string]string{
	ms :=strings.Split(message, "&")
	GetMessage := make(map[string]string)
	
	msIam :=strings.Split(ms[1], "=")
	mx :=strings.Split(msIam[1], ":")
	GetMessage["ip"] =  mx[0]
	GetMessage["port"] =  mx[1]
	
	getId :=strings.Split(ms[2], "=")
	GetMessage["id"] =  getId[1]
	
	getRound :=strings.Split(ms[3], "=")
	GetMessage["round"] =  getRound[1]
		
	getSize :=strings.Split(ms[4], "=")
	GetMessage["size"] =  getSize[1]
		
	getBack :=strings.Split(ms[5], "=")
	GetMessage["back"] =  getBack[1]
		
	fmt.Println("GetMessage-> ",GetMessage)
	return GetMessage
}

func server(s Node) {
	fmt.Printf("Launching server... %s:%s \n" , s.ip,s.port)
	ln, _ := net.Listen("tcp", s.ip+":"+s.port)
	defer ln.Close()
	for {
		conn, _ := ln.Accept()
		message, _ := bufio.NewReader(conn).ReadString('\n') 
		if string(message) != "" {
			fmt.Println("->", string(message))
			doIt(analizMessage(message))
		}
	}	
	
}

func checkNeighborServer(n []Node) bool{
	for i:=0; i < len(n);i++{
		for {
			conn, err := net.Dial("tcp", n[i].ip+":"+n[i].port)
			fmt.Println("Looking for " + n[i].ip+":"+n[i].port)
			time.Sleep(3000 * time.Millisecond)
			if err == nil {
				conn.Close()
				break
			}
		}
	}
	
	return true
}

func sendMessage(s string, n Node){
	conn, _ := net.Dial("tcp", n.ip+":"+n.port)
	defer conn.Close()
	conn.Write([]byte(s))
	fmt.Printf("Message Sent to %s:%s \n" ,n.ip,n.port )	
}

func sendMssToAllNeighbors(ms string){

	for i:=0; i < len(neighbors);i++{
		if Parent.port != neighbors[i].port {
			sendMessage(ms,neighbors[i])
			neighbors[i].send = true	
		}
		
	}
}

func findSize() (string,int){
	size := 0
	for i:=0; i < len(neighbors);i++{
		if Parent.port != neighbors[i].port {
			size = size + neighbors[i].size
		}
	}
	return strconv.Itoa(size+1), size+1
}

func findNewLeadership( l string) bool{
	myReturn:=false
	if status {
		lastLeader,_ := strconv.Atoi(leader)
		newLeader,_ := strconv.Atoi(l)
		if  lastLeader < newLeader {
			fmt.Println("I'm going to change my leader to: ", l )
			leader = l
			myReturn = true
		}
	}else{
		// First leader
		leader = l
	}
	return myReturn
}
func deleteAllActivities(){
	fmt.Println("I'm going to delete all previous activities: " )
	for i:=0; i < len(neighbors);i++{
		fmt.Println("*** Delete all activities about " +neighbors[i].ip + ":" + neighbors[i].port )
		neighbors[i].receive = false
		neighbors[i].send = false
		neighbors[i].size = 0
	}
	iniFlag = false
}

func doIt( ms map[string]string){
	
	fmt.Println(" ++++  ",  ms , leader)
	_,id := findNodeBtwNeighbors(ms["ip"],ms["port"])
	

	if ms["back"] == "true" {
		fmt.Println(" ************** Back Message ",  ms)
		neighbors[id].receive = true
		neighbors[id].size , _ = strconv.Atoi(ms["size"])
	}else{
		thisId , _ := strconv.Atoi(ms["id"])
		round , _ := strconv.Atoi(ms["round"])
		if roundNumber < round || ( roundNumber == round && myRandId < thisId ){
			//select it as parent
			fmt.Println(" ************** select it as parent ",  ms)
			status = false
			//delete all activities
			deleteAllActivities()
			
			fmt.Println("I'm going to change my parent to: ",  ms["ip"],ms["port"])
			fmt.Println("_____________ ",  Parent)
			Parent.ip =  ms["ip"]
			Parent.port =  ms["port"]
			neighbors[id].receive = true
			
			leader = strconv.Itoa(thisId) 
			roundNumber = round
			myRandId = thisId
			
			sendMss := "&Iam="+Iam.ip+":"+Iam.port+"&id="+leader+"&round="+ms["round"]+"&size=0"+ "&back=false"
			neighbors[id].receive = true
			sendMssToAllNeighbors(sendMss)
		}
		
		if roundNumber == round && myRandId == thisId {
			//continue the echo
			fmt.Println(" ************** continue the echo ",  ms)
			neighbors[id].receive = true
			neighbors[id].size , _ = strconv.Atoi(ms["size"])
		}
		
		if roundNumber > round || ( roundNumber == round && myRandId > thisId ){
			//ignire the message
			fmt.Println(" ************** ignire the message ",  ms)
		}	
	}
		
}

func findNodeBtwNeighbors(ip string, port string) (Node , int){
	j := 0
	for i:=0; i < len(neighbors);i++{
		if neighbors[i].ip == ip && neighbors[i].port == port {
				j = i
		}
	}
	return neighbors[j],j
}

func checkReceiveFromAll() bool{
	Myreturn := true
	for _, n := range neighbors {
		if Parent.port != n.port   {
			if !n.receive {
				Myreturn = false
			}
		}
	}
	return Myreturn
}

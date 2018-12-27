package main

import (
	"fmt"
	"log"
	"net"
	"io"
	"os"
	"time"
	"sync"
	"flag"
	"strings"

	// "github.com/golang/protobuf/proto"
	pb "github.com/wangxi19/gotraining/grpc/proto"
	"github.com/wangxi19/utils/sqlutil"
	"google.golang.org/grpc"
)


type mutexBool struct {
	v bool
	m sync.Mutex
}

func (mb *mutexBool) set(val bool) {
	mb.m.Lock()
	defer mb.m.Unlock()
	
	mb.v = val
}

func (mb *mutexBool) get() bool {
	mb.m.Lock()
	defer mb.m.Unlock()
	
	return mb.v
}

type userServer struct {
}

func (s *userServer) GetUserList(stream pb.User_GetUserListServer) error {
	var rsterr error
	timer := time.NewTimer(5000 * time.Millisecond)
	go func() {
		for {
			searchWheres, err := stream.Recv()
			timer.Reset(5000 * time.Millisecond)

			if io.EOF == err {
				rsterr = nil
				return
			}
			
			if nil != err {
				rsterr = err
				return
			}
		
			sqlwhere := &[]string{}

			for _, oneWhere := range searchWheres.Wheres {
				*sqlwhere = append(*sqlwhere, oneWhere.Key + " in ('" + strings.Join(oneWhere.Val, "', '") + "')")
			}

			db, err := dbpool.GetDB(dbname)
			if nil != err {
				rsterr = err
				return
			}

			arrayMap, err := sqlutil.SelectArrayMap(db, "sys_user", "*", "(" + strings.Join(*sqlwhere, ") OR (") + ")", "", -1, -1)

			if nil != err {
				rsterr = err
				return
			}

			var userMapList []*pb.UserList_UserMap
			for _, oneMap := range arrayMap {
				rowmap := map[string]string{}
				for k, v := range oneMap {
					rowmap[k] = string(v[:])
				}
				userMapList = append(userMapList, &pb.UserList_UserMap{RowMap: rowmap})
			}
		
			if err = stream.Send(&pb.UserList{Usermap: userMapList}); nil != err {
				rsterr = err
				return
			}
		}
	}()

	for {
		select {
		case <-timer.C:
				if nil != rsterr {
					fmt.Fprintf(os.Stdout, "[ERROR]: %v\n", rsterr)
				}

				fmt.Fprint(os.Stdout, "timeout, connection will be closed\n")
				return rsterr
		}
	}
}

var (
	dbpool sqlutil.DBPool

	username string
	password string
	host string
	port string
	dbname string
)

func parse() {
	flag.StringVar(&username, "username", "", "db username")
	flag.StringVar(&password, "password", "", "db password")
	flag.StringVar(&host, "host", "", "db server address")
	flag.StringVar(&port, "port", "", "db port")
	flag.StringVar(&dbname, "dbname", "", "db name")
	// flag.Usage = func () {
	// 	flag.Usage()
	// 	fmt.Println("\n\n hello, world")
	// }

	flag.Parse()
}

func main() {
	defer func () {
		if err := recover(); nil != err {
			log.Fatalf("a error occured, %v", err)
		}
	}()
	
	parse()

	err := dbpool.InitDB("postgres", username, password, host, port, dbname, 15)
	if nil != err {
		log.Fatalf("dbpool: %v\n", err)
	}

	lis, err := net.Listen("tcp", "localhost:10000")
	if nil != err {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServer(grpcServer, &userServer{})

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

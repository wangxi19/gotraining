package main

import (
	"log"
	"context"
	// "time"
	"os"
	"os/signal"
	"fmt"
	"flag"
	"reflect"
	"strings"

	"github.com/wangxi19/utils/fileutil"

	pb "github.com/wangxi19/junk/proto"
	"google.golang.org/grpc"
)

var (
	message string
)

func parse() {
	flag.StringVar(&message, "msg", "1", "A message whose be used to sending")
	flag.Parse()
}

func main () {
	defer func() {
		if err := recover(); nil != err {
			log.Fatalf("a error occured, %v", err)
		}
	}()

	parse();

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)

	conn, err := grpc.Dial("localhost:10000", grpc.WithInsecure())
	if nil != err {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	// ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	// defer cancel()

	c := pb.NewUserClient(conn)

	stream, err := c.GetUserList(context.Background())
	if nil != err {
		log.Fatalf("GetUserList: %v", err)
	}

	err = stream.Send(&pb.SearchKey{Key: &pb.SearchKey_Id{Id: message}})
	if nil != err {
		log.Fatalf("Send error: %v", err)
	}

	userLst, err := stream.Recv()
	if nil != err {
		log.Fatalf("Recv error: %v", err)
	}

	idx := 0
	content := ""
	strKeys := []string{}

	for _, oneMap := range userLst.Usermap {
		if idx == 0 {
			keys := reflect.ValueOf(oneMap.RowMap).MapKeys()
			for _, k := range keys {
				strKeys = append(strKeys, k.String())
			}
			fmt.Fprintf(os.Stdout, "%s\n", strings.Join(strKeys, "\t"))
			content += strings.Join(strKeys, "\t") + "\r\n"
			idx++
		}

		vals := []string{}
		for _, k := range strKeys {
			vals = append(vals, oneMap.RowMap[k])
		}
		fmt.Fprintf(os.Stdout, "%s\n", strings.Join(vals, "\t"))
		content += strings.Join(vals, "\t") + "\r\n"
	}
	fileutil.WriteFileString("/root/sql.txt", content, false)
	
	go func() {
		select {
		case <-chSignal:
			os.Exit(0)
		}
	}()
}
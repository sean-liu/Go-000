package main

import (
	"SeanLiu_Go-000/Week09/internal"
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {

	listener, err := net.Listen("tcp", ":10088")
	if err != nil {
		panic(err)
	}

	rm := internal.NewDefaultRoutineManager()

	messageChan := make(chan func() (string, net.Conn), 1)
	rm.Register(func() bool {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			return true
		}

		input := bufio.NewScanner(conn)
		for input.Scan() {
			fmt.Println("received :" + input.Text())
			messageChan <- func() (string, net.Conn) { return input.Text(), conn }
		}

		return false
	})

	rm.Register(func() bool {
		select {
		case funcValue := <-messageChan:
			content, conn := funcValue()
			fmt.Fprint(conn, content)
		default:
			break
		}
		return false
	})

	done, cancel := rm.Start()

	rm.WaitForSystemSignal(func() {
		fmt.Println("signal received, about to cancel")
		listener.Close()
		cancel()
		<-done
		fmt.Println("all subtasks finished")
	})
}

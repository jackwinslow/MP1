package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type process struct {
	id   string
	ip   string
	port string
}

func main() {

	process_num := os.Args[1]
	delays := make([]int, 2)

	file, err := os.Open("./config.txt")
	if err != nil {
		fmt.Println("Failed to open config file!")
		return
	}

	fs := bufio.NewScanner(file)

	// Map to store processes for easy access
	processes := make(map[string]process)

	// Parses delays from config
	fs.Scan()
	unprocessed_delays := strings.Fields(fs.Text())
	delays[0], _ = strconv.Atoi(unprocessed_delays[0])
	delays[1], _ = strconv.Atoi(unprocessed_delays[1])

	// Adds each process to a slice for easy access
	for fs.Scan() {
		fields := strings.Fields(fs.Text())
		p := process{id: fields[0], ip: fields[1], port: fields[2]}
		processes[fields[0]] = p
	}

	// fmt.Print(processes)

	// Initialize receiver
	address := processes[process_num].ip + ":" + processes[process_num].port
	l, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer l.Close()

	// Run receiver listening for new connections
	go unicast_receive(l)

	// Run loop to send messages on user input
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		fields := strings.Fields(text)
		message := ""

		// Checks to see if user included 'send' in input
		if fields[0] != "send" {
			continue
		}
		// Checks to see if inputted destination process exists in config
		if _, ok := processes[fields[1]]; !ok {
			continue
		}
		// Reconstruct message from fields splice
		for i := 2; i < len(fields); i++ {
			message = message + fields[i] + " "
		}
		message = strings.TrimSpace(message)

		go func() {
			t := time.Now()
			formattedTime := t.Format(time.RFC3339) + "\n"

			fmt.Printf("Sent \"%s\" to process %s, system time is %s", message, fields[1], formattedTime)

			rand.Seed(time.Now().UnixNano())
			delay := rand.Intn(delays[1]-delays[0]) + delays[0]
			time.Sleep(time.Duration(delay) * time.Millisecond)

			unicast_send(processes[fields[1]], process_num+" "+message)
		}()

	}
}

func unicast_send(destination process, message string) {
	address := destination.ip + ":" + destination.port
	c, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Fprintf(c, message+"\n")
}

func unicast_receive(source net.Listener) {
	for {
		c, err := source.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go func(c net.Conn) {
			netData, err := bufio.NewReader(c).ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}

			t := time.Now()
			formattedTime := t.Format(time.RFC3339) + "\n"

			fields := strings.Fields(netData)
			source_num := fields[0]
			message := ""
			// Reconstruct message from fields splice
			for i := 1; i < len(fields); i++ {
				message = message + fields[i] + " "
			}
			message = strings.TrimSpace(message)

			fmt.Printf("Received \"%s\" from process %s, system time is %s", message, source_num, formattedTime)
		}(c)
	}
}

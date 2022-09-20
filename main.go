package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type process struct {
	id   int
	ip   string
	port int
}

func main() {

	file, err := os.Open("./config.txt")
	if err != nil {
		fmt.Println("Failed to open config file!")
		return
	}

	s := bufio.NewScanner(file)

	// Slice to store process objects created
	// during config reading
	processes := make([]process, 0)

	// Adds each process to a slice for easy access
	for s.Scan() {
		fields := strings.Fields(s.Text())
		id, _ := strconv.Atoi(fields[0])
		port, _ := strconv.Atoi(fields[2])
		p := process{id: id, ip: fields[1], port: port}
		processes = append(processes, p)
	}

	fmt.Print(processes)
}

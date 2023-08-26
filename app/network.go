package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
)

func Open_port(port int) {
	println("opening port ", port)

	// servAddr := "0.0.0.0:" + port
	// Start the server and listen for incoming connections
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	fmt.Println("Server started, listening on port", port)

	for {
		conn, err := listener.Accept() // Accept incoming connection
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle incoming connection in a new goroutine
		go handleConnection(conn)
	}
}
func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Accepted connection from", conn.RemoteAddr())
	buf := make([]byte, 1024)

	// Read data from the connection
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	receivedMessage := string(buf[:n])
	fmt.Println("Received:", receivedMessage)

	// Prepare and send a reply message
	replyMessage := "Server: Thanks for the message!"
	_, err = conn.Write([]byte(replyMessage))
	if err != nil {
		fmt.Println("Error writing:", err)
	}
	fmt.Println("Sent reply:", replyMessage)
}

func List_of_user(port int) {
	// go through the subnet and list all the users
	//users -> who opened the port
	// do this in each sec
	//global user map strtuct

	//find IPv4
	startIP := getSubnetStartIP()
	endIP := "255"
	execludeIp := getLocalIP()
	// Construct the nmap command with the IP range and port
	cmd := exec.Command("nmap", "--open", "-p", fmt.Sprintf("%d", port), fmt.Sprintf("%s-%s", startIP, endIP), "--exclude", execludeIp)

	// Capture the command's output
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Execute the command
	fmt.Println("Scanning IP", startIP, "-", endIP, "By execluding ", execludeIp)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("stderr:", stderr.String())
		return
	}

	ipAddresses := extractIPAddresses(stdout.String())

	//clear ip_table and allocate again
	Ip_table = Ip_table[:0]

	if len(ipAddresses) == 0 {
		fmt.Println("no clients found")
	} else {
		for _, ip := range ipAddresses {
			Ip_table = addIfNotExists(Ip_table, ip)
		}

		println("total users ")
		for _, eachIp := range Ip_table {
			fmt.Println(eachIp)
		}
	}

}

func getSubnetStartIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error:", err)
		return err.Error()
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			// Skip loopback and non-IPv4 addresses
			if ipNet.IP.IsLoopback() || ipNet.IP.To4() == nil {
				continue
			}

			ones, _ := ipNet.Mask.Size()
			if ones > 0 {
				// Calculate the subnet start IP address
				startIP := ipNet.IP.Mask(ipNet.Mask)

				fmt.Printf("Interface: %s, Subnet Start IP: %s\n", iface.Name, startIP.String())
				return startIP.String()
			}
		}
	}
	return "not found"
}

func extractIPAddresses(scanOutput string) []string {
	ipPattern := `\b(?:\d{1,3}\.){3}\d{1,3}\b`
	re := regexp.MustCompile(ipPattern)
	matches := re.FindAllString(scanOutput, -1)

	return matches
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Get the local address from the connection
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	// Convert the net.IP to a string
	ipString := localAddr.IP.String()

	return ipString
}

func addIfNotExists(slice []string, value string) []string {
	for _, v := range slice {
		if v == value {
			return slice // Value already exists, so return the same slice
		}
	}
	return append(slice, value) // Value doesn't exist, so add it to the slice
}

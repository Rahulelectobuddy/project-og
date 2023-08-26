package main

import (
	"fmt"
	"sync"
	"time"
)

// global IP table to store user
var Ip_table []string
var Port int = 13370

func taskClientDetails() {
	for {
		//create a list of peers with their name
		List_of_user(Port)

		connectToClient() //connect to each clients

		time.Sleep(time.Duration(5) * time.Second) // Sleep for 1 second before the next iteration
	}
}

func main() {
	userName := ""

	//users
	fmt.Print("Total users ", Ip_table)

	//if no userName present ask for that
	if userName == "" {
		fmt.Println("No name please set \nEnter your name :")
		fmt.Scanln(&userName)
	}
	fmt.Println("Your username is", userName)
	// now set userName in header

	// now create a list of users who are available in this subnet.
	// user -> who opened the port *13370*
	// first I should be a user and open the port.

	// Create a WaitGroup to wait for the goroutine to finish
	var wg sync.WaitGroup

	// Add 1 to the WaitGroup for the goroutine
	wg.Add(1)

	// Run Open_port in a separate goroutine
	go func() {
		defer wg.Done() // Mark the goroutine as done when it finishes
		Open_port(Port)
	}()

	// Start the continuous task in a goroutine
	go taskClientDetails()

	// Wait for the goroutine to finish before exiting
	wg.Wait()

}

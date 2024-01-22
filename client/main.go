package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/priyanshu/trainservice-database/proto"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Establish a connection to the server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	client := proto.NewTrainServiceClient(conn)

	// Initialize a scanner for reading user input
	scanner := bufio.NewScanner(os.Stdin)

	// Main loop for user interaction
	for {
		fmt.Print("Commands: add, purchase, view, remove, modify, receipt, exit\n")
		fmt.Print("Enter command: ")

		// Read user command
		if !scanner.Scan() {
			break
		}
		command := strings.TrimSpace(scanner.Text())

		// Process user command
		switch command {
		case "add":
			addUser(client, scanner)
		case "purchase":
			purchaseTicket(client, scanner)
		case "view":
			viewUsersBySection(client, scanner)
		case "remove":
			removeUser(client, scanner)
		case "modify":
			modifySeat(client, scanner)
		case "receipt":
			getReceiptForUser(client, scanner)
		case "exit":
			fmt.Println("Exiting the client.")
			return
		default:
			fmt.Println("Invalid command. Please enter a valid command.")
		}
	}
}

// addUser interacts with the user to add a new user
func addUser(client proto.TrainServiceClient, scanner *bufio.Scanner) {
	user := &proto.User{}

	// Get user details
	fmt.Print("Enter first name: ")
	if !scanner.Scan() {
		return
	}
	user.FirstName = scanner.Text()

	fmt.Print("Enter last name: ")
	if !scanner.Scan() {
		return
	}
	user.LastName = scanner.Text()

	fmt.Print("Enter email: ")
	if !scanner.Scan() {
		return
	}
	user.Email = scanner.Text()

	fmt.Print("Enter new password (Minimum-8 characters): ")
	if !scanner.Scan() {
		return
	}
	user.UserId = scanner.Text()

	// Call gRPC method to add user
	response, err := client.AddUser(context.Background(), user)
	if err != nil {
		log.Fatalf("Error adding user: %v", err)
	}

	// Print or process the response accordingly
	fmt.Printf("Response: %+v\n", response)
}

// purchaseTicket interacts with the user to purchase a ticket
func purchaseTicket(client proto.TrainServiceClient, scanner *bufio.Scanner) {
	ticket := &proto.Ticket{
		User: &proto.User{}, // Ensure a User object is created within Ticket
	}

	// Get ticket details
	fmt.Print("Enter password: ")
	if !scanner.Scan() {
		return
	}
	ticket.User.UserId = scanner.Text()

	fmt.Print("Enter from: ")
	if !scanner.Scan() {
		return
	}
	ticket.From = scanner.Text()

	fmt.Print("Enter to: ")
	if !scanner.Scan() {
		return
	}
	ticket.To = scanner.Text()

	// Call gRPC method to purchase ticket
	request := &proto.PurchaseRequest{Ticket: ticket}
	response, err := client.Purchase(context.Background(), request)
	if err != nil {
		log.Fatalf("Error purchasing ticket: %v", err)
	}

	// Print or process the response accordingly
	fmt.Printf("Response: %+v\n", response)
}

// viewUsersBySection interacts with the user to view users in a section
func viewUsersBySection(client proto.TrainServiceClient, scanner *bufio.Scanner) {
	fmt.Print("Enter section: ")
	if !scanner.Scan() {
		return
	}
	section := scanner.Text()

	// Call gRPC method to view users by section
	request := &proto.ViewUsersBySectionRequest{Section: section}
	response, err := client.ViewUsersBySection(context.Background(), request)
	if err != nil {
		log.Fatalf("Error viewing users by section: %v", err)
	}

	// Process the response
	fmt.Println("Users in Section", section)
	for _, user := range response.User {
		fmt.Printf("UserID: %s, FirstName: %s, LastName: %s, Email: %s\n", user.UserId, user.FirstName, user.LastName, user.Email)
	}

	fmt.Println("Seat Map:")
	for userID, seat := range response.SeatMap {
		fmt.Printf("UserID: %s, Seat: %s\n", userID, seat)
	}
}

// removeUser interacts with the user to remove a user
func removeUser(client proto.TrainServiceClient, scanner *bufio.Scanner) {
	fmt.Print("Enter password to remove: ")
	if !scanner.Scan() {
		return
	}
	userID := scanner.Text()

	// Call gRPC method to remove user
	request := &proto.RemoveUserRequest{UserId: userID}
	response, err := client.RemoveUser(context.Background(), request)
	if err != nil {
		log.Fatalf("Error removing user: %v", err)
	}

	// Print or process the response accordingly
	fmt.Printf("Response: %+v\n", response)
}

// modifySeat interacts with the user to modify a user's seat
func modifySeat(client proto.TrainServiceClient, scanner *bufio.Scanner) {
	fmt.Print("Enter password: ")
	if !scanner.Scan() {
		return
	}
	userID := scanner.Text()

	fmt.Print("Enter new seat: ")
	if !scanner.Scan() {
		return
	}
	newSeat := scanner.Text()

	// Call gRPC method to modify user's seat
	request := &proto.ModifySeatRequest{UserId: userID, NewSeat: newSeat}
	response, err := client.ModifySeat(context.Background(), request)
	if err != nil {
		log.Fatalf("Error modifying seat: %v", err)
	}

	// Print or process the response accordingly
	fmt.Printf("Response: %+v\n", response)
}

// getReceiptForUser interacts with the user to get a receipt for a user
func getReceiptForUser(client proto.TrainServiceClient, scanner *bufio.Scanner) {
	fmt.Print("Enter password for receipt: ")
	if !scanner.Scan() {
		return
	}
	userID := scanner.Text()

	user := &proto.User{UserId: userID}
	// Call gRPC method to get receipt for user
	response, err := client.GetReceiptForUser(context.Background(), user)
	if err != nil {
		log.Fatalf("Error getting receipt: %v", err)
	}

	// Print or process the response accordingly
	fmt.Printf("Response: %+v\n", response)
}

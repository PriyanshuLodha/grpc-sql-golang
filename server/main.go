package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/priyanshu/trainservice-database/pkg/model"
	"github.com/priyanshu/trainservice-database/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// server represents the gRPC server for the TrainService
type server struct {
	mu      sync.Mutex
	users   map[string]*proto.User
	tickets map[string]*proto.Ticket
	proto.UnimplementedTrainServiceServer
}

// AddUser handles the gRPC AddUser request
func (s *server) AddUser(ctx context.Context, req *proto.User) (*proto.ModifySeatResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[req.UserId] = req
	user := &model.User{
		UserID:    req.UserId,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	// Create the user in the database
	user.CreateUser()
	return &proto.ModifySeatResponse{
		Success: true,
		Message: "User added successfully",
	}, nil
}

var (
	count int = 1
)

// Purchase handles the gRPC Purchase request
func (s *server) Purchase(ctx context.Context, req *proto.PurchaseRequest) (*proto.ReceiptResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if req.Ticket == nil || req.Ticket.User == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Ticket or User cannot be nil")
	}

	userID := req.Ticket.User.UserId

	if len(userID) < 7 {
		return nil, status.Errorf(codes.InvalidArgument, "User ID is too short")
	}

	seat := fmt.Sprintf("Section%s-%s", string(userID[0]), userID[len(userID)-1:])
	fmt.Print(seat)
	count++
	if count%2 == 0 {
		req.Ticket.Seat = "A"
	} else {
		req.Ticket.Seat = "B"
	}
	s.tickets[userID] = req.Ticket
	ticket := &model.Ticket{
		From:      req.Ticket.From,
		To:        req.Ticket.To,
		UserID:    req.Ticket.User.UserId,
		PricePaid: req.Ticket.PricePaid,
		Seat:      req.Ticket.Seat,
	}

	// Create the ticket in the database
	ticket.CreateTicket()
	return &proto.ReceiptResponse{
		Ticket: req.Ticket,
	}, nil
}

// ViewUsersBySection handles the gRPC ViewUsersBySection request
func (s *server) ViewUsersBySection(ctx context.Context, req *proto.ViewUsersBySectionRequest) (*proto.ViewUsersBySectionResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	usersInSection := make([]*proto.User, 0)
	seatMap := make(map[string]string)

	// Assume the section is included in the request
	requestedSection := req.Section
	fmt.Print(requestedSection + " ")
	// Loop through all users
	for userID, ticket := range s.tickets {
		// Check if the user has a valid section
		fmt.Print(ticket.Seat)
		if ticket.User != nil {
			fmt.Print("hii")
			// Check if the user's section matches the requested section
			if requestedSection == "" || ticket.Seat == requestedSection {
				usersInSection = append(usersInSection, ticket.User)
				seatMap[userID] = ticket.Seat
			}
		}
	}

	return &proto.ViewUsersBySectionResponse{
		User:    usersInSection,
		SeatMap: seatMap,
	}, nil
}

// RemoveUser handles the gRPC RemoveUser request
func (s *server) RemoveUser(ctx context.Context, req *proto.RemoveUserRequest) (*proto.ModifySeatResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.users, req.UserId)
	delete(s.tickets, req.UserId)
	return &proto.ModifySeatResponse{
		Success: true,
		Message: "User removed successfully",
	}, nil
}

// ModifySeat handles the gRPC ModifySeat request
func (s *server) ModifySeat(ctx context.Context, req *proto.ModifySeatRequest) (*proto.ModifySeatResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ticket, ok := s.tickets[req.UserId]; ok {
		ticket.Seat = req.NewSeat
		s.tickets[req.UserId] = ticket
		return &proto.ModifySeatResponse{
			Success: true,
			Message: "Seat modified successfully",
		}, nil
	}
	model.UpdateUserSection(req.UserId, req.NewSeat)
	return &proto.ModifySeatResponse{
		Success: false,
		Message: "User not found",
	}, nil
}

// GetReceiptForUser handles the gRPC GetReceiptForUser request
func (s *server) GetReceiptForUser(ctx context.Context, req *proto.User) (*proto.ReceiptResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ticket, ok := s.tickets[req.UserId]; ok {
		return &proto.ReceiptResponse{
			Ticket: ticket,
		}, nil
	}
	return nil, status.Errorf(codes.NotFound, "User has not purchased a ticket")
}

func main() {
	// Create a new gRPC server
	s := grpc.NewServer()
	// Register the TrainService server with the gRPC server
	proto.RegisterTrainServiceServer(s, &server{
		users:   make(map[string]*proto.User),
		tickets: make(map[string]*proto.Ticket),
	})

	// Listen for incoming connections on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	fmt.Println("Server listening on :50051")
	// Serve gRPC requests
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

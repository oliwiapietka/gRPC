package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "gobook/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GoBookServer implements the GoBookRPCServer interface.
type GoBookServer struct {
	pb.UnimplementedGoBookRPCServer
	store *Store
}

// NewGoBookServer creates a new GoBookServer with an empty store.
func NewGoBookServer() *GoBookServer {
	return &GoBookServer{
		store: NewStore(),
	}
}

// Register registers a new user on the server.
func (s *GoBookServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username cannot be empty")
	}

	ok := s.store.RegisterUser(req.Username, req.IpAddress)
	if !ok {
		return &pb.RegisterResponse{
			Success: false,
			Message: fmt.Sprintf("User '%s' is already registered", req.Username),
		}, nil
	}

	log.Printf("[REGISTER] User '%s' registered from %s", req.Username, req.IpAddress)
	return &pb.RegisterResponse{
		Success: true,
		Message: fmt.Sprintf("Welcome to GoBook, %s!", req.Username),
	}, nil
}

// Unregister removes a user from the server.
func (s *GoBookServer) Unregister(ctx context.Context, req *pb.UserRequest) (*pb.RegisterResponse, error) {
	ok := s.store.UnregisterUser(req.Username)
	if !ok {
		return &pb.RegisterResponse{
			Success: false,
			Message: fmt.Sprintf("User '%s' is not registered", req.Username),
		}, nil
	}

	log.Printf("[UNREGISTER] User '%s' has left GoBook", req.Username)
	return &pb.RegisterResponse{
		Success: true,
		Message: fmt.Sprintf("Goodbye, %s!", req.Username),
	}, nil
}

// GetOnlineUsers returns a list of all currently registered (online) users.
func (s *GoBookServer) GetOnlineUsers(ctx context.Context, req *pb.Empty) (*pb.OnlineUsersResponse, error) {
	users := s.store.GetOnlineUsers()
	return &pb.OnlineUsersResponse{Users: users}, nil
}

// GetUserInfo returns information about a specific user.
func (s *GoBookServer) GetUserInfo(ctx context.Context, req *pb.UserRequest) (*pb.UserInfoResponse, error) {
	user, exists := s.store.GetUserInfo(req.Username)
	if !exists {
		return &pb.UserInfoResponse{
			Username:  req.Username,
			IpAddress: "",
			Online:    false,
		}, nil
	}

	return &pb.UserInfoResponse{
		Username:  user.Username,
		IpAddress: user.IpAddress,
		Online:    true,
	}, nil
}

// CreatePost creates a new post on the news feed.
func (s *GoBookServer) CreatePost(ctx context.Context, req *pb.NewPostRequest) (*pb.Post, error) {
	if req.Author == "" || req.Content == "" {
		return nil, status.Error(codes.InvalidArgument, "author and content cannot be empty")
	}

	if !s.store.IsRegistered(req.Author) {
		return nil, status.Errorf(codes.PermissionDenied, "user '%s' is not registered", req.Author)
	}

	post := s.store.AddPost(req.Author, req.Content)
	log.Printf("[POST] User '%s' created post #%s", req.Author, post.Id)
	return post, nil
}

// GetNewsFeed returns all posts from the news feed.
func (s *GoBookServer) GetNewsFeed(ctx context.Context, req *pb.Empty) (*pb.NewsFeedResponse, error) {
	posts := s.store.GetNewsFeed()
	return &pb.NewsFeedResponse{Posts: posts}, nil
}

// LikePost adds a like to a post.
func (s *GoBookServer) LikePost(ctx context.Context, req *pb.LikeRequest) (*pb.Post, error) {
	post, alreadyLiked := s.store.LikePost(req.PostId, req.Username)
	if post == nil {
		return nil, status.Errorf(codes.NotFound, "post '%s' not found", req.PostId)
	}

	if alreadyLiked {
		log.Printf("[LIKE] User '%s' already liked post #%s", req.Username, req.PostId)
		return post, nil
	}

	log.Printf("[LIKE] User '%s' liked post #%s", req.Username, req.PostId)
	return post, nil
}

// CommentOnPost adds a comment to a post.
func (s *GoBookServer) CommentOnPost(ctx context.Context, req *pb.CommentRequest) (*pb.Post, error) {
	if req.Content == "" {
		return nil, status.Error(codes.InvalidArgument, "comment content cannot be empty")
	}

	if !s.store.IsRegistered(req.Author) {
		return nil, status.Errorf(codes.PermissionDenied, "user '%s' is not registered", req.Author)
	}

	post := s.store.CommentOnPost(req.PostId, req.Author, req.Content)
	if post == nil {
		return nil, status.Errorf(codes.NotFound, "post '%s' not found", req.PostId)
	}

	log.Printf("[COMMENT] User '%s' commented on post #%s", req.Author, req.PostId)
	return post, nil
}

// SendPM sends a private message to another user.
func (s *GoBookServer) SendPM(ctx context.Context, req *pb.PrivateMessage) (*pb.PMResponse, error) {
	if req.ToUser == "" || req.FromUser == "" {
		return nil, status.Error(codes.InvalidArgument, "from_user and to_user cannot be empty")
	}

	// Check if recipient exists
	_, exists := s.store.GetUserInfo(req.ToUser)
	if !exists {
		return &pb.PMResponse{
			Success: false,
			Message: fmt.Sprintf("User '%s' is not online", req.ToUser),
		}, nil
	}

	req.Timestamp = time.Now().Unix()
	s.store.StorePM(req)

	log.Printf("[PM] '%s' -> '%s'", req.FromUser, req.ToUser)
	return &pb.PMResponse{
		Success: true,
		Message: "Message sent!",
	}, nil
}

// GetMyPMs streams all pending private messages for a user.
func (s *GoBookServer) GetMyPMs(req *pb.UserRequest, stream pb.GoBookRPC_GetMyPMsServer) error {
	msgs := s.store.GetPMs(req.Username)
	for _, msg := range msgs {
		if err := stream.Send(msg); err != nil {
			return err
		}
	}
	return nil
}

package main

import (
	"context"
	"fmt"
	"io"
	"time"

	pb "gobook/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GoBookClient wraps the gRPC client stub and user info.
type GoBookClient struct {
	conn     *grpc.ClientConn
	client   pb.GoBookRPCClient
	username string
	ipAddr   string
}

// NewGoBookClient creates a new client connected to the given server address.
func NewGoBookClient(serverAddr, username string) (*GoBookClient, error) {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}

	client := pb.NewGoBookRPCClient(conn)
	return &GoBookClient{
		conn:     conn,
		client:   client,
		username: username,
		ipAddr:   serverAddr,
	}, nil
}

// Close closes the gRPC connection.
func (c *GoBookClient) Close() {
	c.conn.Close()
}

// Register registers the user with the GoBook server.
func (c *GoBookClient) Register() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.Register(ctx, &pb.RegisterRequest{
		Username:  c.username,
		IpAddress: c.ipAddr,
	})
	if err != nil {
		return fmt.Errorf("registration failed: %v", err)
	}

	if resp.Success {
		fmt.Printf("  %s\n", resp.Message)
	} else {
		fmt.Printf("  Registration failed: %s\n", resp.Message)
	}
	return nil
}

// Unregister removes the user from the GoBook server.
func (c *GoBookClient) Unregister() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.Unregister(ctx, &pb.UserRequest{Username: c.username})
	if err != nil {
		return fmt.Errorf("unregister failed: %v", err)
	}

	fmt.Printf("  %s\n", resp.Message)
	return nil
}

// ViewNewsFeed retrieves and displays the current news feed.
func (c *GoBookClient) ViewNewsFeed() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.GetNewsFeed(ctx, &pb.Empty{})
	if err != nil {
		return fmt.Errorf("failed to get news feed: %v", err)
	}

	if len(resp.Posts) == 0 {
		fmt.Println("\n  News feed is empty. Be the first to post!")
		return nil
	}

	fmt.Println("\n  ============ NEWS FEED ============")
	for _, post := range resp.Posts {
		t := time.Unix(post.Timestamp, 0).Format("15:04:05")
		fmt.Printf("\n  [#%s] %s (%s)\n", post.Id, post.Author, t)
		fmt.Printf("  %s\n", post.Content)
		fmt.Printf("  ♥ %d likes", post.Likes)
		if len(post.Comments) > 0 {
			fmt.Printf(" | %d comments", len(post.Comments))
		}
		fmt.Println()

		for _, comment := range post.Comments {
			ct := time.Unix(comment.Timestamp, 0).Format("15:04:05")
			fmt.Printf("    └─ %s (%s): %s\n", comment.Author, ct, comment.Content)
		}
	}
	fmt.Println("  ====================================")
	return nil
}

// CreatePost creates a new post on the news feed.
func (c *GoBookClient) CreatePost(content string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	post, err := c.client.CreatePost(ctx, &pb.NewPostRequest{
		Author:  c.username,
		Content: content,
	})
	if err != nil {
		return fmt.Errorf("failed to create post: %v", err)
	}

	fmt.Printf("  Post #%s created successfully!\n", post.Id)
	return nil
}

// LikePost likes a post by ID.
func (c *GoBookClient) LikePost(postID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	post, err := c.client.LikePost(ctx, &pb.LikeRequest{
		PostId:   postID,
		Username: c.username,
	})
	if err != nil {
		return fmt.Errorf("failed to like post: %v", err)
	}

	fmt.Printf("  Post #%s now has %d likes ♥\n", post.Id, post.Likes)
	return nil
}

// CommentOnPost adds a comment to a post.
func (c *GoBookClient) CommentOnPost(postID, content string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.client.CommentOnPost(ctx, &pb.CommentRequest{
		PostId:  postID,
		Author:  c.username,
		Content: content,
	})
	if err != nil {
		return fmt.Errorf("failed to comment: %v", err)
	}

	fmt.Println("  Comment added!")
	return nil
}

// GetOnlineUsers retrieves and displays online users.
func (c *GoBookClient) GetOnlineUsers() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.GetOnlineUsers(ctx, &pb.Empty{})
	if err != nil {
		return fmt.Errorf("failed to get online users: %v", err)
	}

	if len(resp.Users) == 0 {
		fmt.Println("\n  No users online.")
		return nil
	}

	fmt.Println("\n  ========= ONLINE USERS =========")
	for i, user := range resp.Users {
		fmt.Printf("  %d. %s (%s)\n", i+1, user.Username, user.IpAddress)
	}
	fmt.Println("  =================================")
	return nil
}

// SendPM sends a private message to another user via the server.
func (c *GoBookClient) SendPM(toUser, content string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.client.SendPM(ctx, &pb.PrivateMessage{
		FromUser:  c.username,
		ToUser:    toUser,
		Content:   content,
		Timestamp: time.Now().Unix(),
	})
	if err != nil {
		return fmt.Errorf("failed to send PM: %v", err)
	}

	if resp.Success {
		fmt.Printf("  PM sent to %s!\n", toUser)
	} else {
		fmt.Printf("  Failed: %s\n", resp.Message)
	}
	return nil
}

// GetMyPMs retrieves all pending private messages.
func (c *GoBookClient) GetMyPMs() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := c.client.GetMyPMs(ctx, &pb.UserRequest{Username: c.username})
	if err != nil {
		return fmt.Errorf("failed to get PMs: %v", err)
	}

	count := 0
	fmt.Println("\n  ======== PRIVATE MESSAGES ========")
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error receiving PM: %v", err)
		}
		count++
		t := time.Unix(msg.Timestamp, 0).Format("15:04:05")
		fmt.Printf("  [%s] From %s: %s\n", t, msg.FromUser, msg.Content)
	}

	if count == 0 {
		fmt.Println("  No new messages.")
	}
	fmt.Println("  ==================================")
	return nil
}

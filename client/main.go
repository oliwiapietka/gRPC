package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	serverAddr := flag.String("server", "localhost:50051", "GoBook server address")
	username := flag.String("username", "", "Your username")
	flag.Parse()

	if *username == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter your username: ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)
		if name == "" {
			log.Fatal("Username cannot be empty!")
		}
		username = &name
	}

	client, err := NewGoBookClient(*serverAddr, *username)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Register with the server
	if err := client.Register(); err != nil {
		log.Fatalf("Failed to register: %v", err)
	}

	// Start interactive CLI
	reader := bufio.NewReader(os.Stdin)
	running := true

	for running {
		fmt.Println()
		fmt.Println("=========================================")
		fmt.Printf("  GoBook - Logged in as: %s\n", *username)
		fmt.Println("=========================================")
		fmt.Println("  1. View News Feed")
		fmt.Println("  2. Create New Post")
		fmt.Println("  3. Like a Post")
		fmt.Println("  4. Comment on a Post")
		fmt.Println("  5. View Online Users")
		fmt.Println("  6. Send Private Message")
		fmt.Println("  7. Check My Messages")
		fmt.Println("  8. Exit")
		fmt.Println("=========================================")
		fmt.Print("  Choose option (1-8): ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			if err := client.ViewNewsFeed(); err != nil {
				fmt.Printf("  Error: %v\n", err)
			}

		case "2":
			fmt.Print("  Enter your post: ")
			content, _ := reader.ReadString('\n')
			content = strings.TrimSpace(content)
			if content != "" {
				if err := client.CreatePost(content); err != nil {
					fmt.Printf("  Error: %v\n", err)
				}
			}

		case "3":
			fmt.Print("  Enter post ID to like: ")
			postID, _ := reader.ReadString('\n')
			postID = strings.TrimSpace(postID)
			if postID != "" {
				if err := client.LikePost(postID); err != nil {
					fmt.Printf("  Error: %v\n", err)
				}
			}

		case "4":
			fmt.Print("  Enter post ID to comment on: ")
			postID, _ := reader.ReadString('\n')
			postID = strings.TrimSpace(postID)
			fmt.Print("  Enter your comment: ")
			comment, _ := reader.ReadString('\n')
			comment = strings.TrimSpace(comment)
			if postID != "" && comment != "" {
				if err := client.CommentOnPost(postID, comment); err != nil {
					fmt.Printf("  Error: %v\n", err)
				}
			}

		case "5":
			if err := client.GetOnlineUsers(); err != nil {
				fmt.Printf("  Error: %v\n", err)
			}

		case "6":
			fmt.Print("  Enter recipient username: ")
			toUser, _ := reader.ReadString('\n')
			toUser = strings.TrimSpace(toUser)
			fmt.Print("  Enter message: ")
			msg, _ := reader.ReadString('\n')
			msg = strings.TrimSpace(msg)
			if toUser != "" && msg != "" {
				if err := client.SendPM(toUser, msg); err != nil {
					fmt.Printf("  Error: %v\n", err)
				}
			}

		case "7":
			if err := client.GetMyPMs(); err != nil {
				fmt.Printf("  Error: %v\n", err)
			}

		case "8":
			client.Unregister()
			running = false
			fmt.Println("  Goodbye!")

		default:
			fmt.Println("  Invalid option. Please choose 1-8.")
		}
	}
}

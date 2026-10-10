package main

import (
	"fmt"
	"sync"
	"time"

	pb "gobook/pb"
)

// Store holds all in-memory data for the GoBook server.
type Store struct {
	mu      sync.RWMutex
	users   map[string]*pb.User            // username -> User
	posts   []*pb.Post                     // ordered list of posts
	pms     map[string][]*pb.PrivateMessage // username -> list of PMs received
	likedBy map[string]map[string]bool     // postID -> set of usernames who liked
	nextID  int                            // counter for post IDs
}

// NewStore creates a new empty Store.
func NewStore() *Store {
	return &Store{
		users:   make(map[string]*pb.User),
		posts:   make([]*pb.Post, 0),
		pms:     make(map[string][]*pb.PrivateMessage),
		likedBy: make(map[string]map[string]bool),
		nextID:  1,
	}
}

// RegisterUser adds a user to the store. Returns false if already registered.
func (s *Store) RegisterUser(username, ipAddr string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[username]; exists {
		return false
	}
	s.users[username] = &pb.User{
		Username:  username,
		IpAddress: ipAddr,
	}
	return true
}

// UnregisterUser removes a user from the store.
func (s *Store) UnregisterUser(username string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[username]; !exists {
		return false
	}
	delete(s.users, username)
	return true
}

// GetOnlineUsers returns all currently registered users.
func (s *Store) GetOnlineUsers() []*pb.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*pb.User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}
	return users
}

// GetUserInfo returns info about a specific user, or nil if not found.
func (s *Store) GetUserInfo(username string) (*pb.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, exists := s.users[username]
	return u, exists
}

// IsRegistered checks if a user is currently registered.
func (s *Store) IsRegistered(username string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.users[username]
	return exists
}

// AddPost creates a new post and adds it to the news feed.
func (s *Store) AddPost(author, content string) *pb.Post {
	s.mu.Lock()
	defer s.mu.Unlock()

	post := &pb.Post{
		Id:        fmt.Sprintf("%d", s.nextID),
		Author:    author,
		Content:   content,
		Timestamp: time.Now().Unix(),
		Likes:     0,
		Comments:  make([]*pb.Comment, 0),
	}
	s.nextID++
	s.posts = append(s.posts, post)
	return post
}

// GetNewsFeed returns all posts in reverse chronological order.
func (s *Store) GetNewsFeed() []*pb.Post {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy in reverse order (newest first)
	result := make([]*pb.Post, len(s.posts))
	for i, p := range s.posts {
		result[len(s.posts)-1-i] = p
	}
	return result
}

// LikePost increments the like count on a post. Returns the updated post, or nil if not found.
// Returns (post, alreadyLiked) — alreadyLiked is true if user already liked this post.
func (s *Store) LikePost(postID, username string) (*pb.Post, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, p := range s.posts {
		if p.Id == postID {
			// Check for duplicate like
			if s.likedBy[postID] == nil {
				s.likedBy[postID] = make(map[string]bool)
			}
			if s.likedBy[postID][username] {
				return p, true // already liked
			}
			s.likedBy[postID][username] = true
			p.Likes++
			return p, false
		}
	}
	return nil, false
}

// CommentOnPost adds a comment to a post. Returns the updated post or nil.
func (s *Store) CommentOnPost(postID, author, content string) *pb.Post {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, p := range s.posts {
		if p.Id == postID {
			comment := &pb.Comment{
				Author:    author,
				Content:   content,
				Timestamp: time.Now().Unix(),
			}
			p.Comments = append(p.Comments, comment)
			return p
		}
	}
	return nil
}

// StorePM stores a private message for the recipient.
func (s *Store) StorePM(pm *pb.PrivateMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.pms[pm.ToUser] = append(s.pms[pm.ToUser], pm)
}

// GetPMs returns and clears all PMs for a given user.
func (s *Store) GetPMs(username string) []*pb.PrivateMessage {
	s.mu.Lock()
	defer s.mu.Unlock()

	msgs := s.pms[username]
	s.pms[username] = nil
	return msgs
}

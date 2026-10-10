# GoBook — Distributed Social Network using gRPC

GoBook is a small distributed social networking platform built with **gRPC** and **Protocol Buffers** in **Go**, created as part of the CS1380 gRPC lab.

## Features

- **Register / Unregister** — join or leave GoBook
- **News Feed** — post messages, view, like and comment
- **Online Users** — see who's currently connected
- **Private Messages (PM)** — send and receive messages between users

## Project Structure

```
gRPC/
├── proto/
│   └── gobook.proto              # Protocol Buffers service definition
├── pb/
│   ├── gobook.pb.go              # Auto-generated message code
│   └── gobook_grpc.pb.go         # Auto-generated gRPC client/server stubs
├── server/
│   ├── main.go                   # Server entry point (port 50051)
│   ├── server.go                 # GoBookRPCServer implementation (10 RPC methods)
│   └── store.go                  # Thread-safe in-memory data store
├── client/
│   ├── main.go                   # Interactive CLI entry point
│   └── client.go                 # gRPC client wrapper
├── server.exe                    # Pre-built server binary
├── client.exe                    # Pre-built client binary
├── go.mod
└── go.sum
```

## Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [protoc](https://github.com/protocolbuffers/protobuf/releases) + Go plugins (only needed to regenerate proto)

## Running GoBook

### 1. Start the server

```bash
# Using pre-built binary:
.\server.exe

# Or run from source:
go run .\server\
```

The server starts on `localhost:50051`.

### 2. Connect clients (open new terminal for each)

```bash
# Using pre-built binary:
.\client.exe --username=Alice

# Or run from source:
go run .\client\ --username=Alice --server=localhost:50051
```

### 3. Client menu options

```
=========================================
  GoBook - Logged in as: Alice
=========================================
  1. View News Feed
  2. Create New Post
  3. Like a Post
  4. Comment on a Post
  5. View Online Users
  6. Send Private Message
  7. Check My Messages
  8. Exit
=========================================
```

## Regenerating proto (if you modify gobook.proto)

```bash
protoc -I proto --go_out=. --go_opt=module=gobook --go-grpc_out=. --go-grpc_opt=module=gobook proto/gobook.proto
```

## gRPC Methods

| RPC Method | Input | Output | Description |
|---|---|---|---|
| `Register` | `RegisterRequest` | `RegisterResponse` | Register user on server |
| `Unregister` | `UserRequest` | `RegisterResponse` | Remove user |
| `GetOnlineUsers` | `Empty` | `OnlineUsersResponse` | List online users |
| `GetUserInfo` | `UserRequest` | `UserInfoResponse` | Get user IP & status |
| `CreatePost` | `NewPostRequest` | `Post` | Add post to news feed |
| `GetNewsFeed` | `Empty` | `NewsFeedResponse` | Get all posts |
| `LikePost` | `LikeRequest` | `Post` | Like a post |
| `CommentOnPost` | `CommentRequest` | `Post` | Comment on a post |
| `SendPM` | `PrivateMessage` | `PMResponse` | Send private message |
| `GetMyPMs` | `UserRequest` | `stream PrivateMessage` | Receive pending PMs |

## Key gRPC Concepts Used

- **Protocol Buffers v3** (`proto3`) for serialization
- **Unary RPCs** for most operations (1 request → 1 response)
- **Server-side streaming** for `GetMyPMs` (1 request → stream of responses)
- **Blocking/synchronous** calls in Go
- **`context.Context`** parameter on every server method
- **`Empty` message** for RPCs with no meaningful input/output

# TTLock Go Client Library

This is a pure Go client library for the [TTLock Open API](https://open.ttlock.com/). It provides a simple and idiomatic way to interact with TTLock smart locks, gateways, and other devices.

## Features

- **Pure Go**: No external dependencies.
- **Automatic Authentication**: Handles OAuth2 token acquisition and automatic refreshing in the background.
- **Comprehensive Coverage**: Supports Lock, eKey, Passcode, and more (work in progress).
- **Error Handling**: Typed errors with specific error codes for robust handling.
- **Feature Flags**: Easy-to-use helpers to check lock capabilities.

## Installation

Since this project uses standard Go modules, you can import it directly:

```go
import "github.com/immofon/ttlock"
```

## Usage

### Initialization

Initialize the client with your TTLock Open Platform credentials. The client will automatically authenticate and start a background goroutine to keep the access token fresh.

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/immofon/ttlock"
)

func main() {
    // Replace with your actual credentials
    clientID := "your_client_id"
    clientSecret := "your_client_secret"
    username := "your_username"
    password := "your_password" // Plain text, library handles MD5 hashing

    // Initialize the client
    // This will immediately attempt to get an access token
    client := ttlock.NewClient(clientID, clientSecret, username, password)

    fmt.Println("Client initialized successfully!")
}
```

### Lock Management

#### List Locks

Retrieve a list of locks associated with the account.

```go
// Get list of locks (page 1, 20 items per page)
resp, err := client.GetLockList(1, 20, "", 0)
if err != nil {
    log.Fatalf("Failed to get lock list: %v", err)
}

for _, lock := range resp.List {
    fmt.Printf("Lock ID: %d, Name: %s, Mac: %s\n", lock.LockID, lock.LockName, lock.LockMac)
    
    // Check if lock supports a specific feature
    if lock.SupportsFeature(ttlock.LockFeaturePasscode) {
        fmt.Println(" - Supports Passcode")
    }
}
```

### eKey Management

#### Send an eKey

Send an electronic key (eKey) to another user.

```go
lockID := 12345
receiver := "receiver_username"
keyName := "Guest Key"
startDate := time.Now().UnixMilli()
endDate := time.Now().Add(24 * time.Hour).UnixMilli()

options := &ttlock.SendKeyOptions{
    Remarks: "Welcome!",
    RemoteEnable: 1, // Enable remote unlock
}

keyResp, err := client.SendKey(lockID, receiver, keyName, startDate, endDate, options)
if err != nil {
    log.Printf("Failed to send key: %v", err)
} else {
    fmt.Printf("Key sent successfully! Key ID: %d\n", keyResp.KeyID)
}
```

### Passcode Management

#### Generate a Random Passcode

Generate a random passcode for a lock.

```go
lockID := 12345
// Generate a 1-day valid passcode
startDate := time.Now().UnixMilli()
endDate := time.Now().Add(24 * time.Hour).UnixMilli()

passcodeResp, err := client.GetRandomPasscode(
    lockID, 
    ttlock.PasscodeTypePeriod, 
    "Airbnb Guest", 
    startDate, 
    endDate,
)

if err != nil {
    log.Printf("Failed to generate passcode: %v", err)
} else {
    fmt.Printf("Passcode: %s\n", passcodeResp.KeyboardPwd)
}
```

## Error Handling

The library provides a typed `Error` struct and helper functions to check for specific error codes.

```go
_, err := client.GetLockList(1, 20, "", 0)
if err != nil {
    // Check for specific error code
    if ttlock.IsErrorCode(err, ttlock.ErrLockFrozen) {
        fmt.Println("The lock is frozen.")
    } else if ttlock.IsErrorCode(err, ttlock.ErrTokenUnauthorized) {
        fmt.Println("Token is invalid or expired.")
    } else {
        // Generic error handling
        fmt.Printf("An error occurred: %v\n", err)
    }
}
```

## Feature Flags

Locks have a `featureValue` field that encodes their capabilities. You can check these using the `SupportsFeature` method on `Lock` or `LockDetail` structs.

```go
if lock.SupportsFeature(ttlock.LockFeatureRemoteUnlockConfig) {
    // Logic for locks that support remote unlock configuration
}

if lock.SupportsFeature(ttlock.LockFeatureFingerprint) {
    // Logic for locks that support fingerprints
}
```

See `feature.go` for a full list of `LockFeature` constants.

## CLI

The project includes a CLI tool located in `cmd/ttlock`.

To build and run the CLI:

```bash
go run ./cmd/ttlock [command] [flags]
```

Available commands:

- `lock`: Get lock details
  - `-id`: Lock ID
- `list-lock`: List locks
  - `-n`: Page number (default: 1)
  - `-s`: Page size (default: 20)
  - `-a`: Lock alias
  - `-g`: Group ID
- `list-passcode`: List passcodes
  - `-id`: Lock ID
  - `-n`: Page number (default: 1)
  - `-s`: Page size (default: 20)
  - `-o`: Order by (1:desc, 2:asc) (default: 1)
  - `-search` / `-q`: Search string
- `genpass`: Generate random passcode
  - `-id`: Lock ID
  - `-t`: Passcode type
  - `-n`: Passcode name
  - `-s`: Start date (YYYYMMDD-HH)
  - `-e`: End date (YYYYMMDD-HH)
- `sendkey`: Send eKey
  - `-id`: Lock ID
  - `-to`: Receiver username
  - `-n`: Key name
  - `-s`: Start date (YYYYMMDD-HH)
  - `-e`: End date (YYYYMMDD-HH)

## License

[MIT](LICENSE)

# Copilot Instructions for ttlock

This repository contains the Go client library for the TTLock API (`github.com/immofon/ttlock`).
The codebase is a pure Go implementation with no external dependencies.

## Architecture & Core Components

- **Client (`client.go`)**: The central entry point.
  - Initialize with `NewClient(clientID, clientSecret)`.
  - Manages `BaseURL` (default: `CNBaseURL`) and `HTTPClient`.
- **Authentication**:
  - `GetAccessToken(username, password)` handles OAuth2.
  - **Important**: The method automatically MD5 hashes the password. Pass the plain text password.
- **Error Handling (`errors.go`)**:
  - API errors are returned as `*Error` struct wrapping an `ErrorCode`.
  - Use `NewError(code)` to create errors.
  - Use `IsErrorCode(err, code)` to check for specific errors (e.g., `ErrLockFrozen`).
- **Feature Flags (`feature.go`)**:
  - Lock capabilities are encoded in a hex string (`featureValue`).
  - Use `HasFeature(featureValue, feature)` or the helper methods `lock.SupportsFeature(feature)` to check capabilities.
  - Features are defined as `LockFeature` constants (e.g., `LockFeaturePasscode`).

## Key Patterns & Conventions

### API Request Pattern

1.  **Endpoint**: Construct URL using `c.BaseURL + "/path"`.
2.  **Parameters**: Use `url.Values`.
    - Always include `clientId` and `accessToken`.
    - Most endpoints require `date` (current Unix milliseconds): `strconv.FormatInt(time.Now().UnixMilli(), 10)`.
3.  **Response**:
    - Unmarshal JSON into a specific response struct (e.g., `LockListResponse`).
    - **Validation**: Check `resp.Errcode != 0`. If non-zero, return `NewError(ErrorCode(resp.Errcode))`.

### Data Structures

- **Responses**: Structs usually contain a `list` field for collections and metadata (`pageNo`, `total`).
- **Models**: Core models like `Lock`, `LockDetail` map directly to JSON fields; both `Lock` and `LockDetail` fields include Chinese inline comments mirroring TTLock response parameter descriptions.

## Development Workflow

- **Dependencies**: Standard library only. No `go.mod` dependencies to manage.
- **Testing**: There are currently **NO** unit tests (`_test.go`).
  - _Action_: When adding new features, consider adding a test file if possible, or verify manually.
- **Formatting**: Follow standard Go conventions (`gofmt`).

## Example Usage

```go
// Initialize
client := ttlock.NewClient("client_id", "client_secret","username", "password")

// Check Feature
if lock.SupportsFeature(ttlock.LockFeatureRemoteUnlockConfig) {
    // Perform remote unlock logic
}

// Iterate Locks
iter := client.IterateLocks(token, "", 0)
for {
    lock, err := iter.Next()
    if err != nil {
        // Handle error
        break
    }
    if lock == nil {
        break
    }
    // Process lock
}
```

### Iterators

For endpoints that support pagination (`pageNo`, `pageSize`), iterator methods are provided to simplify access.

- `IterateLocks(accessToken, lockAlias, groupId)` returns a `*LockIterator`.
- `IteratePasscodes(accessToken, lockID, orderBy, searchStr)` returns a `*PasscodeIterator`.

Usage:

```go
iter := client.IterateLocks(token, "", 0)
for {
    lock, err := iter.Next()
    if err != nil {
        // Handle error
        break
    }
    if lock == nil {
        // End of list
        break
    }
    // Process lock
}
```

After each task, please update the relevant documentation files to reflect any changes made to the codebase, especially the .github/copilot-instructions.md and README.md file to ensure that Copilot has the most up-to-date information about the project structure and functionality.

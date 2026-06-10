# Go Syntax Reference

---

## Project Structure

```
my-project/
├── go.mod        ← module definition (like package.json)
├── main.go       ← entry point
└── handlers.go   ← other .go files in same package
```

```bash
# Initialize a new module
go mod init github.com/zainbharde/url-shortener

# Run
go run main.go

# Build binary
go build -o app main.go
```

---

## Imports

```go
import "fmt"                      // single import

import (                          // grouped import
    "fmt"
    "log"
    "net/http"
    "errors"
)
```

---

## Variables

```go
// Explicit type
var name string = "Zain"
var age int = 22

// Inferred type (most common)
name := "Zain"
age := 22

// Multiple assignment
x, y := 10, 20

// Constants
const MaxRetries = 3

// Zero values — Go initializes everything, no null surprises
var s string   // ""
var n int      // 0
var b bool     // false
```

---

## Arrays and Slices

```go
// Array — fixed size, rarely used directly
var arr [3]int = [3]int{1, 2, 3}

// Slice — dynamic, what you'll use almost always
nums := []int{1, 2, 3, 4, 5}

// Append to a slice
nums = append(nums, 6)

// Slice a slice
nums[1:3]   // [2, 3] — index 1 up to but not including 3

// Length
len(nums)   // 6

// Make a slice with length and capacity
s := make([]int, 0, 10)

// Iterating
for i, v := range nums {
    fmt.Println(i, v)   // index, value
}

// Ignore index
for _, v := range nums {
    fmt.Println(v)
}
```

---

## Maps

```go
// Declare and initialize
ages := map[string]int{
    "Zain": 22,
    "Alex": 30,
}

// Make (empty map)
scores := make(map[string]int)

// Set
scores["math"] = 95

// Get
val := scores["math"]   // 95

// Check if key exists — always do this
val, ok := scores["science"]
if !ok {
    fmt.Println("key not found")
}

// Delete
delete(scores, "math")

// Iterate
for key, val := range scores {
    fmt.Println(key, val)
}
```

---

## Loops

```go
// Go only has `for` — no while keyword

// Standard for loop
for i := 0; i < 5; i++ {
    fmt.Println(i)
}

// While-style
n := 0
for n < 5 {
    n++
}

// Infinite loop
for {
    // break out when needed
    break
}

// Range over slice
for i, v := range []string{"a", "b", "c"} {
    fmt.Println(i, v)
}

// Range over map
for k, v := range map[string]int{"a": 1} {
    fmt.Println(k, v)
}
```

---

## Conditional Logic

```go
// Standard if/else
if age >= 18 {
    fmt.Println("adult")
} else if age >= 13 {
    fmt.Println("teen")
} else {
    fmt.Println("child")
}

// If with init statement — very common in Go
if err := doSomething(); err != nil {
    log.Fatal(err)
}

// Switch
switch day {
case "Monday":
    fmt.Println("start of week")
case "Friday":
    fmt.Println("end of week")
default:
    fmt.Println("midweek")
}
```

---

## Functions

```go
// Basic function
func add(a int, b int) int {
    return a + b
}

// Multiple return values — very common in Go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("cannot divide by zero")
    }
    return a / b, nil
}

// Calling a multi-return function
result, err := divide(10, 2)
if err != nil {
    log.Fatal(err)
}

// Variadic function (variable number of args)
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}
```

---

## Error Handling

```go
// The Go pattern — errors are just values, always check them
result, err := someFunction()
if err != nil {
    // handle it — don't ignore it
    log.Fatal(err)       // logs and exits
    // or
    return err           // pass it up the call stack
    // or
    fmt.Println(err)     // handle gracefully
}

// Creating errors
err := errors.New("something went wrong")

// Formatted errors
err := fmt.Errorf("user %d not found", userID)

// Wrapping errors (adds context while preserving original)
err := fmt.Errorf("createShortCode: %w", originalErr)

// Unwrapping to check error type
if errors.Is(err, sql.ErrNoRows) {
    // handle "not found" case specifically
}
```

---

## User Defined Types — Structs

```go
// Define a struct
type User struct {
    ID        int
    Name      string
    Email     string
    CreatedAt time.Time
}

// Instantiate
u := User{
    ID:    1,
    Name:  "Zain",
    Email: "zainbharde@gmail.com",
}

// Access fields
fmt.Println(u.Name)

// Pointer to struct (common — avoids copying)
u := &User{Name: "Zain"}

// Struct with JSON tags — needed for encoding/decoding JSON
type ShortenRequest struct {
    URL        string `json:"url"`
    CustomCode string `json:"custom_code,omitempty"` // omitempty skips field if empty
}
```

---

## Methods on Structs

```go
type Rectangle struct {
    Width  float64
    Height float64
}

// Value receiver — use when you don't need to modify the struct
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Pointer receiver — use when you need to modify the struct
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

rect := Rectangle{Width: 10, Height: 5}
fmt.Println(rect.Area())   // 50
rect.Scale(2)
fmt.Println(rect.Area())   // 200
```

---

## Pointers

```go
x := 10
p := &x          // p is a pointer to x (holds the memory address)
fmt.Println(*p)  // dereference — prints 10

*p = 20          // modifies x through the pointer
fmt.Println(x)   // 20

// Why this matters — without pointer, original is unchanged
func addOne(n int) {
    n++   // only modifies local copy
}

// With pointer, original IS changed
func addOne(n *int) {
    *n++
}

addOne(&x)
```

---

## Interfaces

```go
// Define behavior, not data
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Any struct that implements these methods satisfies the interface
// No explicit "implements" keyword needed

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

// Now Circle satisfies Shape
func printShapeInfo(s Shape) {
    fmt.Println("Area:", s.Area())
}
```

---

## JSON Encoding / Decoding

```go
import "encoding/json"

// Struct to JSON (encode)
type Response struct {
    ShortCode string `json:"short_code"`
    URL       string `json:"url"`
}

r := Response{ShortCode: "abc123", URL: "https://google.com"}
data, err := json.Marshal(r)
// data is []byte — convert to string with string(data)

// JSON to struct (decode) — what you'll use for incoming requests
var req ShortenRequest
err := json.NewDecoder(r.Body).Decode(&req)
if err != nil {
    // handle bad JSON
}
fmt.Println(req.URL)
```

---

## Basic HTTP Server

```go
import "net/http"

// Handler function signature
func handleShorten(w http.ResponseWriter, r *http.Request) {
    // w is what you write your response to
    // r is the incoming request

    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)        // 201
    w.Write([]byte(`{"short_code":"abc123"}`))
}

// Register routes and start server
func main() {
    http.HandleFunc("/shorten", handleShorten)
    http.HandleFunc("/stats/", handleStats)

    log.Println("server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## Database (database/sql + postgres)

```go
import (
    "database/sql"
    _ "github.com/lib/pq"   // postgres driver — blank import for side effects
)

// Connect
db, err := sql.Open("postgres", "postgres://user:password@localhost/dbname?sslmode=disable")
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// Ping to verify connection
if err := db.Ping(); err != nil {
    log.Fatal(err)
}

// Query single row
var originalURL string
err = db.QueryRow(
    "SELECT original_url FROM links WHERE short_code = $1", code,
).Scan(&originalURL)
if err == sql.ErrNoRows {
    // not found
}

// Query multiple rows
rows, err := db.Query("SELECT short_code, original_url FROM links")
defer rows.Close()
for rows.Next() {
    var code, url string
    rows.Scan(&code, &url)
}

// Insert
_, err = db.Exec(
    "INSERT INTO links (short_code, original_url, clicks) VALUES ($1, $2, $3)",
    code, url, 0,
)

// Update
_, err = db.Exec(
    "UPDATE links SET clicks = clicks + 1 WHERE short_code = $1", code,
)
```

---

## Useful Standard Library Packages

| Package | What it does |
|---|---|
| `fmt` | Printing, formatting strings |
| `log` | Logging with timestamps, `log.Fatal` exits |
| `errors` | Creating and wrapping errors |
| `net/http` | HTTP server and client |
| `encoding/json` | JSON encode/decode |
| `database/sql` | Database interface |
| `math/rand` | Random number generation |
| `strings` | String manipulation |
| `time` | Time, timestamps, durations |
| `os` | Environment variables, file I/O |
| `strconv` | String ↔ int/float conversions |

---

## Things That Will Bite You Early

```go
// 1. Unused imports are a compile error — remove them
import "fmt"   // error if fmt is never used

// 2. Unused variables are a compile error
x := 10   // error if x is never used

// 3. := vs = 
x := 10   // declare AND assign (only for new variables)
x = 20    // assign only (variable must already exist)

// 4. nil vs zero value — Go doesn't have null, but pointers/maps/slices can be nil
var p *int      // nil pointer
var m map[string]int  // nil map — reading ok, writing will panic
m = make(map[string]int)  // initialize before writing

// 5. defer runs at end of function — useful for cleanup
func main() {
    db, _ := sql.Open(...)
    defer db.Close()   // runs when main() returns, no matter what
}
```

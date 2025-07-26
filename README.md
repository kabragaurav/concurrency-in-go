# Concurrency In Golang
by Gaurav Kabra

## Parallelism vs Concurrency
- True parallel tasks execute independently and simultaneously. Thus, require multiple CPUs.
    ![](./assets/ActivityMonitorProcesses.png)
- Concurrent tasks are interleaving and **non-deterministic** in order of execution.
My computer does not have 425 cores, yet able to run 425 processes.

Concurrency is handled by Go Runtime.

![](./assets/GoRuntime.png)

---

## Go Routines
- Abstraction over threads
- In general, `# of go routines > # of threads`
![](./assets/GoRoutinesAndThreads.png)
- Like threads, goroutines share same address space.

```go
go f(x)
```

Above syntax makes `f(x)` run on a goroutine.
E.g.,

```go
package main

import "fmt"

func main() {
	go hello()
	bye()
}

func hello() {
	fmt.Println("hello")
}

func bye() {
	fmt.Println("bye")
}

/**
Can produce different outputs like:

1.
bye

2.
bye
hello
*/
```

This is because a go routine does not block main goroutine.
One of the common ways to keep main goroutine alive is by using `time.Sleep()`:

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	go hi()
	// NEVER use Sleep for aliveness in prod
	time.Sleep(1 * time.Second)
	tata()
}

func hi() {
	fmt.Println("hi")
}

func tata() {
	fmt.Println("tata")
}

/**
In general, produces:
hi
tata
 */
```

### `sync.waitGroup`
- waits for goroutines to finish
- Under the hood, it keep counter for number of goroutines to finish
- All types in `sync` package MUST be passed as pointers to functions

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1) // need to wait for one goroutine before exiting main goroutine
	go hey(&wg)
	wg.Wait()
	seeyou()
}

func hey(wg *sync.WaitGroup) {
	defer wg.Done() // equivalent to wg.Add(-1)
	fmt.Println("hey")
}

func seeyou() {
	fmt.Println("seeyou")
}

/**
Deterministically produces output:
hey
seeyou
*/
```

### Race Condition
- Multiple goroutine try CUD on shared data (critical section) simultaneously
- `-race` flag (race detector) can be used to print exact cause of panic. E.g. `go run -race main.go`
- The `Map` panics on concurrent CUD operations. Hence `sync.Map` should be used

```go
func (m *Map) Load(key interface{}) (value interface{}, ok bool)
func (m *Map) Store(key, value interface{})
func (m *Map) Range(f func(key, value interface{})) bool  // calls f() for all (K, V) pairs
```

- Locks: `sync.Mutex`

```go
func (m *Mutex) Lock()
func (m *Mutex) Unlock()
```

---

## Channels
- to pass values between goroutines that don't directly call each other (yes, channels are usable only by goroutines)
- One way is to pass pointers (communicating by shared memory)
- Channels are other way. They are FIFO queues
- `chan` is reserved keyword and `<-` is channel operator
- A channel is associated to only one data-type

```go
ch := make(chan T)
ch <- data      // sending
data := <- ch   // receiving
```

- Sending and receiving are both blocking
- Types of channels:
  - Unbuffered: Zero capacity. Sender and receiver both must be present (synchronous)
  - Buffered: Pre-defined capacity. Abstraction over array (asynchronous)
- Channels are always passed as pointers and hence no need to pass as `&ch`. Just `ch` suffices
- channels are first-class citizen (part of `builtin.go`) and hence need no import

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string) // unbuffered channel
	go greet(ch)
	time.Sleep(1 * time.Second)
	fmt.Println(<-ch)
}

func greet(ch chan string) {
	fmt.Println("Greet starting")
	ch <- "Hello"
	fmt.Println("Greet ending")
}

/**
Since channel is unbuffered, sender (greet goroutine) wait for receiver (main goroutine)

Hence output is just:
Greet starting
Hello
*/
```

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string, 1) // unbuffered channel
	go aloha(ch)
	time.Sleep(1 * time.Second)
	fmt.Println(<-ch)
}

func aloha(ch chan string) {
	fmt.Println("Aloha starting")
	ch <- "Aloha"
	fmt.Println("Aloha ending")
}

/**
Since channel is buffered, sender (greet goroutine) DOES NOT wait for receiver (main goroutine)

Hence output is:
Aloha starting
Aloha ending
Aloha
*/
```

- Channel directions:
  - Bidirectional (default, `chan T`)
  - Unidirectional (`chan<- T` or `<-chan T`)
  - the default (bidirectional) is implicitly casted to unidirectional while sending/receiving
- Close a channel as `close(ch)`
  - Sending causes panic
  - Receiving returns zero-value of channel data-type once closed channel does not have anymore values
  - Closing again causes panic

#### What will be possible outputs of below programs?
1. 
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)
	greetings := []string{"Hello", "Hi", "Hey", "Hola", "Aloha"}
	go sendToChannel(ch, greetings)
	time.Sleep(2 * time.Second)
	for {
		greeting := <-ch
		fmt.Println("Receive from channel", greeting)
	}
}

func sendToChannel(ch chan string, greetings []string) {
	for _, greeting := range greetings {
		ch <- greeting
	}
}
```

Ans.:

```
Receive from channel Hello
Receive from channel Hi
Receive from channel Hey
Receive from channel Hola
Receive from channel Aloha
fatal error: all goroutines are asleep - deadlock!
```

Reason: `greeting := <-ch` does not know channel is exhausted

2. 
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)
	greetings := []string{"Hello", "Hi", "Hey", "Hola", "Aloha"}
	go sendToChannel(ch, greetings)
	time.Sleep(2 * time.Second)
	for {
		greeting := <-ch
		fmt.Println("Receive from channel", greeting)
	}
}

func sendToChannel(ch chan string, greetings []string) {
	for _, greeting := range greetings {
		ch <- greeting
	}
	close(ch)   // Added close
}
```

Ans.:

```
Receive from channel Hello
Receive from channel Hi
Receive from channel Hey
Receive from channel Hola
Receive from channel Aloha
Receive from channel 
Receive from channel 
Receive from channel
...
```

Reason: Once channel exhausts, receiving from closed channel gives zero-value of channel type (here "" for string)

#### How to fix above code?

Way 1: Using `ok`

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)
	greetings := []string{"Hello", "Hi", "Hey", "Hola", "Aloha"}
	go sendToChannel(ch, greetings)
	time.Sleep(2 * time.Second)
	for {
		greeting, ok := <-ch
		if !ok {
			break
		}
		fmt.Println("Receive from channel", greeting)
	}
}

func sendToChannel(ch chan string, greetings []string) {
	for _, greeting := range greetings {
		ch <- greeting
	}
	close(ch)
}
```

Way 2: Using `for range` loop

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)
	greetings := []string{"Hello", "Hi", "Hey", "Hola", "Aloha"}
	go sendToChannel(ch, greetings)
	time.Sleep(2 * time.Second)
	for greeting := range ch {
		fmt.Println("Receive from channel", greeting)
	}
}

func sendToChannel(ch chan string, greetings []string) {
	for _, greeting := range greetings {
		ch <- greeting
	}
	close(ch)
}
```
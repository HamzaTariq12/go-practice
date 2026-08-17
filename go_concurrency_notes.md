# Go Concurrency Notes

## 1. Two Important Concurrency Concepts

The two important concepts to understand first are:

1. **Deadlock**
2. **Race Condition**

---

# 2. Channels

Channels allow goroutines to communicate with each other.

There are two common types of channels:

- **Unbuffered channel**
- **Buffered channel**

---

## 2.1 Unbuffered Channel

An unbuffered channel has **no storage/capacity**.

```go
ch := make(chan int)
```

Because there is no storage, the sender and receiver must synchronize directly.

### Sending

```go
ch <- 5
```

### Receiving

```go
value := <-ch
```

The important rule is:

> The send blocks until a receiver is ready, and the receive blocks until a sender is ready.

---

## Why Do We Usually Use a Goroutine With an Unbuffered Channel?

Because the sender and receiver need to be able to execute **concurrently**.

For example:

```go
ch := make(chan int)

go func() {
    ch <- 5
}()

value := <-ch

fmt.Println(value)
```

The flow is:

```text
Goroutine                         Main goroutine

    ch <- 5
       |
       |  waits for receiver
       |
       +----------------------->  value := <-ch
                                  |
                                  v
                                  5
```

The goroutine can wait at:

```go
ch <- 5
```

while the main goroutine waits at:

```go
<-ch
```

Once both sides are ready, the value is transferred directly.

### Why is this useful?

An unbuffered channel has **zero capacity**, so the value cannot simply sit inside the channel.

The sender and receiver effectively perform a **direct handoff**.

Think of it like handing someone an object:

```text
Sender                         Receiver

  5  -------------------------->  receives 5
         direct handoff
```

---

## What Happens Without a Goroutine?

Consider:

```go
ch := make(chan int)

ch <- 5

value := <-ch
```

This gets stuck at:

```go
ch <- 5
```

Why?

Because the channel has no storage and there is currently no receiver.

The sender is saying:

> "I want to send 5."

But nobody is receiving yet.

Therefore the send blocks.

The program never reaches:

```go
value := <-ch
```

This can result in:

```text
fatal error: all goroutines are asleep - deadlock!
```

---

## Important Clarification

An unbuffered channel does **not technically require a goroutine**.

What it requires is that the **sender and receiver must be able to execute concurrently**.

In Go, goroutines are the normal way to achieve this.

For example:

```go
ch := make(chan int)

go func() {
    ch <- 5
}()

value := <-ch
```

Here the sender and receiver are running in separate goroutines.

---

# 3. Buffered Channels

A buffered channel has storage/capacity.

```go
ch := make(chan int, 3)
```

The `3` means the channel can hold up to 3 values.

Therefore this works:

```go
ch <- 1
ch <- 2
ch <- 3
```

The values are temporarily stored in the channel.

```text
+-------------------+
|  1  |  2  |  3    |
+-------------------+
      capacity = 3
```

However, this:

```go
ch <- 4
```

will block because the buffer is full.

A receiver must consume a value first:

```go
value := <-ch
```

Then there is space for another value.

---

## Buffered Channel Rules

### Sending

```text
Buffer NOT FULL
        |
        v
    ch <- value
        |
     succeeds
```

```text
Buffer FULL
     |
     v
ch <- value
     |
   blocks
```

### Receiving

```text
Buffer NOT EMPTY
        |
        v
    <-ch
        |
     succeeds
```

```text
Buffer EMPTY
     |
     v
    <-ch
     |
   blocks
```

---

# 4. Unbuffered vs Buffered

| Feature | Unbuffered | Buffered |
|---|---|---|
| Example | `make(chan int)` | `make(chan int, 3)` |
| Storage | None | Has storage |
| Capacity | `0` | Specified capacity |
| Sender waits for receiver | Yes | Only when buffer is full |
| Receiver waits for sender | Yes, if buffer is empty | Yes, if buffer is empty |
| Common use | Direct synchronization | Temporary storage between goroutines |

A useful mental model:

```text
UNBUFFERED

Sender -----------------> Receiver
        direct handoff
```

```text
BUFFERED

Sender ---> [ 1 | 2 | 3 ] ---> Receiver
              buffer
```

---

# 5. Deadlock

A deadlock occurs when goroutines are blocked and there is no possible way for them to make progress.

Example:

```go
ch := make(chan int)

ch <- 5
```

There is no receiver, so the send blocks forever.

Go may report:

```text
fatal error: all goroutines are asleep - deadlock!
```

## Important

A blocked channel operation does **not automatically mean there is a deadlock**.

It is a deadlock when the blocked operation has no possible way to proceed.

For example:

```go
ch := make(chan int)

go func() {
    ch <- 5
}()

value := <-ch
```

The goroutine initially blocks on the send, but the main goroutine eventually receives the value.

Therefore there is no deadlock.

---

# 6. Race Condition

A race condition occurs when multiple goroutines access the same shared data concurrently, at least one access is a write, and the accesses are not properly synchronized.

Example:

```go
var counter int

go func() {
    counter++
}()

go func() {
    counter++
}()
```

The problem is that:

```go
counter++
```

is not one indivisible operation.

Conceptually, it involves:

```text
1. Read counter
2. Add 1
3. Write counter
```

Two goroutines can interfere with each other.

For example:

```text
counter = 0

Goroutine 1              Goroutine 2

read 0                   read 0
add 1                    add 1
write 1                  write 1

Final counter = 1
```

You may have expected:

```text
2
```

but got:

```text
1
```

This is a race condition.

---

# 7. Mutex

A `sync.Mutex` can protect shared data.

```go
var mutex sync.Mutex
var counter int

mutex.Lock()
counter++
mutex.Unlock()
```

When multiple goroutines use the same mutex:

```text
Goroutine 1              Goroutine 2

   Lock                     |
     |                    WAIT
     |                      |
 counter++                WAIT
     |                      |
   Unlock                   |
                            |
                          Lock
                            |
                         counter++
                            |
                          Unlock
```

Only one goroutine can enter the protected section at a time.

---

# 8. WaitGroup

A `sync.WaitGroup` is mainly used to **wait for goroutines to finish**.

It is NOT primarily a solution for race conditions.

Example:

```go
var wg sync.WaitGroup

wg.Add(1)

go func() {
    defer wg.Done()

    // Work
}()

wg.Wait()
```

The methods mean:

```text
wg.Add(1)
    ↓
"I am expecting one goroutine to finish."

wg.Done()
    ↓
"One goroutine has finished."

wg.Wait()
    ↓
"Wait until the counter reaches zero."
```

---

## WaitGroup Does NOT Automatically Prevent Race Conditions

This is still unsafe:

```go
var counter int
var wg sync.WaitGroup

wg.Add(2)

go func() {
    counter++
    wg.Done()
}()

go func() {
    counter++
    wg.Done()
}()

wg.Wait()
```

`WaitGroup` only ensures that both goroutines have finished.

It does not make:

```go
counter++
```

safe.

You would need synchronization such as a `Mutex`, an atomic operation, or a channel-based design.

---

# 9. Mutex vs WaitGroup

| Tool | Main Purpose |
|---|---|
| `sync.Mutex` | Protect shared data |
| `sync.WaitGroup` | Wait for goroutines to finish |
| Channel | Communicate/synchronize between goroutines |

Think of them as:

```text
Mutex
  ↓
"Only one goroutine can access this critical section."

WaitGroup
  ↓
"Wait until these goroutines are finished."

Channel
  ↓
"Send data / synchronize between goroutines."
```

---

# 10. Quick Summary

## Unbuffered Channel

```go
ch := make(chan int)
```

- No storage.
- Capacity is `0`.
- Sender and receiver must synchronize.
- The sender blocks until a receiver is ready.
- The receiver blocks until a sender is ready.
- Usually used with goroutines so sender and receiver can execute concurrently.

## Buffered Channel

```go
ch := make(chan int, 3)
```

- Has storage.
- Capacity is `3`.
- Sender can continue while there is space.
- Sender blocks when the buffer is full.
- Receiver blocks when the buffer is empty.

## Deadlock

> Goroutines are blocked and there is no possible way for them to make progress.

## Race Condition

> Multiple goroutines access shared data concurrently without proper synchronization, causing unpredictable results.

## Mutex

> Protects shared data and critical sections.

## WaitGroup

> Waits for goroutines to finish. It does not automatically prevent race conditions.

## Channel

> Allows goroutines to communicate and synchronize with each other.

---

# 11. The Mental Model to Remember

```text
                 GO CONCURRENCY
                       |
        +--------------+--------------+
        |              |              |
     CHANNEL          MUTEX        WAITGROUP
        |              |              |
   Communication    Protect data    Wait for
   & synchronization               goroutines
        |
   +----+----+
   |         |
Unbuffered  Buffered
   |         |
No storage  Has storage
   |         |
Direct      Temporary
handoff     storage
```

The most important sentence to remember:

> **An unbuffered channel is a direct handoff: the sender cannot complete the send until a receiver is ready.**

That is why goroutines are so commonly used with unbuffered channels.

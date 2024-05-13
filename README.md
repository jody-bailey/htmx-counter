**HTMX Counter**
================
A simple Go application that uses HTMX and server-ide events (SSE) to create a counter that increments, decrements, and resets on the client-side.

**Getting Started**
-------------------

1. Clone this repository: `git clone https://github.com/jody-bailey/htmx-counter.git`
2. Install dependencies: `go mod tidy` (if you have Go Modules enabled)
3. Run the application: `go run main.go`

**How it Works**
---------------

The counter is a simple HTML form that increments, decrements, and resets when submitted to one of the following endpoints:

* `/increment`: Increments the counter on the server-side.
* `/decrement`: Decrements the counter on the server-side.
* `/reset`: Resets the counter to its initial value on the server-side.

When any of these operations are performed, the server sends a Server-Side Event (SSE) notification to the client with the updated count. The client updates the counter display accordingly.

Additionally, the application includes a timer that updates every second and publishes the elapsed time to the "timer" event using SSE. This allows the client to stay in sync with the server's clock.

**Endpoints**
------------

* `/`: The main endpoint that handles form submissions and returns the current counter value.
* `/sse`: The SSE endpoint that sends updates to the client when the counter changes, including timer events.
* `/increment`: Increments the counter on the server-side.
* `/decrement`: Decrements the counter on the server-side.
* `/reset`: Resets the counter to its initial value on the server-side.

**Technical Details**
-------------------

* This application uses Go 1.17 or later.
* HTMX is used for client-side rendering and handling form submissions.
* Server-Side Events (SSE) are used to push updates from the server to the client.
* The counter state is stored in memory; it's not persisted between restarts.


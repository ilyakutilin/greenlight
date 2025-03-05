/*
NOTE FOR THE RATE LIMITING MIDDLEWARE

Using the rate-limiting pattern provided in this package will only work if the API
application is running on a single-machine. If the infrastructure is distributed, with
the application running on multiple servers behind a load balancer, then an alternative
approach is required.

If HAProxy or Nginx as a load balancer or reverse proxy are used, both of these have
built-in functionality for rate limiting that it would probably be sensible to use.
Alternatively, it is possible to use a fast database like Redis to maintain a request
count for clients, running on a server which all the application servers can communicate
with.
*/

package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Middleware recovering the panic.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic
		// as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a panic or
			// not.
			if err := recover(); err != nil {
				// If there was a panic, set a "Connection: close" header on the
				// response. This acts as a trigger to make Go's HTTP server
				// automatically close the current connection after a response has been
				// sent.
				w.Header().Set("Connection", "close")
				// The value returned by recover() has the type any, so we use
				// fmt.Errorf() to normalize it into an error and call our
				// serverErrorResponse() helper. In turn, this will log the error using
				// our custom Logger type at the ERROR level and send the client a 500
				// Internal Server Error response.
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// Middleware to enforce an IP-based rate limit.
func (app *application) rateLimit(next http.Handler) http.Handler {
	// Define a client struct to hold the rate limiter and last seen time for each
	// client.
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	// Declare a mutex and a map to hold the clients' IP addresses and rate limiters.
	var (
		mu sync.Mutex
		// The values of the map are pointers to a client struct.
		clients = make(map[string]*client)
	)

	// Launch a background goroutine which removes old entries from the clients map once
	// every minute.
	go func() {
		for {
			time.Sleep(time.Minute)

			// Lock the mutex to prevent any rate limiter checks from happening while
			// the cleanup is taking place.
			mu.Lock()

			// Loop through all clients. If they haven't been seen within the last three
			// minutes, delete the corresponding entry from the map.
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			// IMPORTANT: Unlock the mutex when the cleanup is complete.
			mu.Unlock()
		}
	}()

	// The function we are returning is a closure.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only carry out the check if rate limiting is enabled.
		if app.config.limiter.enabled {
			// Extract the client's IP address from the request.
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			// Lock the mutex to prevent this code from being executed concurrently.
			mu.Lock()

			// Check to see if the IP address already exists in the map. If it doesn't,
			// then initialize a new rate limiter which allows an average of 2 requests
			// per second, with a maximum of 4 requests in a single ‘burst’, and add
			// the IP address and limiter to the map.
			if _, found := clients[ip]; !found {
				// Create and add a new client struct to the map if it doesn't exist.
				// Use the requests-per-second and burst values from the config struct.
				clients[ip] = &client{
					limiter: rate.NewLimiter(
						rate.Limit(app.config.limiter.rps),
						app.config.limiter.burst,
					),
				}
			}

			// Update the last seen time for the client.
			clients[ip].lastSeen = time.Now()

			// Call the Allow() method on the rate limiter for the current IP address
			// to see if the request is permitted, and if it's not, unlock the mutex and
			// call the rateLimitExceededResponse() helper to return a 429 Too Many
			// Requests response. The code behind the Allow() method is protected by
			// a mutex and is safe for concurrent use.
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}

			// IMPORTANT: Unlock the mutex before calling the next handler in the chain.
			// DON'T use defer to unlock the mutex, as that would mean that the mutex
			// isn't unlocked until all the handlers downstream of this middleware
			// have also returned.
			mu.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}

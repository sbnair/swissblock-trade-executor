# Swissblock Trade Executor

This repository contains my personal approach to the exercise proposed by Swissblock. In summary: "Create a trade execution service that takes an order size and an order price for an asset pair and generates market orders when the liquidity in the order book allows it."

## My Approach
I've opted for creating a simple command-line solution that receives the required parameters, and tries to identify possible trades for as long as defined.
Thus, I've identified the following abstractions:
* The `BookReader`, which defines the interface to implmlement a routine to retrieve Book Items from any given order book. Aside of this, the `bookReaderBinance` is an implementation of this interface that uses the binance WebSocket API to retrieve the order book in streaming mode.
* The `database` layer, based on `gorm` (a popular ORM library in Go) defines the way to access a database via URL. For the moment only sqlite databases are supported, but it would be easibly extendable to other engines.
  On top of this I've defined two models:
  ```mermaid
  erDiagram
    Order ||--o{ Trades: has
  ```
  **Order:** Persists the Order as a composite of the order-related imput parameters
  **Trade:** An order may materialize in one or more trades, depending on the liquidity found in the order book.
  These are included in the `model` package.
* The `Trader` is a service that bears the logic to execute the necessary trades an order can fulfill, based on the items retrieved from a `BookReader`.
* The `Reporting` helper package provide the necessary functions to summarize the resulting trades of an `Order`.

For the entire exercise I've been working in small chunks of time, during different days. But overall I've spent about 6 to 8 hours.

## Dificuties found
* I wasn't certain about the scope of the test, since the text indicates I'm expected to build a service. I initially thought that the expectation was to build a web service capable of executing different trades in parallel, via API calls.
  But I thought that would be too complex for a 6 hour exercise. So I took the simple path and created a minor command line application to execute the trades, providing all the necessay abstractions that could be eventually used later on a web service.
* There is an ongoing issue in `trader/limit_trader.go`, line 40 (`return nil`). This is the point where the `MatchOrder` function is expected to return when the order is completely fullfilled. However, the execution of that return statement will not return. Rather, the routinge gets stuck in the `select` clause. I'm not sure it this may be related to the fact I'm using golang 1.18 and an ubuntu box. But so far I've been unable to identify the cause.

## The Part I enjoyed the most
Defining the `book_reader` was the most challenging part, since I had to decide between using a polling approach or a streaming (via Go channels) to retrieve items from the book order. In the end I opted for the later, as the natural approach when retrieving items asynchrnously.
Since the reader runs in a dedicated goroutine, I had to decide where to place the error handling logic. The code launching the reading loop should be the one informed about any possible error, so I had to include a second channel where the reader could provide errors to the caller.

# Next Steps
* The service as a "proper service". I mean it should be a real service capable of receiving multiple orders in parallel, and try to execute them based on the book readings. Some of the challenges of this approach:
  * Even if the order requests are processed in parallel, the accesses to the order book items should be sincronyzed, via mutexes or other similar mechanisms, to avoid two orders being executed against the same book item.
  * Orders should be dispatched attending certain priority rules. For instance, one rule could be that these orders with the highest limit price should be attended first. For that, I would use a `Priority Queue` (heap data structure) to hold all the orders, so we guarantee the most prioritaire ones are executed first.
* Add telemetry, so the proper metrics can be emited and logged in a metric server, for the sake of observability and alerting.
* Create a proper docker file to hold the binary, and create a docker image to be held in a repository
* Create the proper manifests to run this docker image in either a docker compose service or a Kubernetes cluster.


# Hel - putting the network through hell (work in progress)

> In Norse mythology, Hel is a being who presides over a realm of
> the same name, where she receives a portion of the dead.

Hel provides durability, resilience and fuzz testing for applications
with dependencies for external resources, retrieved over a network.  
This is done by permutating the expected responses by means of
providing data-specific (eg. JSON format instead of XML),
protocol-specific (eg. HTTP 400 Bad Request instead of HTTP 200 OK)
or network-specific (eg. TCP timeout, invalid SSL certificate).

Think of Hel as a Chaos Monkey for network dependencies.

## Use cases
The application A retrieves data from external APIs B and C.  
A thus has a network-based dependency on both B and C, thus any
network-related failure or outage on either B or C, will impact 
A too.  

In order to harden A and thus reduce the impact, the soft spots
must be identified and quantified. This is where Hel comes into
play, allowing the developer to replace B and C with calls to
Hel instead, triggering various kinds of unexpected responses
and failures.

## Test cases
Hel provides the following test cases,

  - Return invalid data
  - Return valid data but invalid metadata (Content-Type vs. Accept-header)
  - Network connection timeout violations

Hel is focused on, but not limited to the HTTP protocol.
Adding support for other protocols such as QUIC, Websockets
or the like should be trivial.

## Usage
Hel is written in Go, allowing it to be executed as a single binary
with no external dependencies.  
Hel is well suited for executing in a container such as Docker.

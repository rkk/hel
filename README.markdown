# Hel - going through hell for the network

> In Norse mythology, Hel is a being who presides over a realm of
> the same name, where she receives a portion of the dead.

Hel provides durability and resilience testing for applications with
dependencies for external resources, retrieved over a network.

## Use cases
The application A retrieves data from external APIs B and C.  
A has a network-based dependency on both B and C, thus any
network-related failure or outage on either B or C, will affect
A too.

As a programmer working on A, I want to ensure that any failure
regarding A or B is properly mitigated in A, minimizing the impact
on the users of A.

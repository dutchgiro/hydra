# Hydra Workers

## What are Hydra Workers?

Hydra Workers are services that Hydra use in order to delegate the balance calculation.

## What are Hydra Worker Libraries?

Hydra Workers are made using libraries that abstract them of the comunication with the Hydra Server and they allow to developers focus in create only the logic of theirs balancers. The libraries provide to developers of all information of the instances collected for the Hydra Probes and the arguments set in the application configuration from the Hydra Server.

The libraries can be replicated in any language that can speak ZeroMQ v3.2.

## How a Hydra Worker?

All Hydra Worker share the same input/output interface. All of them receive and issue an array of instances nested n times.

i.e `[ i1, i2, i3 ] or [[i1, i3], [i4], [i5, i2]] or [[[i1], [i2]], [[i3]], [[i4, i5, i6], [i7, i8]]]`

## How to make a Hydra Worker?

For explain it nothing better than an example, Let's look at the code for make the hydra-worker-pong resend the same array of instances it receives:
```GO
package main

import (
	"os"

	worker "github.com/innotech/hydra-worker-pong/vendors/github.com/innotech/hydra-worker-lib"
)

func main() {
	pongWorker := worker.NewWorker(os.Args)
	fn := func(instances []interface{}, args map[string]interface{}) []interface{} {
		return instances
	}
	pongWorker.Run(fn)
}
```

package main

import (
	"projectXBackend/api"
	"projectXBackend/hazelcast"
)

func main() {

	client := hazelcast.GetHazelcastClient()

	api.StartApp(client)
}

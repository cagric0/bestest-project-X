package main

import (
	"projectXBackend/api"
	"projectXBackend/hazelcast"
)

func main() {

	client := hazelcast.GetHazelcastClient()
	//ctx := context.Background()
	//testMap, _ := client.GetMap(ctx, "TestMap")
	//size, err := testMap.Size(ctx)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(size)
	//v, err := testMap.Get(ctx, "AAAAA")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(v)
	api.StartApp(client)
}

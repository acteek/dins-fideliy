package main

import "fmt"

func main() {
	//dinsEndpoint := "https://my.dins.ru"
	//
	//dinsApi := dins.NewDinsApi(dinsEndpoint)
	//
	//if _, err := dinsApi.GetUser("6ae11e1d81b202ead1733354dce71ba7"); err != nil {
	//	log.Panic(err)
	//}
	type Test struct {
		Id string
	}

	var ol = make(map[int32]Test)

	ol[93] = Test{Id:"eeeee"}
	fmt.Print(ol[92])
}

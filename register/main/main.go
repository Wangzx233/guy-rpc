package main

import (
	"log"
	"net/http"
)

func main() {

	//guy_rpc.HandleHTTP()
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Println(err)
	}
}

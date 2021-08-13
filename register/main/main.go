package main

import (
	"guy-rpc/register/register"
	"log"
	"net/http"
)

func main()  {

	register.HandleHTTP()
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Println(err)
	}
}
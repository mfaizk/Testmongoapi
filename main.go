package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mfaizk/mongoapi/router"
)


func main(){
	fmt.Println("MongoDB Api")

	fmt.Println("SERVER IS GETTING STARTED...")
    r:=router.Router()
	
    log.Fatal(http.ListenAndServe(":4000",r))
	fmt.Println("listening at port 5000")

}
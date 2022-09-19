package main

import (
	"fmt"

	cmp "github.com/sinakeshmiri/arun/internal/adapters/framework/right/compiler"
)

func main(){
a:=`package function

import (
	"net/http"
)

func Function(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("HELLO WORLD!"))
}`

	t,err:=cmp.NewAdapter()
	if err != nil{fmt.Println(err)}
	e,err:=t.Compile(a)
	if err != nil{fmt.Println(err)}
	fmt.Println(e)
}
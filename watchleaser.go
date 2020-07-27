package main

import (
	"fmt"
	"leaser/leaser"
)

func main() {

	leaser := leaser.LeaserConf{TimeOut: 5, SvrName:"test",
		SvrHost: []string{"http://192.168.0.102:2379"}}
	err := leaser.NewWatcher()
	if err != nil {
		fmt.Println(err)
	}
}


package main

import (
	"fmt"
	"github.com/strongjz/leveledup.api/api"
)

func main() {

	fmt.Println("Starting API")
	api := api.Api{}
	api.LevelUp()
}

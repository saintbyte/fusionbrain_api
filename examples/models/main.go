package main

import (
	fbrain "github.com/saintbyte/fusionbrain_api"
	"log"
)

func main() {
	fb := fbrain.NewFusionbrain()
	models, err := fb.GetModels()
	if err != nil {
		log.Fatal(err)
	}
	for _, model := range models {
		log.Println(model)
		log.Println("ID:", model.Id)
		log.Println("Name:", model.Name)
		log.Println("Version:", model.Version)
		log.Println("Type: ", model.Type)
		log.Println("--------------------")
	}
}

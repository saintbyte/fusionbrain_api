package main

import "github.com/saintbyte/fusionbrain_api"

func main() {
	fb := fusionbrain_api.NewFusionbrain()
	//ModelId := "4"
	fb.Generate("Кот с пулеметом и в сапогах", "", "NOSTYLE")
}

package main

import "github.com/saintbyte/fusionbrain_api"

func main() {
	fb := fusionbrain_api.NewFusionbrain()
	//ModelId := "4"
	fb.Generate("Человек смотрит в небо а там огромный НЛО. 4k, Cyberpank", "", "NOSTYLE")
}

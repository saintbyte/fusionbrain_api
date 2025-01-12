package main

import (
	"github.com/saintbyte/fusionbrain_api"
	"log/slog"
	"time"
)

func main() {
	fb := fusionbrain_api.NewFusionbrain()
	result, err := fb.Generate(
		"Человек смотрит в небо а там огромный НЛО. 4k, Cyberpank",
		"",
		"NOSTYLE",
	)
	if err != nil {
		slog.Error("Generate error:", err)
		return
	}
	if result.Censored {
		slog.Error("Generate Censored")
		return
	}
	if result.Status != fusionbrain_api.FusionbrainGenerateStatusINITIAL {
		slog.Error("Generate status:", result.Status)
		return
	}
	slog.Info(result.Status)
	slog.Info(result.Uuid)
	var status = fusionbrain_api.GenerateResponse{}
	var cnt = 0
	for {
		status, err = fb.СheckStatus(result.Uuid)
		if err != nil {
			slog.Error("CheckStatus error:", err)
			return
		}
		if status.Status == fusionbrain_api.FusionbrainGenerateFAIL {
			slog.Error("CheckStatus FAIL")
			return
		}
		if status.Status == fusionbrain_api.FusionbrainGenerateDONE {
			break
		}
		if cnt > 100 {
			slog.Error("Too long , more 100 status request")
			return
		}
		cnt++
		time.Sleep(1 * time.Second)
	}
	for _, image := range status.Images {
		slog.Info(image)
	}
}

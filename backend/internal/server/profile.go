package server

import (
	"encoding/json"
	"net/http"
)

// Profile はトップページのヒーロー部に表示する自己紹介情報。
type Profile struct {
	Name         string   `json:"name"`
	Alias        string   `json:"alias"`
	Label        string   `json:"label"`
	Affiliations []string `json:"affiliations"`
	AvatarURL    string   `json:"avatarUrl"`
}

// DB を持たない方針のため、内容はここに直接置く。
var defaultProfile = Profile{
	Name:  "Yusuke Inoue",
	Alias: "cyokozai | 猪口才",
	Label: "Platform Engineer · Graduate Student",
	Affiliations: []string{
		"千葉工業大学大学院 情報科学研究科 情報科学専攻",
		"CloudNative Days 実行委員",
	},
	AvatarURL: "https://github.com/cyokozai.png",
}

func profileHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(defaultProfile)
	})
}

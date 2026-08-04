package parser

import "testing"

func TestMeaningfulContentLength(t *testing.T) {
	tests := []struct {
		name    string
		readme  string
		wantMin int // expected length is >= wantMin
		wantMax int // expected length is <= wantMax (use a large number for "no upper bound")
	}{
		{
			name:    "only a video link",
			readme:  "# PasarGuard\n\n[![Watch the video](https://img.youtube.com/vi/abc/0.jpg)](https://youtu.be/abcdefg)\n",
			wantMin: 0,
			wantMax: 40, // "PasarGuard" + "Watch the video" ~ well below the 150 threshold
		},
		{
			name:    "bare url only",
			readme:  "https://youtu.be/abcdefg\n",
			wantMin: 0,
			wantMax: 5,
		},
		{
			name: "real description",
			readme: "# Cool Project\n\nThis project is a local proxy that reduces LLM input costs " +
				"by rendering bulky context as compact images, leveraging vision capabilities " +
				"to compress system prompts and heavy documentation into a single request.\n",
			wantMin: 150,
			wantMax: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MeaningfulContentLength(tt.readme)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("MeaningfulContentLength() = %d, want within [%d, %d]", got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestIsEnglishReadme(t *testing.T) {
	tests := []struct {
		name        string
		readme      string
		wantEnglish bool
	}{
		{
			name: "english readme",
			readme: "# MovieBox TUI\n\nSearch, browse, play, and download movies and series from a " +
				"keyboard first terminal interface with isolated search, details and image caches.\n",
			wantEnglish: true,
		},
		{
			name: "english readme with emoji and dashes",
			readme: "# Voleeo API — ⭐ Star the repo\n\nOne native, AI friendly client for HTTP, gRPC, " +
				"WebSocket and GraphQL — the local first client for the AI era.\n",
			wantEnglish: true,
		},
		{
			name: "chinese readme",
			readme: "# iOS 定位修改器\n\n用代理软件的 HTTPS 解密功能，把 Apple 地图定位骗到世界任何角落。" +
				"本仓库将其核心逻辑移植为 JavaScript，适配到五个代理软件，无需越狱也无需电脑。\n",
			wantEnglish: false,
		},
		{
			name: "bilingual readme dominated by chinese",
			readme: "# iOS Location Spoofer\n\nEnglish · 中文\n\n用代理软件的 HTTPS 解密功能，把 Apple " +
				"地图定位骗到世界任何角落。参考项目：本项目基于核心研究，原始项目是用 Go 写的独立应用，" +
				"通过自建 VPN 与 MITM 代理实现定位欺骗。\n",
			wantEnglish: false,
		},
		{
			name:        "cyrillic readme",
			readme:      "# Проєкт\n\nЦе локальний проксі, що зменшує витрати на LLM, рендерячи громіздкий контекст у компактні зображення.\n",
			wantEnglish: false,
		},
		{
			name:        "too few letters to judge",
			readme:      "# 名字\n",
			wantEnglish: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ratio := IsEnglishReadme(tt.readme)
			if got != tt.wantEnglish {
				t.Errorf("IsEnglishReadme() = %v (non-Latin ratio %.2f), want %v", got, ratio, tt.wantEnglish)
			}
		})
	}
}

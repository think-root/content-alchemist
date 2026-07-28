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

package ytdlp

import (
	"testing"
)

func TestParseVideoInfo_DeduplicatesByResolution(t *testing.T) {
	raw := &rawVideoInfo{
		Title:     "Test Video",
		Thumbnail: "https://img.youtube.com/thumb.jpg",
		Duration:  180,
		Formats: []rawFormat{
			// two 1080p video-only: 248 has bigger filesize → should win
			{FormatID: "137", Ext: "mp4", Height: 1080, VCodec: "avc1", ACodec: "none", Filesize: 100},
			{FormatID: "248", Ext: "webm", Height: 1080, VCodec: "vp9", ACodec: "none", Filesize: 200},
			// 360p merged
			{FormatID: "18", Ext: "mp4", Height: 360, VCodec: "avc1", ACodec: "mp4a", Filesize: 50},
			// audio-only
			{FormatID: "140", Ext: "m4a", VCodec: "none", ACodec: "mp4a", ABR: 128},
		},
	}

	info := parseVideoInfo(raw)

	if len(info.Formats) != 3 {
		t.Fatalf("expected 3 formats (1080p deduped to 1, 360p, audio), got %d", len(info.Formats))
	}

	// video formats come first, ordered by first-seen resolution
	f1080 := info.Formats[0]
	if f1080.FormatID != "248" {
		t.Errorf("expected winner 248 (bigger filesize), got %s", f1080.FormatID)
	}
	if !f1080.NeedsAudioMerge {
		t.Error("1080p video-only should be NeedsAudioMerge=true")
	}
	if f1080.Ext != "mkv" {
		t.Errorf("video-only ext should be mkv, got %s", f1080.Ext)
	}

	f360 := info.Formats[1]
	if f360.FormatID != "18" {
		t.Errorf("expected format 18, got %s", f360.FormatID)
	}
	if f360.NeedsAudioMerge || f360.IsAudioOnly {
		t.Error("format 18 should be plain video+audio")
	}

	faudio := info.Formats[2]
	if faudio.FormatID != "140" {
		t.Errorf("expected format 140, got %s", faudio.FormatID)
	}
	if !faudio.IsAudioOnly {
		t.Error("format 140 should be IsAudioOnly=true")
	}
}

func TestParseVideoInfo_PreservesUnknownHeightVariants(t *testing.T) {
	// TikTok-style payload: several progressive MP4s with no height reported.
	// They must not collapse into a single "unknown" entry.
	raw := &rawVideoInfo{
		Title: "A TikTok",
		Formats: []rawFormat{
			{FormatID: "download", Ext: "mp4", Height: 0, VCodec: "h264", ACodec: "aac", Filesize: 0},
			{FormatID: "download_addr", Ext: "mp4", Height: 0, VCodec: "h264", ACodec: "aac", Filesize: 0},
			{FormatID: "play_addr", Ext: "mp4", Height: 0, VCodec: "h264", ACodec: "aac", Filesize: 0},
			{FormatID: "mp3-128", Ext: "mp3", VCodec: "none", ACodec: "mp3", ABR: 128},
		},
	}

	info := parseVideoInfo(raw)

	if len(info.Formats) != 4 {
		t.Fatalf("expected 4 formats (3 video variants + 1 audio), got %d", len(info.Formats))
	}

	ids := map[string]bool{}
	for _, f := range info.Formats {
		ids[f.FormatID] = true
	}
	for _, want := range []string{"download", "download_addr", "play_addr", "mp3-128"} {
		if !ids[want] {
			t.Errorf("expected format %q to be preserved, got %v", want, ids)
		}
	}
}

func TestResolutionLabel_Video(t *testing.T) {
	got := resolutionLabel(rawFormat{Height: 720, VCodec: "avc1"}, false)
	if got != "720p" {
		t.Errorf("expected %q, got %q", "720p", got)
	}
}

func TestResolutionLabel_AudioWithBitrate(t *testing.T) {
	got := resolutionLabel(rawFormat{ABR: 128}, true)
	if got != "128kbps" {
		t.Errorf("expected %q, got %q", "128kbps", got)
	}
}

func TestResolutionLabel_AudioNoBitrate(t *testing.T) {
	got := resolutionLabel(rawFormat{}, true)
	if got != "audio" {
		t.Errorf("expected %q, got %q", "audio", got)
	}
}

func TestResolutionLabel_UnknownVideo(t *testing.T) {
	got := resolutionLabel(rawFormat{VCodec: "avc1"}, false)
	if got != "unknown" {
		t.Errorf("expected %q, got %q", "unknown", got)
	}
}

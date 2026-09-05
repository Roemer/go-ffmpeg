package goffmpeg

import (
	"reflect"
	"testing"
)

func TestEncodingTags(t *testing.T) {
	video := NewVideoCopyEncoding().SetTag("hvc1").GetParameters()
	if want := []string{"-c:v", "copy", "-tag:v", "hvc1"}; !reflect.DeepEqual(video, want) {
		t.Fatalf("video parameters = %v, want %v", video, want)
	}

	audio := NewAudioCopyEncoding().SetIndex(2).SetTag("ac-3").GetParameters()
	if want := []string{"-c:a:2", "copy", "-tag:a:2", "ac-3"}; !reflect.DeepEqual(audio, want) {
		t.Fatalf("audio parameters = %v, want %v", audio, want)
	}
}

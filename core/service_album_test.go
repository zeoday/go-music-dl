package core

import (
	"reflect"
	"testing"
)

func TestAlbumFactoriesAndSourceList(t *testing.T) {
	supported := []string{"netease", "qq", "kugou", "kuwo", "migu", "jamendo", "joox", "qianqian", "soda"}
	for _, source := range supported {
		if fn := GetAlbumSearchFunc(source); fn == nil {
			t.Fatalf("GetAlbumSearchFunc(%q) returned nil", source)
		}
		if fn := GetAlbumDetailFunc(source); fn == nil {
			t.Fatalf("GetAlbumDetailFunc(%q) returned nil", source)
		}
		if fn := GetParseAlbumFunc(source); fn == nil {
			t.Fatalf("GetParseAlbumFunc(%q) returned nil", source)
		}
	}

	if got := GetAlbumSourceNames(); !reflect.DeepEqual(got, supported) {
		t.Fatalf("GetAlbumSourceNames() = %v, want %v", got, supported)
	}
}

func TestGetOriginalLinkSupportsAlbums(t *testing.T) {
	tests := []struct {
		source string
		id     string
		want   string
	}{
		{source: "netease", id: "123", want: "https://music.163.com/#/album?id=123"},
		{source: "qq", id: "abc", want: "https://y.qq.com/n/ryqq/albumDetail/abc"},
		{source: "kugou", id: "456", want: "https://www.kugou.com/album/456.html"},
		{source: "kuwo", id: "789", want: "http://www.kuwo.cn/album_detail/789"},
		{source: "migu", id: "321", want: "https://music.migu.cn/v3/music/album/321"},
		{source: "jamendo", id: "654", want: "https://www.jamendo.com/album/654"},
		{source: "joox", id: "album-id", want: "https://www.joox.com/hk/album/album-id"},
		{source: "qianqian", id: "PS1000000001", want: "https://music.91q.com/album/PS1000000001"},
		{source: "soda", id: "852", want: "https://www.qishui.com/share/album?album_id=852"},
	}

	for _, tt := range tests {
		if got := GetOriginalLink(tt.source, tt.id, "album"); got != tt.want {
			t.Fatalf("GetOriginalLink(%q, %q, album) = %q, want %q", tt.source, tt.id, got, tt.want)
		}
	}
}

func TestDetectSourceSupportsAlbumCapableNewSources(t *testing.T) {
	tests := []struct {
		link string
		want string
	}{
		{link: "https://music.migu.cn/v3/music/album/123", want: "migu"},
		{link: "https://www.jamendo.com/album/456", want: "jamendo"},
		{link: "https://www.joox.com/hk/album/abc", want: "joox"},
		{link: "https://music.91q.com/album/PS0001", want: "qianqian"},
		{link: "https://www.qishui.com/share/album?album_id=777", want: "soda"},
	}

	for _, tt := range tests {
		if got := DetectSource(tt.link); got != tt.want {
			t.Fatalf("DetectSource(%q) = %q, want %q", tt.link, got, tt.want)
		}
	}
}

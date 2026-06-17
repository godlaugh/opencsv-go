package csvcore

import (
	"bytes"
	"testing"
)

func TestDetectEncoding(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want EncodingName
	}{
		{name: "utf8", data: []byte("id,name\n1,Ada\n"), want: EncodingUTF8},
		{name: "utf8 bom", data: append([]byte{0xEF, 0xBB, 0xBF}, []byte("id,name\n")...), want: EncodingUTF8BOM},
		{name: "utf16 le bom", data: []byte{0xFF, 0xFE, 'i', 0, 'd', 0, '\n', 0}, want: EncodingUTF16LE},
		{name: "utf16 be bom", data: []byte{0xFE, 0xFF, 0, 'i', 0, 'd', 0, '\n'}, want: EncodingUTF16BE},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectEncoding(tt.data)
			if got.Name != tt.want {
				t.Fatalf("encoding = %s, want %s", got.Name, tt.want)
			}
		})
	}
}

func TestDetectEncodingLegacyNamedEncodings(t *testing.T) {
	cases := []struct {
		name     string
		encoding string
		text     string
		want     EncodingName
	}{
		{
			name:     "shift jis with kana",
			encoding: "shift-jis",
			text:     "名前,メモ\n太郎,かな\n",
			want:     EncodingShiftJIS,
		},
		{
			name:     "gbk chinese",
			encoding: "gbk",
			text:     "姓名,值\n张三,中文\n",
			want:     EncodingGBK,
		},
		{
			name:     "windows 1252 euro",
			encoding: "windows-1252",
			text:     "name,note\n1,café €\n",
			want:     EncodingWindows1252,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			encoding, err := ParseEncodingName(tt.encoding)
			if err != nil {
				t.Fatal(err)
			}
			encoded, err := EncodeFromUTF8([]byte(tt.text), encoding)
			if err != nil {
				t.Fatal(err)
			}
			detected := DetectEncoding(encoded)
			if detected.Name != tt.want {
				t.Fatalf("encoding = %+v, want %s", detected, tt.want)
			}
			decoded, err := DecodeToUTF8(encoded, detected)
			if err != nil {
				t.Fatal(err)
			}
			if string(decoded) != tt.text {
				t.Fatalf("decoded = %q, want %q", decoded, tt.text)
			}
		})
	}
}

func TestUTF16RoundTrip(t *testing.T) {
	encoding := EncodingInfo{Name: EncodingUTF16LE, HasBOM: true}
	input := []byte("id,name\r\n1,Ada\r\n")
	encoded, err := EncodeFromUTF8(input, encoding)
	if err != nil {
		t.Fatal(err)
	}
	if got := DetectEncoding(encoded); got.Name != EncodingUTF16LE || !got.HasBOM {
		t.Fatalf("encoding = %+v, want utf-16le with BOM", got)
	}
	decoded, err := DecodeToUTF8(encoded, encoding)
	if err != nil {
		t.Fatal(err)
	}
	if string(decoded) != string(input) {
		t.Fatalf("decoded = %q, want %q", decoded, input)
	}
}

func TestLegacySingleByteRoundTripPreservesBytes(t *testing.T) {
	encoding := EncodingInfo{Name: EncodingLegacySingleByte}
	input := []byte{'i', 'd', ',', 'n', 'a', 'm', 'e', '\n', '1', ',', 0xE9, '\n'}
	decoded, err := DecodeToUTF8(input, encoding)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := EncodeFromUTF8(decoded, encoding)
	if err != nil {
		t.Fatal(err)
	}
	if string(encoded) != string(input) {
		t.Fatalf("encoded = %v, want %v", encoded, input)
	}
}

func TestNamedEncodingRoundTrip(t *testing.T) {
	cases := []struct {
		name      string
		encoding  string
		input     string
		wantBytes []byte
	}{
		{
			name:      "shift jis",
			encoding:  "shift-jis",
			input:     "名前\n太郎\n",
			wantBytes: []byte{0x96, 0xbc, 0x91, 0x4f, '\n', 0x91, 0xbe, 0x98, 0x59, '\n'},
		},
		{
			name:      "gbk",
			encoding:  "gbk",
			input:     "姓名\n张三\n",
			wantBytes: []byte{0xd0, 0xd5, 0xc3, 0xfb, '\n', 0xd5, 0xc5, 0xc8, 0xfd, '\n'},
		},
		{
			name:      "windows 1252",
			encoding:  "windows-1252",
			input:     "name\ncafe \u20ac\n",
			wantBytes: []byte{'n', 'a', 'm', 'e', '\n', 'c', 'a', 'f', 'e', ' ', 0x80, '\n'},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			encoding, err := ParseEncodingName(tt.encoding)
			if err != nil {
				t.Fatal(err)
			}
			encoded, err := EncodeFromUTF8([]byte(tt.input), encoding)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(encoded, tt.wantBytes) {
				t.Fatalf("encoded = % x, want % x", encoded, tt.wantBytes)
			}
			decoded, err := DecodeToUTF8(encoded, encoding)
			if err != nil {
				t.Fatal(err)
			}
			if string(decoded) != tt.input {
				t.Fatalf("decoded = %q, want %q", decoded, tt.input)
			}
		})
	}
}

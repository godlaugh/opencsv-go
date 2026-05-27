package csvcore

import "testing"

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

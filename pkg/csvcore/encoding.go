package csvcore

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type EncodingName string

const (
	EncodingUTF8             EncodingName = "utf-8"
	EncodingUTF8BOM          EncodingName = "utf-8-bom"
	EncodingUTF16LE          EncodingName = "utf-16le"
	EncodingUTF16BE          EncodingName = "utf-16be"
	EncodingShiftJIS         EncodingName = "shift-jis"
	EncodingGBK              EncodingName = "gbk"
	EncodingWindows1252      EncodingName = "windows-1252"
	EncodingLegacySingleByte EncodingName = "legacy-single-byte"
)

type EncodingInfo struct {
	Name       EncodingName `json:"name"`
	HasBOM     bool         `json:"hasBom"`
	Confidence float64      `json:"confidence"`
}

func ParseEncodingName(raw string) (EncodingInfo, error) {
	switch strings.ToLower(strings.ReplaceAll(strings.TrimSpace(raw), "_", "-")) {
	case "", "auto":
		return EncodingInfo{}, errors.New("encoding override is empty")
	case "utf8", "utf-8":
		return EncodingInfo{Name: EncodingUTF8, Confidence: 1}, nil
	case "utf8bom", "utf-8-bom", "utf8-bom":
		return EncodingInfo{Name: EncodingUTF8BOM, HasBOM: true, Confidence: 1}, nil
	case "utf16le", "utf-16le", "utf-16-le":
		return EncodingInfo{Name: EncodingUTF16LE, Confidence: 1}, nil
	case "utf16lebom", "utf-16le-bom", "utf-16-le-bom":
		return EncodingInfo{Name: EncodingUTF16LE, HasBOM: true, Confidence: 1}, nil
	case "utf16be", "utf-16be", "utf-16-be":
		return EncodingInfo{Name: EncodingUTF16BE, Confidence: 1}, nil
	case "utf16bebom", "utf-16be-bom", "utf-16-be-bom":
		return EncodingInfo{Name: EncodingUTF16BE, HasBOM: true, Confidence: 1}, nil
	case "shiftjis", "shift-jis", "shift-jisx0213", "shift-jisx-0213", "sjis", "cp932", "ms932":
		return EncodingInfo{Name: EncodingShiftJIS, Confidence: 1}, nil
	case "gbk", "gb2312", "cp936":
		return EncodingInfo{Name: EncodingGBK, Confidence: 1}, nil
	case "windows-1252", "windows1252", "cp1252":
		return EncodingInfo{Name: EncodingWindows1252, Confidence: 1}, nil
	case "latin1", "latin-1", "iso-8859-1":
		return EncodingInfo{Name: EncodingLegacySingleByte, Confidence: 1}, nil
	case "legacy-single-byte", "legacy", "single-byte":
		return EncodingInfo{Name: EncodingLegacySingleByte, Confidence: 1}, nil
	default:
		return EncodingInfo{}, fmt.Errorf("unsupported encoding %q", raw)
	}
}

func DetectEncoding(sample []byte) EncodingInfo {
	switch {
	case bytes.HasPrefix(sample, []byte{0xEF, 0xBB, 0xBF}):
		return EncodingInfo{Name: EncodingUTF8BOM, HasBOM: true, Confidence: 1}
	case bytes.HasPrefix(sample, []byte{0xFF, 0xFE}):
		return EncodingInfo{Name: EncodingUTF16LE, HasBOM: true, Confidence: 1}
	case bytes.HasPrefix(sample, []byte{0xFE, 0xFF}):
		return EncodingInfo{Name: EncodingUTF16BE, HasBOM: true, Confidence: 1}
	}

	if utf8.Valid(sample) {
		return EncodingInfo{Name: EncodingUTF8, Confidence: 0.98}
	}

	leZeros, beZeros, pairs := 0, 0, len(sample)/2
	for i := 0; i+1 < len(sample); i += 2 {
		if sample[i] == 0 && sample[i+1] != 0 {
			beZeros++
		}
		if sample[i] != 0 && sample[i+1] == 0 {
			leZeros++
		}
	}
	if pairs > 8 {
		if float64(leZeros)/float64(pairs) > 0.35 {
			return EncodingInfo{Name: EncodingUTF16LE, Confidence: 0.85}
		}
		if float64(beZeros)/float64(pairs) > 0.35 {
			return EncodingInfo{Name: EncodingUTF16BE, Confidence: 0.85}
		}
	}

	return EncodingInfo{Name: EncodingLegacySingleByte, Confidence: 0.45}
}

func DecodeToUTF8(input []byte, encoding EncodingInfo) ([]byte, error) {
	switch encoding.Name {
	case EncodingUTF8:
		if utf8.Valid(input) {
			return input, nil
		}
		return bytes.ToValidUTF8(input, []byte("?")), nil
	case EncodingUTF8BOM:
		input = bytes.TrimPrefix(input, []byte{0xEF, 0xBB, 0xBF})
		if utf8.Valid(input) {
			return input, nil
		}
		return bytes.ToValidUTF8(input, []byte("?")), nil
	case EncodingUTF16LE:
		return decodeUTF16(input, binary.LittleEndian, encoding.HasBOM)
	case EncodingUTF16BE:
		return decodeUTF16(input, binary.BigEndian, encoding.HasBOM)
	case EncodingShiftJIS, EncodingGBK, EncodingWindows1252:
		codec, ok := textEncoding(encoding.Name)
		if !ok {
			return nil, errors.New("unsupported encoding")
		}
		return transformBytes(input, codec.NewDecoder())
	case EncodingLegacySingleByte:
		runes := make([]rune, len(input))
		for i, b := range input {
			runes[i] = rune(b)
		}
		return []byte(string(runes)), nil
	default:
		return nil, errors.New("unsupported encoding")
	}
}

func EncodeFromUTF8(input []byte, encoding EncodingInfo) ([]byte, error) {
	switch encoding.Name {
	case EncodingUTF8:
		return input, nil
	case EncodingUTF8BOM:
		return append([]byte{0xEF, 0xBB, 0xBF}, input...), nil
	case EncodingUTF16LE:
		return encodeUTF16(string(input), binary.LittleEndian, encoding.HasBOM), nil
	case EncodingUTF16BE:
		return encodeUTF16(string(input), binary.BigEndian, encoding.HasBOM), nil
	case EncodingShiftJIS, EncodingGBK, EncodingWindows1252:
		codec, ok := textEncoding(encoding.Name)
		if !ok {
			return nil, errors.New("unsupported encoding")
		}
		return transformBytes(input, codec.NewEncoder())
	case EncodingLegacySingleByte:
		output := make([]byte, 0, len(input))
		for _, r := range string(input) {
			if r > 255 {
				return nil, errors.New("legacy single-byte encoding cannot represent one or more runes")
			}
			output = append(output, byte(r))
		}
		return output, nil
	default:
		return nil, errors.New("unsupported encoding")
	}
}

func decodeUTF16(input []byte, order binary.ByteOrder, hasBOM bool) ([]byte, error) {
	if hasBOM && len(input) >= 2 {
		input = input[2:]
	}
	if len(input)%2 == 1 {
		input = input[:len(input)-1]
	}

	units := make([]uint16, 0, len(input)/2)
	for i := 0; i+1 < len(input); i += 2 {
		units = append(units, order.Uint16(input[i:i+2]))
	}
	return []byte(string(utf16.Decode(units))), nil
}

func encodeUTF16(input string, order binary.ByteOrder, hasBOM bool) []byte {
	units := utf16.Encode([]rune(input))
	output := make([]byte, 0, len(units)*2+2)
	if hasBOM {
		if order == binary.LittleEndian {
			output = append(output, 0xFF, 0xFE)
		} else {
			output = append(output, 0xFE, 0xFF)
		}
	}
	var buf [2]byte
	for _, unit := range units {
		order.PutUint16(buf[:], unit)
		output = append(output, buf[:]...)
	}
	return output
}

func textEncoding(name EncodingName) (encoding.Encoding, bool) {
	switch name {
	case EncodingShiftJIS:
		return japanese.ShiftJIS, true
	case EncodingGBK:
		return simplifiedchinese.GBK, true
	case EncodingWindows1252:
		return charmap.Windows1252, true
	default:
		return nil, false
	}
}

func transformBytes(input []byte, transformer transform.Transformer) ([]byte, error) {
	return io.ReadAll(transform.NewReader(bytes.NewReader(input), transformer))
}

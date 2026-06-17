package csvcore

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
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

	return detectLegacyEncoding(sample)
}

func detectLegacyEncoding(sample []byte) EncodingInfo {
	type candidate struct {
		name  EncodingName
		score float64
	}

	legacyDecoded, _ := DecodeToUTF8(sample, EncodingInfo{Name: EncodingLegacySingleByte})
	shiftJISDecoded, shiftJISErr := DecodeToUTF8(sample, EncodingInfo{Name: EncodingShiftJIS})
	gbkDecoded, gbkErr := DecodeToUTF8(sample, EncodingInfo{Name: EncodingGBK})
	windows1252Decoded, windows1252Err := DecodeToUTF8(sample, EncodingInfo{Name: EncodingWindows1252})
	shiftJISStats := collectDecodedTextStats(shiftJISDecoded)
	gbkStats := collectDecodedTextStats(gbkDecoded)
	windows1252Stats := collectDecodedTextStats(windows1252Decoded)
	if shiftJISErr == nil && shiftJISStats.replacements == 0 && shiftJISStats.kana > 0 {
		return EncodingInfo{Name: EncodingShiftJIS, Confidence: 0.78}
	}
	if gbkErr == nil && gbkStats.replacements == 0 && gbkStats.cjk >= 2 && (shiftJISStats.replacements > 0 || shiftJISStats.halfwidthKana > 0) {
		return EncodingInfo{Name: EncodingGBK, Confidence: 0.74}
	}
	if windows1252Err == nil && windows1252Stats.replacements == 0 && windows1252Stats.westernSymbols > 0 && (shiftJISStats.replacements > 0 || gbkStats.replacements > 0) {
		return EncodingInfo{Name: EncodingWindows1252, Confidence: 0.74}
	}

	best := candidate{name: EncodingLegacySingleByte, score: decodedTextScore(legacyDecoded, EncodingLegacySingleByte, sample)}
	for _, name := range []EncodingName{EncodingShiftJIS, EncodingGBK, EncodingWindows1252} {
		decoded, err := DecodeToUTF8(sample, EncodingInfo{Name: name})
		if err != nil {
			continue
		}
		score := decodedTextScore(decoded, name, sample)
		if score > best.score {
			best = candidate{name: name, score: score}
		}
	}
	if best.name == EncodingLegacySingleByte || best.score < decodedTextScore(legacyDecoded, EncodingLegacySingleByte, sample)+8 {
		return EncodingInfo{Name: EncodingLegacySingleByte, Confidence: 0.45}
	}
	confidence := 0.62
	if best.score >= 40 {
		confidence = 0.78
	}
	return EncodingInfo{Name: best.name, Confidence: confidence}
}

type decodedTextStats struct {
	cjk            int
	kana           int
	halfwidthKana  int
	latinExtended  int
	replacements   int
	controls       int
	printableASCII int
	westernSymbols int
}

func collectDecodedTextStats(decoded []byte) decodedTextStats {
	text := string(decoded)
	stats := decodedTextStats{}
	for _, r := range text {
		switch {
		case r == unicode.ReplacementChar:
			stats.replacements++
		case r == '\n' || r == '\r' || r == '\t':
		case unicode.IsControl(r):
			stats.controls++
		case r >= 0x20 && r <= 0x7e:
			stats.printableASCII++
		case isKana(r):
			stats.kana++
		case isHalfwidthKana(r):
			stats.halfwidthKana++
		case isCJK(r):
			stats.cjk++
		case isWesternEncodingSymbol(r):
			stats.westernSymbols++
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			stats.latinExtended++
		}
	}
	return stats
}

func decodedTextScore(decoded []byte, name EncodingName, raw []byte) float64 {
	stats := collectDecodedTextStats(decoded)
	score := 0.0
	score += float64(stats.printableASCII) * 0.3
	score += float64(stats.latinExtended) * 2
	score += float64(stats.cjk) * 3
	score += float64(stats.kana) * 5
	score += float64(stats.westernSymbols) * 4
	score -= float64(stats.halfwidthKana) * 4
	score -= float64(stats.replacements) * 25
	score -= float64(stats.controls) * 18

	switch name {
	case EncodingShiftJIS:
		score += float64(stats.kana) * 18
		score += float64(stats.cjk) * 2
	case EncodingGBK:
		score += float64(stats.cjk) * 8
		score -= float64(stats.kana+stats.halfwidthKana) * 15
	case EncodingWindows1252:
		score += float64(countC1Bytes(raw)) * 2
		score += float64(stats.latinExtended) * 4
		score += float64(stats.westernSymbols) * 16
		score -= float64(stats.cjk+stats.kana+stats.halfwidthKana) * 10
	case EncodingLegacySingleByte:
		score -= float64(countC1Bytes(raw)) * 10
	}
	return score
}

func countC1Bytes(input []byte) int {
	count := 0
	for _, b := range input {
		if b >= 0x80 && b <= 0x9f {
			count++
		}
	}
	return count
}

func isKana(r rune) bool {
	return (r >= 0x3040 && r <= 0x30ff) || (r >= 0x31f0 && r <= 0x31ff)
}

func isHalfwidthKana(r rune) bool {
	return r >= 0xff61 && r <= 0xff9f
}

func isCJK(r rune) bool {
	return (r >= 0x3400 && r <= 0x4dbf) || (r >= 0x4e00 && r <= 0x9fff)
}

func isWesternEncodingSymbol(r rune) bool {
	switch r {
	case '€', '‚', 'ƒ', '„', '…', '†', '‡', 'ˆ', '‰', 'Š', '‹', 'Œ', 'Ž',
		'‘', '’', '“', '”', '•', '–', '—', '˜', '™', 'š', '›', 'œ', 'ž', 'Ÿ':
		return true
	default:
		return false
	}
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

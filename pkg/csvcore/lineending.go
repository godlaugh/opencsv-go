package csvcore

import "bytes"

type LineEndingName string

const (
	LineEndingLF   LineEndingName = "lf"
	LineEndingCRLF LineEndingName = "crlf"
	LineEndingCR   LineEndingName = "cr"
)

type LineEndingInfo struct {
	Name       LineEndingName `json:"name"`
	Sequence   string         `json:"sequence"`
	Confidence float64        `json:"confidence"`
}

func DetectLineEnding(sample []byte) LineEndingInfo {
	crlf := bytes.Count(sample, []byte("\r\n"))
	withoutCRLF := bytes.ReplaceAll(sample, []byte("\r\n"), nil)
	lf := bytes.Count(withoutCRLF, []byte("\n"))
	cr := bytes.Count(withoutCRLF, []byte("\r"))

	total := crlf + lf + cr
	if total == 0 {
		return LineEndingInfo{Name: LineEndingLF, Sequence: "\n", Confidence: 0.5}
	}

	if crlf >= lf && crlf >= cr {
		return LineEndingInfo{Name: LineEndingCRLF, Sequence: "\r\n", Confidence: float64(crlf) / float64(total)}
	}
	if cr >= lf {
		return LineEndingInfo{Name: LineEndingCR, Sequence: "\r", Confidence: float64(cr) / float64(total)}
	}
	return LineEndingInfo{Name: LineEndingLF, Sequence: "\n", Confidence: float64(lf) / float64(total)}
}

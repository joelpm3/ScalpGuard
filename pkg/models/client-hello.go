package models

import (
	"fmt"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/gospider007/ja3"
	"github.com/gospider007/requests"
	"github.com/joelpm3/ScalpGuard/pkg/algorithms"
)

// Assuming a generic structure for an HTTP/2 setting for demonstration purposes.
type Setting struct {
	ID  uint16
	Val uint32
}

type ParsedClientHello struct {
	NegotiatedProtocol    string
	TlsVersion            uint16
	UserAgent             string
	OrderHeaders          []string
	Cookies               string
	Tls                   ja3.TlsData
	Ja3                   string
	Ja3n                  string
	Ja3H                  string
	Ja3nH                 string
	Http2Fingerprint      string
	Http2FingerprintHash  string
}

func ParseClientHello(ctx *gin.Context) (ParsedClientHello, error) {
	fpData, ok := ja3.GetFpContextData(ctx.Request.Context())
	if !ok {
		return ParsedClientHello{}, fmt.Errorf("unable to fingerprint TLS handshake")
	}

	h2Ja3Spec := fpData.H2Ja3Spec()
	connectionState := fpData.ConnectionState() // Corrected to match the expected return value

	initialSettingsStr := buildInitialSettingsStr(h2Ja3Spec.InitialSetting)
	orderHeadersStr := abbreviateOrderHeaders(h2Ja3Spec.OrderHeaders)

	http2Fingerprint := fmt.Sprintf("%s|0|%s", initialSettingsStr, orderHeadersStr)

	tlsData, err := fpData.TlsData()
	if err != nil {
		return ParsedClientHello{}, err // Adjust based on the actual method signature
	}

	result := ParsedClientHello{
		NegotiatedProtocol:    connectionState.NegotiatedProtocol,
		TlsVersion:            connectionState.Version,
		UserAgent:             ctx.Request.UserAgent(),
		OrderHeaders:          fpData.OrderHeaders(),
		Cookies:               requests.Cookies(ctx.Request.Cookies()).String(),
		Tls:                   tlsData,
		Http2Fingerprint:      http2Fingerprint,
		Http2FingerprintHash:  algorithms.Ja3Digest(http2Fingerprint),
	}

	result.Ja3, result.Ja3n = tlsData.Fp()
	result.Ja3H = algorithms.Ja3Digest(result.Ja3)
	result.Ja3nH = algorithms.Ja3Digest(result.Ja3n)

	return result, nil
}

func buildInitialSettingsStr(settings []ja3.Setting) string {
	var sb strings.Builder
	for i, setting := range settings {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("%d:%d", setting.Id, setting.Val))
	}
	return sb.String()
}

func abbreviateOrderHeaders(headers []string) string {
	headerAbbr := map[string]string{
		":method":   "m",
		":authority": "a",
		":scheme":    "s",
		":path":      "p",
	}
	var sb strings.Builder
	for _, header := range headers {
		if abbr, ok := headerAbbr[header]; ok {
			if sb.Len() > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(abbr)
		}
	}
	return sb.String()
}

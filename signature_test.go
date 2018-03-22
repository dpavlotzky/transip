package transip

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"testing"
	"encoding/pem"
	"crypto/x509"
)

func TestParamEncode(t *testing.T) {
	in := []kV{
		kV{"__method", "getDomainNames"},
		kV{"__service", "DomainService"},
		kV{"__hostname", "api.transip.nl"},
		kV{"__timestamp", "1492760973"},
		kV{"__nonce", "58f9b98ddd3999.86051758"},
	}
	expect := []byte(`__method=getDomainNames&__service=DomainService&__hostname=api.transip.nl&__timestamp=1492760973&__nonce=58f9b98ddd3999.86051758`)

	if out := urlencode(in); bytes.Compare(out, expect) != 0 {
		t.Errorf("Mismatch out=%s\nexpect=%s", out, expect)
	}
}

func TestSHA512ASN1(t *testing.T) {
	in := []byte(`__method=getDomainNames&__service=DomainService&__hostname=api.transip.nl&__timestamp=1492851509&__nonce=58fb1b35916f25.33598874`)
	expect := []byte{
		0x30, 0x51, 0x30, 0x0d, 0x06, 0x09, 0x60, 0x86, 0x48, 0x01, 0x65, 0x03, 0x04, 0x02, 0x03, 0x05, 0x00, 0x04,
		0x40, 0x60, 0x87, 0x8e, 0x6a, 0x11, 0x4a, 0x92, 0xc7, 0x62, 0x24, 0x8a, 0xb3, 0xa1, 0xc0, 0x28, 0xd6, 0xfb,
		0x80, 0xe0, 0x00, 0x8c, 0xce, 0xd4, 0x22, 0x97, 0x31, 0x89, 0xb8, 0x12, 0x38, 0xbb, 0xeb, 0xb1, 0x7d, 0x8e,
		0x69, 0xd9, 0x9b, 0x29, 0x62, 0x22, 0x50, 0xb0, 0x7f, 0x4c, 0x3d, 0x5f, 0xcc, 0x6b, 0x9d, 0x82, 0xf8, 0xd4,
		0xe4, 0xba, 0xfb, 0x0e, 0xe1, 0x42, 0x6d, 0xdf, 0xe6, 0xef, 0x12,
	}

	out := sha512ASN1(in)
	if bytes.Compare(out, expect) != 0 {
		t.Errorf("Mismatch out=%x\nexpect=%x", out, expect)
	}
}

func TestBase64(t *testing.T) {
	expect := `Uy1zdZC/AIpCamHbxJD9vzCWobocTG1/8aHn3wcDcl8h4r+i4GGGG9XTw1bam7HJQdEZzc+TAPXt55ldRjfdLlR0pC9FAn6FsvRk6OpxvrjlsbObUUYM/XtzpMCCupt/FskV3wsoCAXl7Yj27LxOU2duK0uJq/0aLVYwe8C65FE12egOgB7ghfvkN3DHnApmYYvnm6T5qQ7gUX/i/UYHJJueTfpup1CXCpFsebVY+8F5eSuLRVVMHCP49MwP4y+SKU25vmRaUSDUmE4ZdEJzvUYdBUEmZ8KDgGj3gu/qBZH9CRRcx2o4z4jjMS4XLxlz4mwCztjs9wnu2PCvscM01Q==`
	in := []byte{
		0x53,
		0x2d, 0x73, 0x75, 0x90, 0xbf, 0x00, 0x8a, 0x42, 0x6a, 0x61, 0xdb, 0xc4, 0x90, 0xfd, 0xbf, 0x30, 0x96, 0xa1,
		0xba, 0x1c, 0x4c, 0x6d, 0x7f, 0xf1, 0xa1, 0xe7, 0xdf, 0x07, 0x03, 0x72, 0x5f, 0x21, 0xe2, 0xbf, 0xa2, 0xe0,
		0x61, 0x86, 0x1b, 0xd5, 0xd3, 0xc3, 0x56, 0xda, 0x9b, 0xb1, 0xc9, 0x41, 0xd1, 0x19, 0xcd, 0xcf, 0x93, 0x00,
		0xf5, 0xed, 0xe7, 0x99, 0x5d, 0x46, 0x37, 0xdd, 0x2e, 0x54, 0x74, 0xa4, 0x2f, 0x45, 0x02, 0x7e, 0x85, 0xb2,
		0xf4, 0x64, 0xe8, 0xea, 0x71, 0xbe, 0xb8, 0xe5, 0xb1, 0xb3, 0x9b, 0x51, 0x46, 0x0c, 0xfd, 0x7b, 0x73, 0xa4,
		0xc0, 0x82, 0xba, 0x9b, 0x7f, 0x16, 0xc9, 0x15, 0xdf, 0x0b, 0x28, 0x08, 0x05, 0xe5, 0xed, 0x88, 0xf6, 0xec,
		0xbc, 0x4e, 0x53, 0x67, 0x6e, 0x2b, 0x4b, 0x89, 0xab, 0xfd, 0x1a, 0x2d, 0x56, 0x30, 0x7b, 0xc0, 0xba, 0xe4,
		0x51, 0x35, 0xd9, 0xe8, 0x0e, 0x80, 0x1e, 0xe0, 0x85, 0xfb, 0xe4, 0x37, 0x70, 0xc7, 0x9c, 0x0a, 0x66, 0x61,
		0x8b, 0xe7, 0x9b, 0xa4, 0xf9, 0xa9, 0x0e, 0xe0, 0x51, 0x7f, 0xe2, 0xfd, 0x46, 0x07, 0x24, 0x9b, 0x9e, 0x4d,
		0xfa, 0x6e, 0xa7, 0x50, 0x97, 0x0a, 0x91, 0x6c, 0x79, 0xb5, 0x58, 0xfb, 0xc1, 0x79, 0x79, 0x2b, 0x8b, 0x45,
		0x55, 0x4c, 0x1c, 0x23, 0xf8, 0xf4, 0xcc, 0x0f, 0xe3, 0x2f, 0x92, 0x29, 0x4d, 0xb9, 0xbe, 0x64, 0x5a, 0x51,
		0x20, 0xd4, 0x98, 0x4e, 0x19, 0x74, 0x42, 0x73, 0xbd, 0x46, 0x1d, 0x05, 0x41, 0x26, 0x67, 0xc2, 0x83, 0x80,
		0x68, 0xf7, 0x82, 0xef, 0xea, 0x05, 0x91, 0xfd, 0x09, 0x14, 0x5c, 0xc7, 0x6a, 0x38, 0xcf, 0x88, 0xe3, 0x31,
		0x2e, 0x17, 0x2f, 0x19, 0x73, 0xe2, 0x6c, 0x02, 0xce, 0xd8, 0xec, 0xf7, 0x09, 0xee, 0xd8, 0xf0, 0xaf, 0xb1,
		0xc3, 0x34, 0xd5,
	}
	out := base64.StdEncoding.EncodeToString(in)
	if out != expect {
		t.Errorf("Mismatch sig=%s\nexpect=%s", out, expect)
	}
}

func TestSign(t *testing.T) {
	expect := `3xhFhDsp1H2%2Ba4LidS2hXuQQjNmsxIGlEkgastg9jO5BRvqTLUXppAXjhieoq4P%2Bf8E5BF9%2FZ4SnYWKkQUvMbG%2BfoMnmK%2BGL6CwYeyn%2FLglZMNrFdoMw18PRH1iW3quvF2yxVcnCJXT%2FRwXUQu4TbH7f8kYH12R20hNI6HlDDx%2FQbLPtMyS9nMAmhhebtmjnutlmZS%2BCK%2Bh9jllaGcMiCorYGBD5aXHu%2FUdlsxmzSQ8acvtDa%2BfBkZ%2BDcvi%2B9fafzMCxd1OT5%2BP9vj93u1i9VmL0Bz3Z88Uj%2FxCnTXH0VdZZGyWuqKNNjawXs1wwoJ5%2FrDhMCyVdqt%2Fb8X2%2Bt%2F%2Bpww%3D%3D`
	sig, e := Sign(privKey(), []KV{
			{"__method", "getDomainNames"},
			{"__service", "DomainService"},
			{"__hostname", "api.transip.nl"},
			{"__timestamp", "1492851509"},
			{"__nonce", "58fb1b35916f25.33598874"},
		})
	if e != nil {
		t.Fatal(e)
	}
	if sig != expect {
		t.Errorf("Got:\n%s\n\nExpected:\n%s", sig, expect)
	}
}

func privKey() *rsa.PrivateKey {
	keyContents := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA7xQerL6xF0GGYxFYXHTGwX59jyDyCIR6QVDoIMcEAHg1tGj8
tWT3ydjdoWJLnwpzjN4KXCIxMu3ybEqHo/YDBHMGfvqwFulFQCxwA3Yf8MLG4Ig4
q6njNb9DHBTAZwhcVxtk3PjvCk7PuKetqb9+bXANh9gUglH8CPQcxRKF9gdvMR3Q
I9O7s+VJmsKZ++QQsiuMw5cHBP7/1WoixBnDJ/J37Q5xuuPYPQ7BmbGZZnkiH99R
XV5yTXRRYwvofQmJJH4HkMAsS+Vim4n43AuZijUKuHs6ctzLvckiTWwyKt4GDMDM
fk6R9t5xYJn83NOF+CbuiXleXtViiQz+ZQ5RCwIDAQABAoIBAQDjTnfTugJZoB0L
d+RRE14dfgwW1zYHTx2FmEz7TPzLDX/SJbePJ45HxP8Df5dygNdX6YxkCMZKK92/
hCTuiOpZgpt8gxCE2AjVeOqO//JiUG4R8LIg1IeIBG7j9f7wdwyEbTE6vxtW65On
dxUwPTcRCeZzb8ggF57PTHlGDdR8EsLiHyEVlb8+hiXbSowq5Ydr+cv38c1Pz10y
lwroXmtn7n8+wJCjVxxABDCyEQI7/H00U6eVetc1LAwJQtQJFvAwGNCymMcsmXq7
L4WTrSlxRSIgWwJtk8rsnaY2UVNuz+4TWfZ5HrIu2LED8Ml7meHmV/bji7uYBp8+
z8Y8OfHhAoGBAPz8xZa/SJECy5DCpFCoo2MQgPpxZTJCJ2+BEaWIjlEvXIlL+Hmk
8VU9Jea0e6UVKc7EVqqbHlZUx1N5q9TaWOrQfZG29tFaUJLX+wDKseccSJFd6+hU
yoRjxwa533kQfJpNiSYjg32apvXBUVyrHlCyKcJWuHxdx5hLH/5NlNhTAoGBAPHs
8n43Ga28LfFFw5bx8LFfTQvPZ+sWljOol+bU1+vq9EABTJkbJlJLALS+LFCKePCR
1hatd8r8GvRoTKeXPVqIpu7fcSWcHSJPFwfln/KT2/GnBbuBi5V8NHdX/rOpr/Jo
5dm7Ie+VOPiCbS1gcL8L64zlzp3vhrKYg4f3uS1pAoGAPCtWRzM5bBvRFJ3mfLSP
H4mWU3pSyjBHttJowwkGaDKufI0QDMZ5C3/ems9ENRAigGXcAvmfroK9YZInlxlT
Wo25v8VXUJV9Yl9x+E89Hq1waPqAmCJKhFBCzsu4Zc/RAtX8D5EUvfPhT8PpuPON
4z1shyce+51GUmdTtaT4CLcCgYEA6csYlUzebf1bUL0gxXDOMDtvE6i+Pnw3b3jQ
Q20RtZX7sRcQVS3dnM2KwyC9ZqBLPAFTqdq919ZGnkdlPNh1nFZPLK2WhMgXh55z
HViVeDHX7fKBIbGRmUbM1UCejjXAKT2iUwX7R7MnuVEh/SYkDxyP6Dv1rr9ZpqKp
Ce8mb3ECgYAQarY6LDm/Y4CHNSsqWgngev02hpOXFXEekBxjjcUy5CnuENf6EJ2G
E+Hzb9AqohkZQWPpTqsuYOhdsS6N7udF/gnX0vsqBqu3e/Z9lV3FCbQk5sqKqu0G
u7Isy/Q8xjHsJG6KP2/pMvMG42lhB3b+GKSmxVM999U30SdgorE3Kw==
-----END RSA PRIVATE KEY-----`)

	block, _ := pem.Decode(keyContents)
	privKey,_ := x509.ParsePKCS1PrivateKey(block.Bytes)
	return privKey
}

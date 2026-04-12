package aws

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

// Sign signs an AWS request using Signature Version 4.
func Sign(req *http.Request, body []byte, accessKey, secretKey, region, service string, now time.Time) {
	date := now.UTC().Format("20060102")
	timestamp := now.UTC().Format("20060102T150405Z")

	req.Header.Set("X-Amz-Date", timestamp)

	// 1. Create canonical request
	canonicalRequest := createCanonicalRequest(req, body)
	hashedCanonicalRequest := sha256Hash(canonicalRequest)

	// 2. Create string to sign
	credentialScope := fmt.Sprintf("%s/%s/%s/aws4_request", date, region, service)
	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s",
		timestamp,
		credentialScope,
		hashedCanonicalRequest,
	)

	// 3. Calculate signature
	signingKey := getSigningKey(secretKey, date, region, service)
	signature := hex.EncodeToString(hmacSHA256(signingKey, stringToSign))

	// 4. Create authorization header
	authHeader := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		accessKey,
		credentialScope,
		getSignedHeaders(req),
		signature,
	)
	req.Header.Set("Authorization", authHeader)
}

func createCanonicalRequest(req *http.Request, body []byte) string {
	// Method
	method := req.Method

	// URI
	uri := req.URL.Path
	if uri == "" {
		uri = "/"
	}

	// Query string
	query := req.URL.Query()
	keys := make([]string, 0, len(query))
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var queryParts []string
	for _, k := range keys {
		for _, v := range query[k] {
			queryParts = append(queryParts, fmt.Sprintf("%s=%s", k, v))
		}
	}
	canonicalQuery := strings.Join(queryParts, "&")

	// Headers
	var signedHeaders []string
	var canonicalHeaders []string
	headerKeys := make([]string, 0, len(req.Header))
	for k := range req.Header {
		headerKeys = append(headerKeys, strings.ToLower(k))
	}
	sort.Strings(headerKeys)
	for _, k := range headerKeys {
		v := req.Header.Get(k)
		canonicalHeaders = append(canonicalHeaders, fmt.Sprintf("%s:%s", k, strings.TrimSpace(v)))
		signedHeaders = append(signedHeaders, k)
	}

	// Payload hash
	payloadHash := sha256Hash(string(body))

	return fmt.Sprintf("%s\n%s\n%s\n%s\n\n%s\n%s",
		method,
		uri,
		canonicalQuery,
		strings.Join(canonicalHeaders, "\n")+"\n",
		strings.Join(signedHeaders, ";"),
		payloadHash,
	)
}

func getSignedHeaders(req *http.Request) string {
	var keys []string
	for k := range req.Header {
		keys = append(keys, strings.ToLower(k))
	}
	sort.Strings(keys)
	return strings.Join(keys, ";")
}

func getSigningKey(secret, date, region, service string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+secret), date)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, service)
	kSigning := hmacSHA256(kService, "aws4_request")
	return kSigning
}

func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

func sha256Hash(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

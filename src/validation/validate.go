package validation

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"strings"
)

// IsValidPayload checks if the github payload's hash fits with the hash computed by GitHub sent as a header
func IsValidPayload(secret, headerHash string, payload []byte) bool {
	hash := HashPayload(secret, payload)
	return hmac.Equal(
		[]byte(hash),
		[]byte(headerHash),
	)
}

// HashPayload computes the hash of payload's body according to the webhook's secret token see https://developer.github.com/webhooks/securing/#validating-payloads-from-github
// returning the hash as a hexadecimal string
func HashPayload(secret string, payloadBody []byte) string {
	hm := hmac.New(sha1.New, []byte(secret))
	hm.Write(payloadBody)
	sum := hm.Sum(nil)
	return fmt.Sprintf("%x", sum)
}

func IsValidPayloadSignature(secret, signatureHeader string, body []byte) (bool, error) {
	// Check header is valid
	signature_parts := strings.SplitN(signatureHeader, "=", 2)
	if len(signature_parts) != 2 {
		return false, fmt.Errorf("Invalid signature header: '%s' does not contain two parts (hash type and hash)", signatureHeader)
	}

	// Ensure secret is a sha1 hash
	signature_type := signature_parts[0]
	signature_hash := signature_parts[1]
	if signature_type != "sha1" {
		return false, fmt.Errorf("Signature should be a 'sha1' hash not '%s'", signature_type)
	}

	hash := HashPayload(secret, body)

	if !IsValidPayload(secret, signature_hash, body) {
		return false, fmt.Errorf("Payload did not come from GitHub, because secret is %s so hash should be %s", secret, hash)
	}

	return true, nil
}
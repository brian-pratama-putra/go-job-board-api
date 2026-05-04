package others

import (
	"crypto/sha256"
	"fmt"
)

func GenerateChecksum(p_payload string) string {
	v_hash := sha256.Sum256([]byte(p_payload))
	return fmt.Sprintf("%x", v_hash)
}

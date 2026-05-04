package others

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func ResponseService(p_c *gin.Context, p_method string, p_err_code interface{}, p_err_msg string, p_param map[string]interface{}) {
	v_request_datetime := ""
	if v_body, v_exists := p_c.Get("parsed_body"); v_exists {
		if v_map, v_ok := v_body.(map[string]interface{}); v_ok {
			if v_dt, v_ok2 := v_map["datetime"].(string); v_ok2 {
				v_request_datetime = v_dt
			}
		}
	}

	v_new_datetime  := NowJakartaStr()
	v_response      := map[string]interface{}{
		"err_code": p_err_code,
		"err_msg":  p_err_msg,
	}
	for v_k, v_v := range p_param {
		v_response[v_k] = v_v
	}

	v_payload_parts := []string{}
	for _, v_v := range v_response {
		if v_str, v_ok := v_v.(string); v_ok {
			v_payload_parts = append(v_payload_parts, v_str)
		}
	}
	v_payload_parts     = append(v_payload_parts, v_new_datetime)
	v_secret_response   := os.Getenv("SECRET_KEY_RESPONSE")
	v_payload_str       := strings.Join(v_payload_parts, "#") + "#" + v_secret_response
	v_new_checksum      := GenerateChecksum(v_payload_str)

	v_response["datetime"]  = v_new_datetime
	v_response["checksum"]  = v_new_checksum

	v_public_key := p_c.GetHeader("PublicKey")
	if v_public_key != "" {
		v_secret_header     := os.Getenv("SECRET_KEY_HEADER")
		v_sign_payload      := v_public_key + v_secret_header + v_request_datetime
		p_c.Header("SignatureKey", GenerateChecksum(v_sign_payload))
	}

	p_c.JSON(http.StatusOK, v_response)
}

package others

func GetStr(p_data map[string]interface{}, p_key string) string {
	if v_val, v_ok := p_data[p_key]; v_ok && v_val != nil {
		if v_str, v_ok2 := v_val.(string); v_ok2 {
			return v_str
		}
	}
	return ""
}

func GetSliceStr(p_data map[string]interface{}, p_key string) []string {
	v_result := []string{}
	if v_val, v_ok := p_data[p_key]; v_ok && v_val != nil {
		if v_slice, v_ok2 := v_val.([]interface{}); v_ok2 {
			for _, v_item := range v_slice {
				if v_str, v_ok3 := v_item.(string); v_ok3 {
					v_result = append(v_result, v_str)
				}
			}
		}
	}
	return v_result
}

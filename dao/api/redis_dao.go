package dao

import (
	"context"
	"encoding/json"
	"time"

	"go-job-board-api/config"
)

func CacheGet(p_key string) (map[string]interface{}, error) {
	v_client        := config.GetRedisClient()
	v_result, v_err := v_client.Get(context.Background(), p_key).Result()
	if v_err != nil {
		return nil, v_err
	}
	var v_data map[string]interface{}
	if v_err := json.Unmarshal([]byte(v_result), &v_data); v_err != nil {
		return nil, v_err
	}
	return v_data, nil
}

func CacheSet(p_key string, p_value interface{}, p_ttl int) error {
	v_client        := config.GetRedisClient()
	v_bytes, v_err  := json.Marshal(p_value)
	if v_err != nil {
		return v_err
	}
	return v_client.SetEx(context.Background(), p_key, string(v_bytes), time.Duration(p_ttl)*time.Second).Err()
}

func CacheDelete(p_key string) error {
	v_client := config.GetRedisClient()
	return v_client.Del(context.Background(), p_key).Err()
}

func CacheGetList(p_key string) ([]map[string]interface{}, error) {
	v_client        := config.GetRedisClient()
	v_result, v_err := v_client.Get(context.Background(), p_key).Result()
	if v_err != nil {
		return nil, v_err
	}
	var v_data []map[string]interface{}
	if v_err := json.Unmarshal([]byte(v_result), &v_data); v_err != nil {
		return nil, v_err
	}
	return v_data, nil
}

func CacheSetList(p_key string, p_value interface{}, p_ttl int) error {
	v_client        := config.GetRedisClient()
	v_bytes, v_err  := json.Marshal(p_value)
	if v_err != nil {
		return v_err
	}
	return v_client.SetEx(context.Background(), p_key, string(v_bytes), time.Duration(p_ttl)*time.Second).Err()
}

func CacheDeletePattern(p_pattern string) error {
	v_client    := config.GetRedisClient()
	v_ctx       := context.Background()
	var v_cursor uint64
	for {
		v_keys, v_next_cursor, v_err := v_client.Scan(v_ctx, v_cursor, p_pattern, 100).Result()
		if v_err != nil {
			return v_err
		}
		if len(v_keys) > 0 {
			v_client.Del(v_ctx, v_keys...)
		}
		v_cursor = v_next_cursor
		if v_cursor == 0 {
			break
		}
	}
	return nil
}

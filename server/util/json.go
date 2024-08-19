package util

import hjson "github.com/cloudwego/hertz/pkg/common/json"

func ToJsonString(v interface{}) string {
	bytes, _ := hjson.Marshal(v)
	return string(bytes)
}

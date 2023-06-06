package utils

import (
	"encoding/json"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"reflect"
)

func JsonDataList(resp interface{}) []map[string]interface{} {
	var list []proto.Message
	if reflect.TypeOf(resp).Kind() == reflect.Slice {
		s := reflect.ValueOf(resp)
		for i := 0; i < s.Len(); i++ {
			ele := s.Index(i)
			list = append(list, ele.Interface().(proto.Message))
		}
	}

	result := make([]map[string]interface{}, 0)
	for _, v := range list {
		m := ProtoToMap(v, false)
		result = append(result, m)
	}
	return result
}

func JsonDataOne(pb proto.Message) map[string]interface{} {
	return ProtoToMap(pb, false)
}

func ProtoToMap(pb proto.Message, idFix bool) map[string]interface{} {
	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: false,
	}

	s, _ := marshaler.MarshalToString(pb)
	out := make(map[string]interface{})
	json.Unmarshal([]byte(s), &out)
	if idFix {
		if _, ok := out["id"]; ok {
			out["_id"] = out["id"]
			delete(out, "id")
		}
	}
	return out
}

func VerifyByIDCard(idCard string) bool {
	if len(idCard) != 18 {
		return false
	}
	weight := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	validate := []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
	sum := 0
	for i := 0; i < len(weight); i++ {
		sum += weight[i] * int(byte(idCard[i])-'0')
	}
	m := sum % 11
	return validate[m] == idCard[17]
}

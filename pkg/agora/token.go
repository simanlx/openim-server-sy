package agora

import (
	"encoding/json"
	"errors"
	"fmt"
	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
	"log"
	"net/http"
	"strconv"
)

type rtc_int_token_struct struct {
	Uid_rtc_int  uint32 `json:"uid"`         // 用户的uid
	Channel_name string `json:"ChannelName"` // 频道名称
	Role         uint32 `json:"role"`        // 用户角色
}

var rtc_token string
var int_uid uint32
var channel_name string

var role_num uint32
var role rtctokenbuilder.Role

func rtcTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" && r.Method != "OPTIONS" {
		http.Error(w, "Unsupported method. Please check.", http.StatusNotFound)
		return
	}

	var t_int rtc_int_token_struct
	var unmarshalErr *json.UnmarshalTypeError
	int_decoder := json.NewDecoder(r.Body)
	int_err := int_decoder.Decode(&t_int)
	if int_err == nil {

		int_uid = t_int.Uid_rtc_int
		channel_name = t_int.Channel_name
		role_num = t_int.Role
		switch role_num {
		case 1:
			role = rtctokenbuilder.RolePublisher
		case 2:
			role = rtctokenbuilder.RoleSubscriber
		}
	}
	if int_err != nil {
		fmt.Println(int_err)
		if errors.As(int_err, &unmarshalErr) {
			errorResponse(w, "Bad request. Wrong type provided for field "+unmarshalErr.Value+unmarshalErr.Field+unmarshalErr.Struct, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad request.", http.StatusBadRequest)
		}
		return
	}

	errorResponse(w, rtc_token, http.StatusOK)
	log.Println(w, r)
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["token"] = message
	resp["code"] = strconv.Itoa(httpStatusCode)
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func main() {
	// 使用 int 型 uid 生成 RTC Token
	http.HandleFunc("/fetch_rtc_token", rtcTokenHandler)
	fmt.Printf("Starting server at port 8082\n")

	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatal(err)
	}
}

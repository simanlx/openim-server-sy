package notify

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestChargeNotify(t *testing.T) {

	// url ?redpacketID =10039
	ChargeNotifyReq := ChargeNotifyReq{
		Version:         "1",
		TranCode:        "2",
		MerOrderId:      "2",
		MerId:           "3",
		MerAttach:       "4",
		Charset:         "5",
		SignType:        "6",
		ResultCode:      "7",
		ErrorCode:       "7",
		ErrorMsg:        "7",
		OrderId:         "7",
		TranAmount:      "7",
		SubmitTime:      "7",
		TranFinishTime:  "7",
		BusinessType:    "7",
		FeeAmount:       "7",
		BankOrderId:     "7",
		RealBankOrderId: "7",
		DivideAcctDtl:   "7",
		SignValue:       "7",
		OperationID:     "7",
	}

	co, _ := json.Marshal(ChargeNotifyReq)
	fmt.Println(string(co))
}

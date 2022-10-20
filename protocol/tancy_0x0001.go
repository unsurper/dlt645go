package protocol

import (
	"encoding/hex"
	"strconv"
)

type Dlt_0x0091 struct {
	//修改前仪表参数
	DeviceID string
	//修改前仪表参数
	WP float64
}

func (entity *Dlt_0x0091) MsgID() MsgID {
	return Msgdlt_0x0091
}

func (entity *Dlt_0x0091) Encode() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (entity *Dlt_0x0091) Decode(data []byte) (int, error) {
	bytea := data[1:7]
	for i, j := 0, len(bytea)-1; i < j; i, j = i+1, j-1 {
		bytea[i], bytea[j] = bytea[j], bytea[i]
	}
	entity.DeviceID = hex.EncodeToString(bytea)
	byteb := make([]byte, 4)
	for i := 0; i < 4; i++ {
		byteb[i] = data[14+i] - 0x33
	}
	var err error
	entity.WP, err = stringtoWP(hex.EncodeToString(byteb))
	if err != nil {
		return 0, err
	}

	return len(data), nil
}

func stringtoWP(s string) (float64, error) {
	a0, _ := strconv.ParseFloat(s[0:2], 64)
	a1, _ := strconv.ParseFloat(s[2:4], 64)
	a2, _ := strconv.ParseFloat(s[4:6], 64)
	a3, _ := strconv.ParseFloat(s[6:8], 64)
	res := a0*0.01 + a1 + a2*100 + a3*10000
	return res, nil
}

package get

//Only for autenticated users

import (
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
)

func Init(t *pb.Request) (response *pb.Response) {

	if t.UID == nil {
		response = ErrorReturn(t, 401, "000011", "You are not authorised")
		return response
	}

	switch *t.Param {
	case "getuserbyid":
		response = GetUserByID(t)
	default:
		response = ErrorReturn(t, 406, "000012", "Wrong request")

	}

	return response

}

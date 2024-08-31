package get

//Only for autenticated users

import (
	. "auth/model"
	"fmt"

	"github.com/getsentry/sentry-go"
	. "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
)

type UserResponse struct {
	Name string `json:"name"`
	Mail string `json:"email"`
	UID  string `json:"uid"`
}

func GetUserByID(t *pb.Request) (response *pb.Response) {
	ans := make(map[string]interface{})
	args := ToMapStringInterface(t.Args)
	p := bluemonday.UGCPolicy()

	if args["uid"] == nil {
		return ErrorReturn(t, 500, "0000012", "Missing uid")
	}

	uid := p.Sanitize(fmt.Sprintf("%v", args["uid"]))

	//Check DB and table config
	db, err := ConnectDBv2()
	if err != nil {
		if viper.GetBool("server.sentry") {
			sentry.CaptureException(err)
		} else {
			SetErrorLog(err.Error())
		}
		return ErrorReturn(t, 500, "000027", err.Error())
	}

	var userExist Users

	rows := db.Conn.Debug().Where(`uid = ?`, uid).First(&userExist)

	if rows.RowsAffected == 0 {
		// return error. user name is exist in db users
		return ErrorReturn(t, 400, "000003", "There is no such user")
	}

	respuser := &UserResponse{
		Name: userExist.Name,
		Mail: userExist.Mail,
		UID:  userExist.UID,
	}

	ans["user"] = respuser

	response = Interfacetoresponse(t, ans)
	return response

}

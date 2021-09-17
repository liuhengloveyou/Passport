package face

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/liuhengloveyou/passport/common"
	"github.com/liuhengloveyou/passport/protos"
	"github.com/liuhengloveyou/passport/service"
	"github.com/liuhengloveyou/passport/sessions"

	gocommon "github.com/liuhengloveyou/go-common"
)


func TenantUserAdd(w http.ResponseWriter, r *http.Request) {
	sessionUser := r.Context().Value("session").(*sessions.Session).Values[common.SessUserInfoKey].(protos.User)
	if sessionUser.TenantID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		logger.Error("TenantUserAdd TenantID ERR")
		return
	}

	req := &protos.UserReq{}
	if err := readJsonBodyFromRequest(r, req, 1024); err != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		logger.Error("TenantUserAdd param ERR: ", err)
		return
	}
	logger.Infof("TenantUserAdd body: %#v\n", req)
	if len(req.Roles) > 10 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		logger.Error("TenantUserAdd param roles ERR")
		return

	}

	if req.UID == 0 {
		nuid, e := service.AddUserService(req)
		if e != nil {
			gocommon.HttpJsonErr(w, http.StatusOK, e)
			logger.Error("TenantUserAdd AddUserService ERR: ", e)
			return
		}

		req.UID = nuid
	}
	if err := service.TenantUserAdd(req.UID, sessionUser.TenantID, req.DepIds, req.Roles, req.Disable); err != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		logger.Error("TenantUserAdd service ERR: ", err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, req.UID)

	return
}

func TenantUserDel(w http.ResponseWriter, r *http.Request) {
	sessionUser := r.Context().Value("session").(*sessions.Session).Values[common.SessUserInfoKey].(protos.User)
	if sessionUser.TenantID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		logger.Error("TenantUserDel session ERR")
		return
	}

	req := &protos.Tenant{} // 只用一个UID字段
	if err := readJsonBodyFromRequest(r, req, 1024); err != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		logger.Error("TenantUserDel param ERR: ", err)
		return
	}
	logger.Infof("TenantUserDel body: %#v\n", req)

	_, e := service.TenantUserDel(req.UID, sessionUser.TenantID)
	if e != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, e)
		logger.Error("TenantUserDel service ERR: ", e)
	}

	gocommon.HttpJsonErr(w, http.StatusOK, common.ErrOK)

	return
}


func TenantUserGet(w http.ResponseWriter, r *http.Request) {
	sessionUser := r.Context().Value("session").(*sessions.Session).Values[common.SessUserInfoKey].(protos.User)
	if sessionUser.TenantID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		logger.Error("TenantUserGet TenantID ERR")
		return
	}

	var err error

	r.ParseForm()
	hasTotal, _ := strconv.ParseUint(r.FormValue("hasTotal"), 10, 64)
	nickname := strings.TrimSpace(r.FormValue("nickname"))
	uidStr := r.FormValue("uids")

	var uids []uint64
	if uidStr != "" {
		uidss := strings.Split(uidStr, ",")
		if len(uidss) <= 0 {
			gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
			logger.Error("TenantUserGet uids ERR: ", uidStr)
			return
		}
		uids = make([]uint64, len(uidss))

		for i, ouids := range uidss {
			if uids[i], err = strconv.ParseUint(ouids, 10, 64); err != nil {
				gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
				logger.Error("TenantUserGet uids ERR: ", uidStr, uidStr, uidss)
				return
			}
		}
	}

	page, _ := strconv.ParseUint(r.FormValue("page"), 10, 64)
	pageSize, _ := strconv.ParseUint(r.FormValue("pageSize"), 10, 64)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 1
	}
	if pageSize > 1000 {
		pageSize = 1000
	}

	rr, e := service.TenantUserGet(sessionUser.TenantID, page, pageSize, nickname, uids, hasTotal == 1)
	if e != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, e)
		logger.Error("TenantUserGet db ERR: ", e)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, rr)

	return
}

func TenantUserSetDepartment(w http.ResponseWriter, r *http.Request) {
	sessionUser := r.Context().Value("session").(*sessions.Session).Values[common.SessUserInfoKey].(protos.User)
	if sessionUser.TenantID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	var req protos.SetDepartmentReq
	if err := readJsonBodyFromRequest(r, &req, 1024); err != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		logger.Error("TenantUserSetDepartment param ERR: ", err)
		return
	}
	logger.Infof("TenantUserSetDepartment param: %v", req)

	if req.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		logger.Error("TenantUserSetDepartment param ERR: ", req)
		return
	}

	if err := service.TenantUserSetDepartment(req.UID, sessionUser.TenantID, req.DepIds); err != nil {
		logger.Errorf("TenantUserSetDepartment %v %s\n", req, err.Error())
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpJsonErr(w, http.StatusOK, common.ErrOK)
	return
}

func TenantUserDisableByUID(w http.ResponseWriter, r *http.Request) {
	sessionUser := r.Context().Value("session").(*sessions.Session).Values[common.SessUserInfoKey].(protos.User)
	if sessionUser.TenantID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		logger.Error("TenantID ERR")
		return
	}

	var req protos.DisableUserReq
	if err := readJsonBodyFromRequest(r, &req, 1024); err != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		logger.Error("TenantUserDisableByUID param ERR: ", err)
		return
	}

	if req.UID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		logger.Error("TenantUserDisableByUID param ERR: ", req)
		return
	}

	if req.Disable != 0 && req.Disable != 1 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		logger.Error("TenantUserDisableByUID param ERR: ", req)
		return
	}

	logger.Infof("TenantUserDisableByUID %v\n", req)

	if err := service.TenantUserDisabledService(req.UID, sessionUser.TenantID, int8(req.Disable)); err != nil {
		logger.Errorf("TenantUserDisableByUID %v %s\n", req, err.Error())
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpJsonErr(w, http.StatusOK, common.ErrOK)
	return
}

func TenantUserModifyExtInfo(w http.ResponseWriter, r *http.Request) {
	sessionUser := r.Context().Value("session").(*sessions.Session).Values[common.SessUserInfoKey].(protos.User)
	if sessionUser.TenantID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		logger.Error("TenantID ERR")
		return
	}

	var req protos.UserExtReq
	if err := readJsonBodyFromRequest(r, &req, 1024); err != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		logger.Error("TenantUserModifyExtInfo param ERR: ", err)
		return
	}
	logger.Infof("TenantUserModifyExtInfo %v\n", req)

	if err := service.TenantUpdateUserExt(req.UID, sessionUser.TenantID, req.K, req.V); err != nil {
		logger.Errorf("TenantUserModifyExtInfo %v %s\n", req, err.Error())
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpJsonErr(w, http.StatusOK, common.ErrOK)
	return
}

func TenantUserModifyPWDByUID(w http.ResponseWriter, r *http.Request) {
	sessionUser := r.Context().Value("session").(*sessions.Session).Values[common.SessUserInfoKey].(protos.User)
	if sessionUser.TenantID <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrNoAuth)
		return
	}

	var req map[string]interface{}
	if err := readJsonBodyFromRequest(r, &req, 1024); err != nil {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		logger.Error("modifyPWDByID param ERR: ", err)
		return
	}

	uid, ok := req["uid"].(float64)
	if !ok || uint64(uid) <= 0 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		logger.Error("modifyPWDByUID param ERR: ", req)
		return
	}

	pwd, ok := req["pwd"].(string)
	pwd = strings.TrimSpace(pwd)
	if len(pwd) < 4 || len(pwd) > 16 {
		gocommon.HttpJsonErr(w, http.StatusOK, common.ErrParam)
		logger.Error("modifyPWDByUID param ERR: ", req)
		return
	}

	logger.Infof("modifyPWDByUID %v %s\n", uid, pwd)

	if _, err := service.SetUserPWD(uint64(uid), sessionUser.TenantID, pwd); err != nil {
		logger.Errorf("modifyPWDByUID %v %s %s\n", uid, pwd, err.Error())
		gocommon.HttpJsonErr(w, http.StatusOK, err)
		return
	}

	gocommon.HttpErr(w, http.StatusOK, 0, "OK")
	return
}

package gf_identity_core

import (
	"context"
	"net/http"
	gf_core "github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
// LOGOUT_PIPELINE

func LogoutPipeline(pGFsessionID gf_core.GF_ID,
	pDomainForAuthCookiesStr string,
	pResp       http.ResponseWriter,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (string, *gf_core.GFerror) {
	
	//---------------------
	// DB
	session, gfErr := DBsqlGetSession(pGFsessionID, pCtx, pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	//---------------------
	
	logoutSuccessRedirectURLstr := session.LogoutSuccessRedirectURLstr

	//---------------------
	// DB
	gfErr = dbSQLdeleteSession(pGFsessionID, pCtx, pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	//---------------------
	// unset all session cookies
	DeleteCookies(pDomainForAuthCookiesStr, pResp)

	//---------------------

	return logoutSuccessRedirectURLstr, nil
}
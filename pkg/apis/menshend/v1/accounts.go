package v1

import (
    vault "github.com/hashicorp/vault/api"
    "github.com/emicklei/go-restful"
    "time"
    "fmt"
    mfilters "github.com/nebtex/menshend/pkg/filters"
    mutils "github.com/nebtex/menshend/pkg/utils"
)

//AuthResource ...
type AuthResource struct {
}

//LoginStatus store the vault token state relative to  menshend
type LoginStatus struct {
    IsLogged         bool `json:"isLogged"`
    IsAdmin          bool `json:"isAdmin"`
    CanImpersonate   bool `json:"canImpersonate"`
    SessionExpiresAt int64 `json:"sessionExpiresAt"`
}

//Register ...
func (a *AuthResource) Register(container *restful.Container) {
    ws := new(restful.WebService).
        Consumes(restful.MIME_JSON).
        Produces(restful.MIME_JSON)
    
    ws.Path("/v1/account").
        Doc("login service")
    
    ws.Route(ws.GET("").To(a.accountStatus).
        Doc("get login status").
        Operation("loginStatus").
        Writes(LoginStatus{}))
    
    ws.Route(ws.DELETE("").To(a.logout).
        Doc("get login status").
        Operation("logout"))
    container.Add(ws)
    
}

//logout try to revoke the vault token, just  a wrap over the vault enpoint
func (a *AuthResource) logout(request *restful.Request, response *restful.Response) {
    defer func() {}()
    user := mfilters.GetTokenFromRequest(request)
    vc, err := vault.NewClient(vault.DefaultConfig())
    mutils.HttpCheckPanic(err, mutils.InternalError)
    vc.SetToken(user)
    err = vc.Auth().Token().RevokeSelf(user)
    mutils.HttpCheckPanic(err, mutils.PermissionError)
}

//accountStatus, if the token is active this will return some important info like
//admin and impersonate capabilities
func (*AuthResource)accountStatus(request *restful.Request, response *restful.Response) {
    var creationTimeMillisecond int64
    var ttl int64
    
    defer func() {
        r := recover()
        fmt.Println(r)
        if r != nil {
            ls := LoginStatus{
                IsLogged: false,
                IsAdmin: false,
                CanImpersonate: false,
                SessionExpiresAt: 0,
            }
            mutils.HttpCheckPanic(response.WriteEntity(ls), mutils.InternalError)
        }
    }()
    user := mfilters.GetTokenFromRequest(request)
    vc, err := vault.NewClient(vault.DefaultConfig())
    mutils.CheckPanic(err)
    vc.SetToken(user)
    secret, err := vc.Auth().Token().LookupSelf()
    mutils.HttpCheckPanic(err, mutils.NotAuthorized)
    CheckSecretFailIfIsNull(secret)
    if (secret.WrapInfo != nil) {
        creationTimeMillisecond = secret.WrapInfo.CreationTime.UnixNano() / int64(time.Millisecond)
        ttl = int64(secret.WrapInfo.TTL) * 1000
    }
    
    ls := LoginStatus{
        IsLogged: true,
        IsAdmin: mfilters.IsAdmin(user),
        CanImpersonate: mfilters.CanImpersonate(user),
        SessionExpiresAt: creationTimeMillisecond + ttl,
    }
    mutils.HttpCheckPanic(response.WriteEntity(ls), mutils.InternalError)
}

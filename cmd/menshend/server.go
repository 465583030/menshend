package main

import (
    "net/http"
    "github.com/nebtex/menshend/pkg/apis/menshend/v1"
    "fmt"
    "github.com/Sirupsen/logrus"
    "github.com/rakyll/statik/fs"
    
    . "github.com/nebtex/menshend/statik"
    "github.com/gorilla/mux"
)

func mainHandler(response http.ResponseWriter, request *http.Request) {
    //detect  menshend host
    //use proxy server
    
}

func proxyServer() http.Handler {
    return PanicHandler(GetSubDomainHandler(v1.BrowserDetectorHandler(
        NeedLogin(RoleHandler(GetServiceHandler(ProxyHandler()))))))
}

func ui() http.Handler {
    r := mux.NewRouter()
    statikFS, _ := fs.New()
    r.PathPrefix("/static").Handler(http.FileServer(statikFS))
    r.PathPrefix("/").Handler(Index())
    return r
}
func uiLogin(){
    //add csfr protection
    // use flash for error
    // redirect to service
    //check same origin policy
    //read form, if github use gomniauth
    //token use token login
    //vault, rados, ldap
}
func menshendServer(address, port string) error {
    // /ui
    // /uilogin
    // /uilogout
    // /v1 - api
    //http.Handle("/api", v1.APIHandler())
    http.Handle("/", ui())
     
    logrus.Infof("Server listing on %s:%s", address, port)
    return http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), nil)
    
    //r.Host("{subdomain:[a-z\\-]+}." + Config.HostWithoutPort()).Handler(handler)
}

/*
func setToken(u *User, expiresIn int64, response *restful.Response, hasCSRF bool) {
    expireAt := MakeTimestampMillisecond()
    if expiresIn == 0 {
        expireAt += Config.DefaultTTL
    } else {
        expireAt += expiresIn
    }
    u.SetExpiresAt(expireAt)

    if !hasCSRF {
        response.AddHeader("X-Menshend-Token", u.GenerateJWT())

    } else {
        ct := &http.Cookie{Path: "/", Name: "X-Menshend-Token", Value: u.GenerateJWT(),
            Expires: time.Unix(u.ExpiresAt / 1000, 0),
            HttpOnly:true }

        ct.Domain = "." + Config.HostWithoutPort()

        if Config.Scheme == "https" {
            ct.Secure = true
        }
        http.SetCookie(response.ResponseWriter, ct)

    }
}
func setExpirationTime(u *User, expiresIn int64) *User {
    expireAt := MakeTimestampMillisecond()
    if expiresIn == 0 {
        expireAt += Config.DefaultTTL
    } else {
        expireAt += expiresIn
    }
    u.SetExpiresAt(expireAt)
    return u

}
*/



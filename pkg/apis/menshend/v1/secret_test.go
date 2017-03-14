package v1

import (
    "testing"
    . "github.com/smartystreets/goconvey/convey"
    "net/http/httptest"
    "github.com/emicklei/go-restful"
    "net/http"
    . "github.com/nebtex/menshend/pkg/utils/test"
    . "github.com/nebtex/menshend/pkg/config"
    . "github.com/nebtex/menshend/pkg/users"
    . "github.com/nebtex/menshend/pkg/utils"
    . "github.com/nebtex/menshend/pkg/apis/menshend"
    
    "encoding/json"
    
    "io/ioutil"
    vault "github.com/hashicorp/vault/api"
    "github.com/ansel1/merry"
)

func Test_SecretEndpoint(t *testing.T) {
    VaultConfig.Address = "http://localhost:8200"
    Convey("should return info about the envirenment", t, func(c C) {
        CleanVault()
        PopulateVault()
        httpReq, err := http.NewRequest("GET", "/v1/secret/roles/ml-team/gitlab./" + Config.VaultPath + "/roles/ml-team/gitlab.", nil)
        So(err, ShouldBeNil)
        httpReq.Header.Set("Content-Type", "application/json")
        user, err := NewUser("myroot")
        So(err, ShouldBeNil)
        user.TokenLogin()
        user.SetExpiresAt(GetNow() + 3600)
        httpReq.Header.Add("X-Menshend-Token", user.GenerateJWT())
        httpWriter := httptest.NewRecorder()
        wsContainer := restful.NewContainer()
        s := SecretResource{}
        s.Register(wsContainer)
        wsContainer.ServeHTTP(httpWriter, httpReq)
        jsres, err := ioutil.ReadAll(httpWriter.Body)
        So(err, ShouldBeNil)
        secret := &vault.Secret{}
        err = json.Unmarshal(jsres, secret)
        So(err, ShouldBeNil)
        So(httpWriter.Body, ShouldNotBeNil)
        So(httpWriter.Code, ShouldEqual, 200)
        
    })
    
    Convey("fails if secret does not exists", t, func(c C) {
        CleanVault()
        PopulateVault()
        defer func() {
            r := recover()
            if (r == nil) {
                t.Error("did not panicked")
                t.Fail()
            }
            switch x := r.(type) {
            case error:
                c.So(merry.Is(x, NotFound), ShouldBeTrue)
            default:
                t.Errorf("%v", x)
                t.Fail()
            }
        }()
        httpReq, err := http.NewRequest("GET", "/v1/secret/roles/ml-team/gitlab./secret/gitlab/password", nil)
        So(err, ShouldBeNil)
        httpReq.Header.Set("Content-Type", "application/json")
        user, err := NewUser("myroot")
        So(err, ShouldBeNil)
        user.TokenLogin()
        user.SetExpiresAt(GetNow() + 3600)
        httpReq.Header.Add("X-Menshend-Token", user.GenerateJWT())
        httpWriter := httptest.NewRecorder()
        wsContainer := restful.NewContainer()
        s := SecretResource{}
        s.Register(wsContainer)
        wsContainer.ServeHTTP(httpWriter, httpReq)
        
    })
}
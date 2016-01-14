package userAuth
import (
	"github.com/go-martini/martini"
	"fmt"
	"github.com/hoysoft/JexGO/sessions"
	"net/http"
	"net/url"
	"io/ioutil"
	"bytes"
	"time"
	"encoding/json"
	"strings"
	"errors"
	"github.com/hoysoft/JexGO/utils"
)

var (
// RedirectParam is the query string parameter that will be set
// with the page the user was trying to visit before they were
// intercepted.
	RedirectParam string = "next"
	PathRegister string = "/register"
	PathLogin string = "/login"
	PathLogout string = "/logout"
)

type Config struct {
	// ClientID is the application's ID.
	ClientID     string

	// ClientSecret is the application's secret.
	ClientSecret string

	// Endpoint contains the resource server's token endpoint
	// URLs. These are constants specific to each server and are
	// often available via site-specific packages, such as
	// google.Endpoint or github.Endpoint.
	Endpoint     Endpoint

	// RedirectURL is the URL to redirect users going through
	// the OAuth flow, after the resource owner's URLs.
	RedirectURL  string

	// Scope specifies optional requested permissions.
	Scopes       []string

	Callback     func(string)
}

type Endpoint struct {
	RegisterURL string
	LoginURL    string
	LogoutURL   string
	AuthURL     string
	TokenURL    string
}

func NewOAuth2Provider(conf *Config) martini.Handler {

	return func(s sessions.Session, c martini.Context, w http.ResponseWriter, r *http.Request) {
		//fmt.Println("prover:", r)
		conf.RedirectURL = "http://" + r.Host + PathCallback
		if r.Method == "GET" {
			switch r.URL.Path {
			case PathRegister:
				handleOAuth2Register(conf, s, w, r)
			case PathLogin:
				loginHandle(conf, s, w, r)
			case PathLogout:
				logoutHandle(conf, c, s, w, r)
			case PathCallback:
				callbackhandle(conf, c, s, w, r)
			case PathError:
				handleOAuth2Error(s, w, r);
			}
		}
		if r.Method == "POST" {
			switch r.URL.Path {

			case PathCallback:
				callbackhandle(conf, c, s, w, r)

			}
		}
		tk := unmarshallToken(s)
		if tk != nil {
			// check if the access token is expired
			if tk.Expired() && tk.Refresh() == "" {
				s.Delete(keyToken)
				tk = nil
			}
		}
		// Inject tokens.
		c.MapTo(tk, (*IToken)(nil))
	}
}

func loginHandle(f *Config, s sessions.Session, w http.ResponseWriter, r *http.Request) {
	next := extractPath(r.URL.Query().Get(RedirectParam))
	s.Set("_RedirectURL", next)
//	if len(f.ClientID) > 0 && len(f.ClientSecret) > 0 {
//		http.Redirect(w, r, f.authCodeURL(), 302)
//	}else {
		path := fmt.Sprintf("%s?redirect_uri=%s", f.Endpoint.LoginURL, f.RedirectURL)
		http.Redirect(w, r, path, 302)
//	}
}

//func getCode(url string) (string,error){
//	 bytestr,err:=spider.HttpGetString(url)
//	if err!=nil {
//		return err
//	}
//
//}

func logoutHandle(f *Config, c martini.Context, s sessions.Session, w http.ResponseWriter, r *http.Request) {
	s.Delete(keyToken)
	path := fmt.Sprintf("%s?client_id=%s&client_secret=%s", f.Endpoint.LogoutURL, f.ClientID, f.ClientSecret)
	utils.HttpGetString(path)
	//	fmt.Println("oauth logout result:",string(str))
	f.ClientID = ""
	f.ClientSecret = ""
	c.Invoke(Logout)
	http.Redirect(w, r, "/", 302)
}

func handleOAuth2Error(s sessions.Session, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "登陆失败")
}

func handleOAuth2Register(f *Config, s sessions.Session, w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("%s?redirect_uri=%s", f.Endpoint.RegisterURL, f.RedirectURL)
	http.Redirect(w, r, path, 302)
}

func callbackhandle(f *Config, c martini.Context, s sessions.Session, w http.ResponseWriter, r *http.Request) {
	rurl, _ := s.Get("_RedirectURL").(string)
	rurl=extractPath(rurl)
	if (len(r.URL.Query().Get("code")) > 0) {
		//获取token
		tk, error := f.authTokenURL(r.URL.Query().Get("code"))

		if error == nil && tk.Valid() {
			val, _ := json.Marshal(tk)
			s.Set(keyToken, val)
			fmt.Println("登陆成功")
			s.AddFlash("登陆成功")
			c.Invoke(oAuthUserLoginCallback)
			if len(rurl)==0 {
				rurl="/"
			}
			fmt.Println("rul:",rurl)
			http.Redirect(w, r, rurl, 302)
			return
		}else {
			s.AddFlash( "登陆失败")
			http.Redirect(w, r, PathError, 302)
			return
		}
	}
	if len(r.URL.Query().Get("client_id")) > 0 {
		f.ClientID = r.URL.Query().Get("client_id")
		f.ClientSecret = r.URL.Query().Get("client_secret")
		http.Redirect(w, r, f.authCodeURL(), 302)
		return
	}


	//	fmt.Println("call:",r)
	//	if (len(r.URL.Query().Get("code"))>0) {
	//		fmt.Println("callFUN:q token")
	//		rurl,_:=  s.Get("_RedirectURL").(string)
	//		//获取token
	//		tk, error := f.authTokenURL(r.URL.Query().Get("code"))
	//		if error==nil && tk.Valid() {
	//			val, _ := json.Marshal(tk)
	//			s.Set(keyToken, val)
	//			s.AddFlash("success","登陆成功")
	//			c.Invoke(oAuthUserLogin)
	//
	//
	//			 http.RedirectHandler(rurl, 302)
	//		 	 return
	//		}else{
	//			s.AddFlash("warning","登陆失败")
	//			http.Redirect(w, r, PathError, 302)
	//			return
	//		}
	//	}else{
	//		fmt.Println("callFUN:1111")
	//		//获取code
	//		if len(r.URL.Query().Get("client_id"))>0 {
	//			fmt.Println("callFUN:q code")
	//			f.ClientID = r.URL.Query().Get("client_id")
	//			f.ClientSecret = r.URL.Query().Get("client_secret")
	//			http.Redirect(w, r, f.authCodeURL(), 302)
	//			return
	//		}
	//		http.Redirect(w, r, "/", 302)
	//	}
	//	fmt.Println("callFUN:sppp")
}


func (c *Config) authCodeURL() string {
	var buf bytes.Buffer
	buf.WriteString(c.Endpoint.AuthURL)
	v := url.Values{
		"response_type": {"code"},
		"client_id":     {c.ClientID},
		"redirect_uri":  CondVal(c.RedirectURL),
		"scope":         CondVal(strings.Join(c.Scopes, " ")),
	}

	if strings.Contains(c.Endpoint.AuthURL, "?") {
		buf.WriteByte('&')
	} else {
		buf.WriteByte('?')
	}
	buf.WriteString(v.Encode())
	fmt.Println(buf.String())
	return buf.String()
}

func (c *Config) authTokenURL(code string) (*token, error) {
	//	测试地址：https://oauthtest.lecloud.com/accesstoken?grant_type=authorization_code&code=b86e7f2393bce2eaeb67a37f401ef5ea&client_id=clientid-dabingge&client_secret=clientsecret-dabingge&redirect_uri=http://www.baidu.com
	urlstr := c.Endpoint.TokenURL + "?grant_type=authorization_code&code=" + code + "&client_id=" + c.ClientID + "&client_secret=" + c.ClientSecret //+ "&redirect_uri="+c.RedirectURL
	fmt.Println(urlstr)
	res, err := http.Get(urlstr)
	if err != nil {
		return nil, err
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	//json str 转struct
	tk := token{}
	//	fmt.Println(string(result))
	if err := json.Unmarshal(result, &tk); err == nil && len(tk.Error) == 0 {
		tk.Expiry = time.Now().Add(tk.Expires_in * time.Second)

		return &tk, nil
	}
	if len(tk.Error) > 0 {
		return nil, errors.New(tk.Error)
	}
	return nil, err
}

func CondVal(v string) []string {
	if v == "" {
		return nil
	}
	return []string{v}
}

func extractPath(next string) string {
	n, err := url.Parse(next)
	if err != nil {
		return "/"
	}
	return n.Path
}
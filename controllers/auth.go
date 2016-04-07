package controllers

import (
	"fmt"
	"os"
	"io"
	"strings"
	"github.com/markbates/goth/gothic"
	"net/http"
	//"html/template"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/digitalocean"
	"github.com/markbates/goth/providers/dropbox"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/instagram"
	"github.com/markbates/goth/providers/lastfm"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/spotify"
	"github.com/markbates/goth/providers/twitch"
	"github.com/markbates/goth/providers/twitter"
	"github.com/markbates/goth/providers/box"
	"github.com/markbates/goth/providers/salesforce"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/yammer"
	"github.com/markbates/goth/providers/onedrive"
	"github.com/markbates/goth"
	"net/url"
	"github.com/eraserxp/coedit/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/session"
)

var globalSessions *session.Manager


func init()  {
	gothic.Store = sessions.NewFilesystemStore(os.TempDir(), []byte("goth-example"))
	goth.UseProviders(
		twitter.New(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), "http://localhost:3000/auth/twitter/callback"),
		// If you'd like to use authenticate instead of authorize in Twitter provider, use this instead.
		// twitter.NewAuthenticate(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), "http://localhost:3000/auth/twitter/callback"),

		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), "http://localhost:3000/auth/facebook/callback"),
		gplus.New("281140391713-b1dskle4dtsi6nn4ce01tbkpcp3aovs6.apps.googleusercontent.com", "cIM92vsFvLyfhIZASmAo2ZaE", "http://localhost:8080/auth/gplus/callback"),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:3000/auth/github/callback"),
		spotify.New(os.Getenv("SPOTIFY_KEY"), os.Getenv("SPOTIFY_SECRET"), "http://localhost:3000/auth/spotify/callback"),
		linkedin.New(os.Getenv("LINKEDIN_KEY"), os.Getenv("LINKEDIN_SECRET"), "http://localhost:3000/auth/linkedin/callback"),
		lastfm.New(os.Getenv("LASTFM_KEY"), os.Getenv("LASTFM_SECRET"), "http://localhost:3000/auth/lastfm/callback"),
		twitch.New(os.Getenv("TWITCH_KEY"), os.Getenv("TWITCH_SECRET"), "http://localhost:3000/auth/twitch/callback"),
		dropbox.New(os.Getenv("DROPBOX_KEY"), os.Getenv("DROPBOX_SECRET"), "http://localhost:3000/auth/dropbox/callback"),
		digitalocean.New(os.Getenv("DIGITALOCEAN_KEY"), os.Getenv("DIGITALOCEAN_SECRET"), "http://localhost:3000/auth/digitalocean/callback", "read"),
		bitbucket.New(os.Getenv("BITBUCKET_KEY"), os.Getenv("BITBUCKET_SECRET"), "http://localhost:3000/auth/bitbucket/callback"),
		instagram.New(os.Getenv("INSTAGRAM_KEY"), os.Getenv("INSTAGRAM_SECRET"), "http://localhost:3000/auth/instagram/callback"),
		box.New(os.Getenv("BOX_KEY"), os.Getenv("BOX_SECRET"), "http://localhost:3000/auth/box/callback"),
		salesforce.New(os.Getenv("SALESFORCE_KEY"), os.Getenv("SALESFORCE_SECRET"), "http://localhost:3000/auth/salesforce/callback"),
		amazon.New(os.Getenv("AMAZON_KEY"), os.Getenv("AMAZON_SECRET"), "http://localhost:3000/auth/amazon/callback"),
		yammer.New(os.Getenv("YAMMER_KEY"), os.Getenv("YAMMER_SECRET"), "http://localhost:3000/auth/yammer/callback"),
		onedrive.New(os.Getenv("ONEDRIVE_KEY"), os.Getenv("ONEDRIVE_SECRET"), "http://localhost:3000/auth/onedrive/callback"),


	)

	//set a global session
	globalSessions, _ = session.NewManager("memory", `{"cookieName":"gosessionid", "enableSetCookie,omitempty": true, "gclifetime":3600, "maxLifetime": 3600, "secure": false, "sessionIDHashFunc": "sha1", "sessionIDHashKey": "", "cookieLifeTime": 3600, "providerConfig": ""}`)
	go globalSessions.GC()
}

//I need to add query variable to make goth works with beego
func addQueryVars(r *http.Request, vars map[string]string) {
	parts, i := make([]string, len(vars)), 0
	for key, value := range vars {
		parts[i] = url.QueryEscape(":"+key) + "=" + url.QueryEscape(value)
		i++
	}
	q := strings.Join(parts, "&")
	if r.URL.RawQuery == "" {
		r.URL.RawQuery = q
	} else {
		r.URL.RawQuery += "&" + q
	}
}

func preprocessUrl(req *http.Request)  {
	array := strings.Split(req.URL.Path, "/")
	queries := map[string]string {
		"provider": array[2],
	}
	addQueryVars(req, queries)
}

type AuthHandler struct {

}

func (this *AuthHandler) ServeHTTP(res http.ResponseWriter, req *http.Request)  {
	preprocessUrl(req)
	auth(res, req)
}

type AuthCallbackHandler struct {

}

func (this *AuthCallbackHandler) ServeHTTP(res http.ResponseWriter, req *http.Request)  {
	preprocessUrl(req)
	authCallback(res, req)
}

func authCallback(res http.ResponseWriter, req *http.Request) {

	// print our state string to the console. Ideally, you should verify
	// that it's the same string as the one you set in `setState`
	fmt.Println("State: ", gothic.GetState(req))
	fmt.Println( "request method: " + req.Method)

	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		fmt.Fprintln(res, err)
		return
	}

	//t, _ := template.New("foo").Parse(userTemplate)

	account := &models.Account{ user.Email, ""}
	fmt.Println( account.CheckExist( ) )

	//if everything is fine, set the session for the current user
	sess, err := globalSessions.SessionStart(res, req)
	if err != nil {
		fmt.Println("set error,", err)
	}
	defer sess.SessionRelease(res)
	err = sess.Set("username", user.Email)
	if err != nil {
		fmt.Println("set error,", err)
	}

	http.Redirect( res, req, "/user/" + user.Email, http.StatusFound)
	//t.Execute(res, user)
}

func auth(res http.ResponseWriter, req *http.Request)  {
	gothic.BeginAuthHandler(res, req)
}

var userTemplate = `
<p>Name: {{.Name}}</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
`
type AccountController struct {
	beego.Controller
}

type FileList struct{
	Name string
	Url  string
}

//check whether the current user has logged in or not
func isLogin(this *AccountController) bool {
	//try to retrieve the session
	w := this.Ctx.ResponseWriter
	r := this.Ctx.Request
	sess, err := globalSessions.SessionStart(w, r)
	if (err != nil) {
		return false
	}
	username := sess.Get("username")
	fmt.Println("get username from session: " + username.(string))
	user := this.Ctx.Input.Param(":uemail")

	return username == user
}

func (a *AccountController) Get() {
	//if not logged in, redirect to the main page
	if (!isLogin(a)) {
		a.Redirect("/", 302)
	}

	user_email := a.Ctx.Input.Param(":uemail")

	fmt.Println("AccountController for " + user_email)
	a.Data["Email"] = user_email

	account := &models.Account{ user_email, ""}
	a.Data["Options"] = account.SearchDocument()


	//	c.Data["Website"] = "beego.me"
	//	c.Data["Email"] = "astaxie@gmail.com"
	a.TplName = "user.tpl"
}

type RequestUserListHandler struct {

}

func (this *RequestUserListHandler) ServeHTTP(res http.ResponseWriter, req *http.Request)  {
	RequestUserList(res, req)
}

func RequestUserList(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
		case "GET":
			sess, _ := globalSessions.SessionStart(res, req)
			username := sess.Get("username")
			account := &models.Account{ username.(string), ""}

			jsonlist := account.SearchDocument()
			fmt.Println(jsonlist)
			io.WriteString(res, jsonlist)

		default:

	}
}

type LogoutHandler struct {

}

func (this *LogoutHandler) ServeHTTP(res http.ResponseWriter, req *http.Request)  {
	Logout(res, req)
}

func Logout (res http.ResponseWriter, req *http.Request) {
	globalSessions.SessionDestroy( res, req)
}
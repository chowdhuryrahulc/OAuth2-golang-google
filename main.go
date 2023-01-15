package main

// source youtube: https://www.youtube.com/watch?v=OdyXIi6DGYw
// source github: https://github.com/plutov/packagemain/blob/master/11-oauth2/main.go


// CONCLUSION: 
// 		1) go to google devloper console and do the stuff (described below)
//			dont forget to add redirect_url. 
// 			And you get client_id and secret_key from google devloper console
// 		2) init, add scopes bout the number of access you want. Eg: username, password
// 		3) .env file: put client_id and secreat_key in ,env. Exactly as done. Without comma or ""
// 		4) Flow: goes to handleMain. Where the UI of button is located
// 		5) When user clicks button, sends to "/login". Which triggers handleGoogleLogin function
// 		6) After that.... "/callback".
// 			The OAuth2 redirect_url which we added in google devloper console redirects the user to that url
// 			thats how handleGoogleCallback() gets executed. 
// 			Where we get user info. It gets in this case printed in browser
// 		7) Return siring is given below, So we can store that stuff

/*
enviournment variables are key value pairs
in terminal, type: echo $key	Value will be returned

go to google devloper console --> new project --> create --> 
oauth concent screen --> external --> create --> add app name, app logo, domain(Dont Know What and how) etc --> 
scopes --> emailid, profile, openid --> test (add your email etc)
credentials --> create credentials --> oauth client id (from dropdown) --> web application 
server side, so web. Not to change to client/app side. That requires androidmanifest credentials --> get client_id and secret key
Client Id: 1026515211019-msg227mh02p7udmpddneutn0crem80vq.apps.googleusercontent.com
Client Secret: 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
DID NOT DO ENOUGH FOR PRODUCTION. RESEARCH MORE WHEN SENDING TO PROD. NOW IT IS TESTING

!problem1: Missing required parameter: client_id
Sol: godotenv.Load(".env")
flow: init happens before the windows dilog to allow access(means before http.serve)
and handleGoogleLogin is also called. When we press the googlelogin button in front-end 

!problem2: The OAuth client was not found.
Sol: remove , from .env file (SOLVED)

!problem3: You can’t sign in because this app sent an invalid request. You can try again later or contact the developer about this issue
!Error 400: redirect_uri_mismatch
Sol: we should have put "redirect_uri" in google devloper conole while creating OAuth server
source: https://stackoverflow.com/questions/11485271/google-oauth-2-authorization-error-redirect-uri-mismatch
1st answer
https://console.cloud.google.com/apis/credentials?project=oauth2-golang-trial
edit OAuth2ClientId


? PROJECT SUCCESSFUL
RETURN: 
Content: {
  "id": "113689114832707015366",
  "email": "chowdhuryrahulc@gmail.com",
  "verified_email": true,
  "name": "Rahul Chowdhury",
  "given_name": "Rahul",
  "family_name": "Chowdhury",
  "picture": "https://lh3.googleusercontent.com/a/AEdFTp6DniEf9IOaQUmsg_WEMJPxpxhKvx9LyNEPqJfHqw=s96-c",
  "locale": "en-GB"
}


*/


import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"github.com/joho/godotenv"
)

var (
	googleOauthConfig *oauth2.Config
	// TODO: randomize it
	oauthStateString = "pseudo-random"
)

func init() {
	godotenv.Load(".env")				// loads .env file to our file sothat we can access client_id and client_secret
	// fmt.Println("init")
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

func main() {
	//todo How does the user gets redirected when he presses google-login button? And to which url?
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleGoogleLogin)
	http.HandleFunc("/callback", handleGoogleCallback)
	fmt.Println(http.ListenAndServe(":8080", nil))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	//!What hapens when user click google-signin?
	//(IMP) Sol: redirected to "/login". As shown below. Which leads to handleGoogleLogin()function getting executed
	var htmlIndex = `<html>
<body>
	<a href="/login">Google Log In</a>
</body>
</html>`
//fmt.Println(w) // &{0xc000001ae0 0xc000156000 {} 0xcf8640 false false false false 0 {0 0} 0xc00004e6c0 {0xc00015e000 map[] false false} map[] false 0 -1 0 false false [] 0 [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0 0 0 0] [0 0 0] 0xc00001a7e0 0}
//fmt.Println(htmlIndex) // <html><body><a href="/login">Google Log In</a></body></html>
	fmt.Fprintf(w, htmlIndex)
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	// fmt.Println("handle logiGoogleLogin")
	fmt.Println(url) // https://accounts.google.com/o/oauth2/auth?client_id=1026515211019-msg227mh02p7udmpddneutn0crem80vq.apps.googleusercontent.com&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.profile&state=pseudo-random
	fmt.Println(http.StatusTemporaryRedirect)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("calback")
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "Content: %s\n", content)
}

func getUserInfo(state string, code string) ([]byte, error) {
	fmt.Println("get user info")
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, nil
}
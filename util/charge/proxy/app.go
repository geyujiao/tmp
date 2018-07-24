package proxy

import "sync"

const (
	Success = iota
	ReceFail
	SendFail
	ServerClose
)

var (
	app  *App
	once sync.Once
)

type App struct {
	AppId        string
	AppSecret    string
	AppMd5Secret string
	Url          string
}

func NewClient() *App {
	once.Do(func() {
		app = &App{
			AppId:        "1045",
			AppSecret:    "zhwwtoo786bbsome",
			AppMd5Secret: "o61uswq6",
			Url: "",
		}
	})

	return app
}

func (app *App) Set(appid, appsecret, appmd5secret string) {
	app.AppId = appid
	app.AppSecret = appsecret
	app.AppMd5Secret = appmd5secret
}
func (app *App) SetUrl(url string) {
	app.Url = url
}

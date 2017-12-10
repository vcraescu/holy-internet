package holyinternet

import (
	"github.com/getlantern/systray"
	"github.com/vcraescu/holy-internet/pkg/icon"
	"log"
)

type App struct {
	MailerDaemon *MailerDaemon
	pause        bool
	errorCh      chan bool
}

type pauseMenuItem struct {
	*systray.MenuItem
	state bool
}

func (item *pauseMenuItem) Toggle() {
	item.state = !item.state
	if item.state {
		item.SetTitle("Resume")
		item.SetTooltip("Resume your prayings")
		return
	}

	item.SetTitle("Pause")
	item.SetTooltip("Pause your prayings")
}

func (a *App) Pause() {
	a.pause = true
}

func (a *App) Resume() {
	a.pause = false
}

func (a App) IsPaused() (bool) {
	return a.pause
}

func (a *App) Failed() {
	a.errorCh <- true
}

func (a *App) DiscardFailure() {
	a.errorCh <- false
}

func (a *App) Run(onReady func(app *App), onExit func(app *App)) {
	systray.Run(func() {
		onAppReady(a)
		onReady(a)
	}, func() {
		onExit(a)
	})
}

func onAppReady(a *App) {
	systray.SetIcon(icon.HolyActive)
	systray.SetTooltip("Holy Internet")

	mPause := &pauseMenuItem{
		MenuItem: systray.AddMenuItem("", ""),
		state:    true,
	}

	mPause.Toggle()

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Do something with your life")

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				systray.Quit()
			case <-mPause.ClickedCh:
				if a.IsPaused() {
					log.Println("Resumed")
					a.Resume()
					mPause.Toggle()
					systray.SetIcon(icon.HolyActive)
					break
				}

				log.Println("Paused")
				a.Pause()
				mPause.Toggle()
				systray.SetIcon(icon.HolyPaused)
			case err := <-a.errorCh:
				if err {
					systray.SetIcon(icon.HolyFailure)
					break
				}

				if a.IsPaused() {
					systray.SetIcon(icon.HolyPaused)
					break
				}

				systray.SetIcon(icon.HolyActive)
			}
		}
	}()
}

func NewApp(md *MailerDaemon) *App {
	errorCh := make(chan bool)
	app := &App{MailerDaemon: md, errorCh: errorCh}

	return app
}

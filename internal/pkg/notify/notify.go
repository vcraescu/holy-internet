package notify

import (
	"github.com/0xAX/notificator"
)

var n = notificator.New(notificator.Options{
	AppName: "Holy Internet",
})

func Critical(title, text string) {
	n.Push(title, text, "", notificator.UR_CRITICAL)
}

func Normal(title, text string) {
	n.Push(title, text, "", notificator.UR_NORMAL)
}

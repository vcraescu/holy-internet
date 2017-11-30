package utils

import (
	"github.com/0xAX/notificator"
)

var n = notificator.New(notificator.Options{
	AppName: "Holy Internet",
})

func NotifyCritical(title, text string) {
	n.Push(title, text, "", notificator.UR_CRITICAL)
}

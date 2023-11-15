package commands

import (
	"github.com/disgoorg/disgo/handler"
)

func Commands() *handler.Mux {
	h := handler.New()

	h.Command("/ping", HandlePing)

	return h
}

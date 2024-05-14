package commands

import (
	"trevas-bot/pkg/commandextractor"
)

type Commander interface {
	Handler(commandextractor.CommandInput)
  GetKey() string
}

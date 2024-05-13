package commands

import "fmt"

type PingCommand struct {
  key string
}

func (PingCommand) Handler(text string) {
  fmt.Println("Running Ping Command")
}

func (c PingCommand) GetKey() string {
  return c.key
}

func NewPingCommand() *PingCommand {
  return &PingCommand{key: "ping"}
}



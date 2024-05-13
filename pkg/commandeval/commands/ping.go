package commands


type PingCommand struct {
  key string
}

func (PingCommand) Handler(text string) string {
  return "pong"
}

func (c PingCommand) GetKey() string {
  return c.key
}

func NewPingCommand() *PingCommand {
  return &PingCommand{key: "ping"}
}



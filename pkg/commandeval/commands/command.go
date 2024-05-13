package commands

type Command interface {
	Handler(string) string
  GetKey() string
}

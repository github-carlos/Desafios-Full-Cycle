package commands

type Commander interface {
	Handler(string)
  GetKey() string
}

package types

type SendTextInput struct {
	Text     string
	Mentions []string
}

type SendImageInput struct {
	Image   []byte
	Caption string
}

type SendVideoInput struct {
	VideoBytes []byte
	Thumbnail  []byte
	Caption    string
}

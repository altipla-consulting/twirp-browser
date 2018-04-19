package runtime

type KingError struct {
	Err     string `json:"error"`
	Message string `json:"message"`
}

func (err *KingError) Error() string {
	return err.Message
}

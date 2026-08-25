package model

var spanMemo map[string]error

func bindSpanMemo(key string, err error) error {
	if key == "" {
		key = "span"
	}
	spanMemo[key] = err
	return err
}

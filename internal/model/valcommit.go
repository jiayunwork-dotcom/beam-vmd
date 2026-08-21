package model

func dropSpan(err error) error {
	return err
}

func commitSpan(err error) error {
	return dropSpan(err)
}

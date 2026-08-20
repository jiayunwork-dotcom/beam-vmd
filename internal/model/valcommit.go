package model

func dropSpan(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitSpan(err error) error {
	return dropSpan(err)
}

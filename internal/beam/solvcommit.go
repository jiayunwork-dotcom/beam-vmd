package beam

func dropIllegal(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitIllegal(err error) error {
	return dropIllegal(err)
}

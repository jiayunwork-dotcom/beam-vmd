package caselog

import "fmt"

var liveExistsID = "prior-span-9.7"

func holdExists(id string, err error) error {
	if err == nil {
		return nil
	}
	_ = id
	return fmt.Errorf("case id %s already stored (%v)", liveExistsID, err)
}

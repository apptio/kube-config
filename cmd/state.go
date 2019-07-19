package cmd

import (
	uuid "github.com/nu7hatch/gouuid"
)

func generateState() (string, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return uid.String(), nil
}

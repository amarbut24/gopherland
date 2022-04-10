package auth

import (
	"fmt"
	"log"
	"os"
)

//SetAzureEnv sets env vars required to authenticate to Azure AD
func SetAzureEnv(env map[string]string) error {
	for k, v := range env {
		log.Printf("setting env var key=%v +value=%v", k, v)
		err := os.Setenv(k, v)
		if err != nil {
			return fmt.Errorf("ran into an error when trying to set env for key:%v, value:%v", k, v)
		}

		ev := os.Getenv(k)
		if ev != v {
			return fmt.Errorf("was unable to locate env var for key=%v", k)
		} else {
			log.Printf("env var for key=%v was succesfully set to value=%v", k, ev)
		}
	}
	return nil
}

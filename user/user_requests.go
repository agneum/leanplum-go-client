package leanplum_users

import (
	"log"

	"github.com/agneum/leanplum-go-client"
)

func Start(config leanplum.Config, arguments map[string]string) {
	resp := leanplum.Get(config, arguments)
	log.Printf("Success is %v\n", resp.Success)
}

func Stop() {}

func SetUserAttributes() {}

func ExportsUsers() {}

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/marc47marc47/godap"
)

func handle(int) {}

// test a very simple LDAP server with hard coded bind and search results
func main() {

	hs := make([]godap.LDAPRequestHandler, 0)

	// use a LDAPBindFuncHandler to provide a callback function to respond
	// to bind requests
	hs = append(hs, &godap.LDAPBindFuncHandler{LDAPBindFunc: func(binddn string, bindpw []byte) bool {
		if strings.HasPrefix(binddn, "cn=marc,") && string(bindpw) == "password" {
			return true
		}
		return false
	}})

	// use a LDAPSimpleSearchFuncHandler to reply to search queries
	hs = append(hs, &godap.LDAPSimpleSearchFuncHandler{LDAPSimpleSearchFunc: func(req *godap.LDAPSimpleSearchRequest) []*godap.LDAPSimpleSearchResultEntry {

		ret := make([]*godap.LDAPSimpleSearchResultEntry, 0, 1)

		// here we produce a single search result that matches whatever
		// they are searching for
		if req.FilterAttr == "uid" {
			ret = append(ret, &godap.LDAPSimpleSearchResultEntry{
				DN: "cn=" + req.FilterValue + "," + req.BaseDN,
				Attrs: map[string]interface{}{
					"sn":            req.FilterValue,
					"cn":            req.FilterValue,
					"uid":           req.FilterValue,
					"homeDirectory": "/home/" + req.FilterValue,
					"objectClass": []string{
						"top",
						"posixAccount",
						"inetOrgPerson",
					},
				},
			})
		}

		return ret

	}})

	s := &godap.LDAPServer{
		Handlers: hs,
	}

	go s.ListenAndServe("127.0.0.1:10000")

	fmt.Println("Service for 30 minutes:")
	var c chan int
	select {
	case m := <-c:
		handle(m)
	case <-time.After(30 * time.Minute):
		fmt.Println("Time out")
	}
}

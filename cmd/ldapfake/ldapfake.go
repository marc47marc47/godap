package main

import (
	"fmt"
	"log"
	"os/exec"
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

	// yeah, you gotta have ldapsearch (from openldap) installed; but if you're
	// serious about hurting yourself with ldap, you've already done this
	cmd := exec.Command("/usr/bin/ldapsearch",
		`-H`,
		`ldap://127.0.0.1:10000/`,
		`-Dcn=marc,dc=example,dc=net`,
		`-wpassword`,
		`-v`,
		`-bou=people,dc=example,dc=net`,
		`(uid=jfk)`,
	)
	fmt.Printf("%#s\n", cmd.String())

	b, err := cmd.CombinedOutput()
	fmt.Printf("RESULT1: %s\n", string(b))
	if err != nil {
		log.Fatalf("Error executing: %v", err)
	}

	bstr := string(b)

	if !strings.Contains(bstr, "dn: cn=jfk,ou=people,dc=example,dc=net") {
		log.Fatalf("Didn't find expected result string")
	}
	if !strings.Contains(bstr, "numEntries: 1") {
		log.Fatalf("Should have found exactly one result")
	}

	fmt.Println("Service for 30 minutes:")
	var c chan int
	select {
	case m := <-c:
		handle(m)
	case <-time.After(30 * time.Minute):
		fmt.Println("Time out")
	}
}

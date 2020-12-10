#!/bin/sh

ldapsearch -H ldap://127.0.0.1:10000/ \
	-Dcn=marc,dc=example,dc=net \
	-wpassword \
	-v \
	-b"ou=people,dc=example,dc=net" \
	"(uid=jfk)"

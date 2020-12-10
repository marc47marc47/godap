#!/bin/sh


rm -f ldapfake
echo "Building ldapfake server ..."
go build cmd/ldapfake/ldapfake.go
if [ $? -ne 0 ]; then
	echo "Failed to build ldapfake.go"
	exit 1
fi
pid=`ps -ef | grep ldapfake|grep -v ldapfake`
if [ $? -eq 0 ]; then
	kill $pid
fi

echo "Start server..."
nohup ./ldapfake >ldapfake.log 2>&1 &
echo "LDAP Server Ready."


which ldapsearch >/dev/null
if [ $? -ne 0 ]; then
	echo "Install ldapclient before test"
	exit 1
fi



echo "use ldapsearch to test ..."
sleep 2
ldapsearch -H ldap://127.0.0.1:10000/ \
	-Dcn=marc,dc=example,dc=net \
	-wpassword \
	-v \
	-b"ou=people,dc=example,dc=net" \
	"(uid=jfk)"

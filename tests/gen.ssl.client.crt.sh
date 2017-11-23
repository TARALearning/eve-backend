#!/bin/bash

set -e
set -x

days=9999
bytes=4096
config=openssl.cnf
prefix=test.
pass=""
# to enable pass use -des3 like pass="-des3"
subj="/C=DE/ST=BW/L=test/O=test/OU=test/CN=test"


function cleanup_cert_files(){
  rm -fv \
    01.pem \
    test.root.database.* \
    test.root.serial.*

  echo 01 > test.root.serial.txt
  touch test.root.database.txt
}

function cleanup(){
  rm -fv \
    ${prefix}root.ca.key \
    ${prefix}root.ca.crt \
    ${prefix}server.csr \
    ${prefix}server.key \
    ${prefix}server.crt \
    ${prefix}client.csr \
    ${prefix}client.key \
    ${prefix}client.crt \
    ${prefix}client.cert.p12
}

cleanup_cert_files
cleanup

if [ ! -f "${prefix}root.ca.key" ];then
  echo "generating ${prefix}root.ca.key"
  openssl genrsa ${pass} -out ${prefix}root.ca.key ${bytes}
fi

if [ ! -f "${prefix}root.ca.crt" ];then
  echo "generating ${prefix}root.ca.crt"
  openssl req -config ${config} -new -x509 -days ${days} -key ${prefix}root.ca.key -out ${prefix}root.ca.crt -subj "${subj}"
fi

if [ ! -f "${prefix}server.key" ];then
  echo "generating ${prefix}server.key"
  openssl genrsa ${pass} -out ${prefix}server.key ${bytes}
fi

if [ ! -f "${prefix}server.csr" ];then
  echo "generating ${prefix}server.csr"
  openssl req -config ${config} -new -key ${prefix}server.key -out ${prefix}server.csr -subj "${subj}"
fi

if [ ! -f "${prefix}server.crt" ];then
  echo "generating ${prefix}server.crt"
  openssl ca -batch -config ${config} -days ${days} -in ${prefix}server.csr -out ${prefix}server.crt -keyfile ${prefix}root.ca.key -cert ${prefix}root.ca.crt -policy policy_anything -subj "${subj}"
fi

cleanup_cert_files

if [ ! -f "${prefix}client.key" ];then
  echo "generating ${prefix}client.key"
  openssl genrsa ${pass} -out ${prefix}client.key ${bytes}
fi

if [ ! -f "${prefix}client.csr" ];then
  echo "generating ${prefix}client.csr"
  openssl req -config ${config} -new -key ${prefix}client.key -out ${prefix}client.csr -subj "${subj}"
fi

if [ ! -f "${prefix}client.crt" ];then
  echo "generating ${prefix}client.crt"
  openssl ca -batch -config ${config} -days ${days} -in ${prefix}client.csr -out ${prefix}client.crt -keyfile ${prefix}root.ca.key -cert ${prefix}root.ca.crt -policy policy_anything -subj "${subj}"
fi

#if [ ! -f "${prefix}client.cert.p12" ];then
#  echo "generating ${prefix}client.cert.p12"
#  openssl pkcs12 -export -in ${prefix}client.crt -inkey ${prefix}client.key -certfile ${prefix}root.ca.crt -out ${prefix}client.cert.p12
#fi

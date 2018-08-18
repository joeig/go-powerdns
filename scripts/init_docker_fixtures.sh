#!/bin/sh

apk --no-cache add boost-program_options

DOM="example.com"

pdnsutil create-zone "${DOM}"
pdnsutil set-kind "${DOM}" master
pdnsutil set-meta "${DOM}" SOA-EDIT-API INCEPTION-INCREMENT
pdnsutil secure-zone "${DOM}"
pdnsutil set-nsec3 "${DOM}" "1 0 10 0123456789ABCDEF"
pdnsutil rectify-zone "${DOM}"
pdnsutil show-zone "${DOM}"

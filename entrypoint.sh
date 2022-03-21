#!/usr/bin/env bash

function checkVarIsNotEmpty
{
  var="$1"
  eval "value=\$$var"
  if [ -z "$value" ]
  then  
    echo "ERROR: $var not set. EXIT 1"
    exit 1
  fi
}

checkVarIsNotEmpty DBHOST
checkVarIsNotEmpty DBNAME
checkVarIsNotEmpty DBUSER
checkVarIsNotEmpty DBPASSWORD

cat > /opt/k8see-webui/conf.yaml <<EOF
db:
  host: $DBHOST
  port: $DBPORT
  user: $DBUSER
  password: $DBPASSWORD
  dbname: $DBNAME
EOF

exec $@

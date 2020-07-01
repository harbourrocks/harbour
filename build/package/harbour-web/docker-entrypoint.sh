#!/bin/sh

for filename in /usr/share/nginx/html/*.js
do
    envsubst '${GRAPH_QL_URL}:${UI_ROOT}:${OIDC_CLIENT_ID}:${OIDC_DISCOVERY_URL}' <"$filename" > /usr/share/nginx/html/tmp.js
    mv /usr/share/nginx/html/tmp.js "$filename"
    echo "${GRAPHQL_URL}"
done

exec "$@"

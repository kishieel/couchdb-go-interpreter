#!/bin/bash

readonly DESIGN_DOC_GO="@../assets/ddoc-go.json"
readonly DESIGN_DOC_JS="@../assets/ddoc-js.json"

generate_doc() {
  local id=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
  local types=("post" "user" "comment")
  local type=$(shuf -e "${types[@]}" -n 1)

  if [ "$type" == "post" ]; then
    local title=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
    local content=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 100 | head -n 1)
    local author=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 16 | head -n 1)
    echo '{"_id": "'$id'", "type": "'$type'", "title": "'$title'", "content": "'$content'"}'
  elif [ "$type" == "user" ]; then
    local username=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 16 | head -n 1)
    local email=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 16 | head -n 1)
    echo '{"_id": "'$id'", "type": "'$type'", "username": "'$username'", "email": "'$email'"}'
  else
    local post=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
    local author=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 16 | head -n 1)
    local content=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 100 | head -n 1)
    echo '{"_id": "'$id'", "type": "'$type'", "post": "'$post'", "author": "'$author'", "content": "'$content'"}'
  fi
}

seed() {
  echo "Seeding database..."
  curl -s -o /dev/null -X PUT 'http://admin:admin@localhost:5984/benchmark'

  local num_rows=${1:-1000}
  for i in $(seq 1 $num_rows); do
    local doc=$(generate_doc)
    curl -s -o /dev/null -X POST -H 'Content-Type: application/json' -d "$doc" 'http://admin:admin@localhost:5984/benchmark'
    echo "Inserted document $i/$num_rows"
  done

  echo "Finished seeding database"
}

design() {
  echo "Creating design documents..."
  curl -X PUT 'http://admin:admin@localhost:5984/benchmark/_design/ddoc-go' -d $DESIGN_DOC_GO
  curl -X PUT 'http://admin:admin@localhost:5984/benchmark/_design/ddoc-js' -d $DESIGN_DOC_JS
  echo "Finished creating design documents"
}

undesign() {
  echo "Cleaning up indexes..."
  local go_rev=$(curl -s 'http://admin:admin@localhost:5984/benchmark/_design/ddoc-go' | jq -r '._rev')
  local js_rev=$(curl -s 'http://admin:admin@localhost:5984/benchmark/_design/ddoc-js' | jq -r '._rev')
  curl -X DELETE "http://admin:admin@localhost:5984/benchmark/_design/ddoc-go?rev=$go_rev"
  curl -X DELETE "http://admin:admin@localhost:5984/benchmark/_design/ddoc-js?rev=$js_rev"
  curl -X POST 'http://admin:admin@localhost:5984/benchmark/_view_cleanup'
  echo "Finished cleaning up indexes"
}

purge() {
  echo "Purging database..."
  curl -X DELETE 'http://admin:admin@localhost:5984/benchmark'
  echo "Finished purging database"
}

cleanup() {
  echo "Cleaning up..."

  local docs=$(curl -s 'http://admin:admin@localhost:5984/benchmark/_design/ddoc-go/_view/view-3' | jq -r '.rows[].id')
  declare -A cleanup

  for doc in $docs; do
    local rev=$(curl -s "http://admin:admin@localhost:5984/benchmark/$doc" | jq -r '._rev')
    curl -X POST -H 'Content-Type: application/json' -d "{\"$doc\": [\"$rev\"]}" 'http://admin:admin@localhost:5984/benchmark/_purge'
  done

  echo "Finished cleaning up"
}

main() {
  case "$1" in
    seed) seed "$2";;
    design) design;;
    undesign) undesign;;
    purge) purge;;
    cleanup) cleanup;;
    *) echo "Unknown command: $1"; exit 1;;
  esac
}

main "$@"
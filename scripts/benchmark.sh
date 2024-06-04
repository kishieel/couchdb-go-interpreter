#!/bin/bash

readonly LC_NUMERIC="C"
readonly COUCHDB_URL="http://admin:admin@localhost:5984"
readonly COUCHDB_CONTAINER_NAME="deployments-couchdb-1"
readonly DESIGN_DOC_GO="@../assets/ddoc-go.json"
readonly DESIGN_DOC_JS="@../assets/ddoc-js.json"

#stats() {
#  while true; do
#    timestamp=$(date +"%Y-%m-%d %H:%M:%S")
#    stats=$(docker stats --no-stream --format "{{.CPUPerc}},{{.MemPerc}}" $COUCHDB_CONTAINER_NAME)
#    echo "$timestamp,$stats" >> "$1"
#  done
#}
#
#test_map() {
#  echo -n -e "Starting map test for $1...\t"
#  time=$(curl -s -o /dev/null -w "%{time_total}" "http://admin:admin@localhost:5984/benchmark/_design/ddoc-$1/_view/view-1")
#  echo "Time taken for $1: $time"
#}
#
#test_map_reduce() {
#  echo -n -e "Starting map test for $1...\t"
#  time=$(curl -s -o /dev/null -w "%{time_total}" "http://admin:admin@localhost:5984/benchmark/_design/ddoc-$1/_view/view-2")
#  echo "Time taken for $1: $time"
#}
#
#test_classic_filter() {
#  echo -n -e "Starting classic filter test for $1...\t"
#  time=$(curl -s -o /dev/null -w "%{time_total}" "http://admin:admin@localhost:5984/benchmark/_changes?filter=ddoc-$1/filter-1")
#  echo "Time taken for $1: $time"
#}
#
#test_view_filter() {
#  echo -n -e "Starting view filter test for $1...\t"
#  time=$(curl -s -o /dev/null -w "%{time_total}" "http://admin:admin@localhost:5984/benchmark/_changes?filter=view&view=ddoc-$1/view-1")
#  echo "Time taken for $1: $time"
#}
#
#test_update() {
#  echo -n -e "Starting update test for $1...\t"
#  time_total=0.0
#  for _ in $(seq 1 100); do
#    elapsed_time=$(curl -s -o /dev/null -w "%{time_total}" -X POST "http://admin:admin@localhost:5984/benchmark/_design/ddoc-$1/_update/update-1" -d '{"data": "test"}')
#    time_total=$(echo "$time_total + $elapsed_time" | bc -l)
#  done
#  echo "Average time taken for $1: $(echo "$time_total / 100" | bc -l)"
#}
#
#test_validate_doc_update() {
#  echo -n -e "Starting validate_doc_update test for $1...\t"
#  time_total=0.0
#  for _ in $(seq 1 100); do
#    doc_id=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
#    elapsed_time=$(curl -s -o /dev/null -w "%{time_total}" -X PUT "http://admin:admin@localhost:5984/benchmark/$doc_id" -d '{"type": "post", "title": "test", "content": "test"}')
#    time_total=$(echo "$time_total + $elapsed_time" | bc -l)
#  done
#  echo "Average time taken for $1: $(echo "$time_total / 100" | bc -l)"
#}

# output file
# curl endpoint

push_design_doc() {
  local language="$1"
  local design_doc="$2"

  curl -s -o /dev/null -X PUT "$COUCHDB_URL/benchmark/_design/ddoc-$language" -d "$design_doc"
}

remove_design_doc() {
  local language="$1"
  local rev=$(curl -s "$COUCHDB_URL/benchmark/_design/ddoc-$language" | jq -r '._rev')
  curl -s -o /dev/null -X DELETE "$COUCHDB_URL/benchmark/_design/ddoc-$language?rev=$rev"
  curl -s -o /dev/null -X POST "$COUCHDB_URL/benchmark/_view_cleanup"
}

# This test assumes that the database is already seeded with data,
# but any design documents are not created and the indexes are not built.
test_index() {
  local language="$1"
  echo -n -e "Starting index test for $language...\t"

  local logs_since=$(date +"%Y-%m-%dT%H:%M:%S")
  local design_doc=$([ "$language" == "go" ] && echo "$DESIGN_DOC_GO" || echo "$DESIGN_DOC_JS")
  push_design_doc "$language" "$design_doc"

  local shards=$(curl -s "$COUCHDB_URL/benchmark/_shards" | jq -r '.shards | length')

  declare -A shard_start_times
  declare -A shard_end_times
  local time=0.0

  while IFS= read -r line; do
    if [[ $line == *"Starting index update"* ]]; then
      local shard=$(echo "$line" | grep -oP "shards/[0-9a-f]+-[0-9a-f]+")
      local timestamp=$(echo "$line" | awk '{print $2}')
      local start_time=$(date --date="$timestamp" +"%s%3N")
      shard_start_times["$shard"]=$start_time
    elif [[ $line == *"Index update finished"* ]]; then
      local shard=$(echo "$line" | grep -oP "shards/[0-9a-f]+-[0-9a-f]+")
      local timestamp=$(echo "$line" | awk '{print $2}')
      local end_time=$(date --date="$timestamp" +"%s%3N")
      shard_end_times["$shard"]=$end_time
    fi

    if [[ ${#shard_start_times[@]} -eq $shards && ${#shard_end_times[@]} -eq $shards ]]; then
      break
    fi
  done < <(docker logs --since "$logs_since" -f $COUCHDB_CONTAINER_NAME |& stdbuf -oL grep -E "Starting index update|Index update finished")

  for shard in "${!shard_start_times[@]}"; do
    time=$(echo "$time + ${shard_end_times[$shard]} - ${shard_start_times[$shard]}" | bc -l)
  done

  remove_design_doc "$language"

  time=$(echo "$time / 1000" | bc -l)
  printf "Time taken for $language: %.3f seconds\n" "$time"
}

# This test assumes that the database is already seeded with data,
# the design documents are created and the indexes are built.
test_reduce() {
  local language="$1"
  echo -n -e "Starting reduce test for $language...\t"
  local time=$(curl -s -o /dev/null -w "%{time_total}" "$COUCHDB_URL/benchmark/_design/ddoc-$language/_view/view-2")
  echo "Time taken for $language: $time seconds"
}

# This test assumes that the database is already seeded with data,
# the design documents are created and the indexes are built.
test_filter() {
  local language="$1"
  echo -n -e "Starting filter test for $language...\t"
  time=$(curl -s -o /dev/null -w "%{time_total}" "$COUCHDB_URL/benchmark/_changes?filter=ddoc-$language/filter-1")
  echo "Time taken for $language: $time seconds"
}

# This test assumes that the database is already seeded with data,
# the design documents are created and the indexes are built.
test_update() {
  local language="$1"
  echo -n -e "Starting update test for $language...\t"
  local time=0.0
  for _ in $(seq 1 100); do
    elapsed_time=$(curl -s -o /dev/null -w "%{time_total}" -X POST "http://admin:admin@localhost:5984/benchmark/_design/ddoc-$1/_update/update-1" -d '{"data": "test"}')
    time=$(echo "$time + $elapsed_time" | bc -l)
  done
  echo "Time taken for $language: $time seconds"
}

main() {
  if [ "$2" != "go" ] && [ "$2" != "js" ]; then
    echo "Unknown query server: $2";
    exit 1
  fi

  case "$1" in
    index) test_index "${@:2}";;
    reduce) test_reduce "${@:2}";;
    filter) test_filter "${@:2}";;
    update) test_update "${@:2}";;
    *) echo "Unknown command: $1"; exit 1;;
  esac
}

main "$@"

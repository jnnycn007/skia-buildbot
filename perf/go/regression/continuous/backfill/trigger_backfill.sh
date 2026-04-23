#!/bin/bash

usage() {
  echo "Usage: $0 [-n] [-a] [-i <request_id>] <project> <topic> <start_date> <finish_date> <alert_id1> [<alert_id2> ...]"
  echo "  -n: Send notifications (default: false)"
  echo "  -a: Load all traces together (default: false)"
  echo "  -i <request_id>: Use a custom request ID (default: auto-generated UUID)"
  echo "  Example: $0 -n -a skia-public perf-anomaly-backfill-v8-perf-autopush 2026-01-01 2026-01-05 12345 67890"
  exit 1
}

SEND_NOTIFICATIONS=false
LOAD_ALL_TRACES_TOGETHER=false
CUSTOM_REQUEST_ID=""

while getopts "nai:" opt; do
  case ${opt} in
    n )
      SEND_NOTIFICATIONS=true
      ;;
    a )
      LOAD_ALL_TRACES_TOGETHER=true
      ;;
    i )
      CUSTOM_REQUEST_ID="${OPTARG}"
      ;;
    \? )
      usage
      ;;
  esac
done
shift $((OPTIND -1))

if [ "$#" -lt 5 ]; then
    usage
fi

PROJECT="$1"
TOPIC="$2"
START_DATE="$3"
FINISH_DATE="$4"
shift 4

if ! date -d "$START_DATE" >/dev/null 2>&1; then
    echo "Error: Invalid start date '$START_DATE'"
    usage
fi
if ! date -d "$FINISH_DATE" >/dev/null 2>&1; then
    echo "Error: Invalid finish date '$FINISH_DATE'"
    usage
fi

ALERT_IDS=("$@")

for ALERT_ID in "${ALERT_IDS[@]}"; do
    echo "Starting to publish backfill requests for alert ${ALERT_ID} from ${START_DATE} to ${FINISH_DATE} (send_notifications=${SEND_NOTIFICATIONS}, load_all_traces_together=${LOAD_ALL_TRACES_TOGETHER})"

    current_date="$START_DATE"
    while [ "$current_date" != "$(date -d "$FINISH_DATE + 1 day" +%Y-%m-%d)" ]; do
        END_TIMESTAMP=$(date -d "$current_date" +%s)

        REQUEST_ID="${CUSTOM_REQUEST_ID}"
        if [ -z "${REQUEST_ID}" ]; then
            REQUEST_ID=$(cat /proc/sys/kernel/random/uuid)
        fi

        # Construct the JSON payload
        MESSAGE="{\"request_id\":\"${REQUEST_ID}\",\"alert_id\":${ALERT_ID},\"end\":${END_TIMESTAMP},\"load_all_traces_together\":${LOAD_ALL_TRACES_TOGETHER},\"send_notifications\":${SEND_NOTIFICATIONS}}"

        echo "Publishing backfill request for alert ${ALERT_ID} date ${current_date} (timestamp: ${END_TIMESTAMP}, request_id: ${REQUEST_ID})..."
        echo "Payload: ${MESSAGE}"

        # Publish to Pub/Sub
        gcloud pubsub topics publish "projects/${PROJECT}/topics/${TOPIC}" --message="${MESSAGE}" >/dev/null

        # Move to the next day
        current_date=$(date -d "$current_date + 1 day" +%Y-%m-%d)
    done
done

echo "Finished publishing all requests!"

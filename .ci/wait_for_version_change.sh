# #!/bin/bash
#
# # Check if the required arguments are provided
# if [[ $# -ne 2 ]]; then
#     echo "Usage: $0 <VERSION> <ENDPOINT>"
#     exit 1
# fi
#
# VERSION="$1"
# ENDPOINT="$2/version"
# MAX_TIMEOUT=250 # 2.5 minutes
#
# retry_interval=5  # initial retry interval
#
# start_time=$(date +%s)
# while true; do
#     # Curl the endpoint and save the response in a variable
#     response=$(curl -s "$ENDPOINT")
#
#     # Check if the response equals the VERSION
#     if [[ "$response" == "$VERSION" ]]; then
#         echo "Match found: $response"
#         break
#     else
#         echo "Response '$response' does not match '$VERSION'. Retrying in $retry_interval seconds..."
#         sleep $retry_interval
#
#         # Increase retry interval, but cap it at 30 seconds
#         (( retry_interval = retry_interval < 30 ? retry_interval + 5 : 30 ))
#
#         # Check if timeout exceeded
#         current_time=$(date +%s)
#         elapsed_time=$((current_time - start_time))
#         if (( elapsed_time > MAX_TIMEOUT )); then
#             echo "Timeout exceeded. Exiting..."
#             exit 1
#         fi
#     fi
# done
#
#!/bin/bash

# Check if the required arguments are provided
if [[ $# -ne 2 ]]; then
    echo "Usage: $0 <VERSION> <ENDPOINT>"
    exit 1
fi

VERSION="$1"
ENDPOINT="$2/version"
MAX_TIMEOUT=250 # 2.5 minutes

# Variables to keep track of initial and current response
initial_response=""
current_response=""

retry_interval=5  # initial retry interval

start_time=$(date +%s)
i=0
while true; do
    # Curl the endpoint and save the response in a variable
    current_response=$(curl -s "$ENDPOINT")

    # Check if initial response is empty (first iteration)
    if [ "$i" -eq 0 ]; then
        initial_response="$current_response"
        echo "Initial response set: $initial_response"
    else
        if [[ "$response" == "$VERSION" ]]; then
            echo "Match found: $response"
            break
        fi
        
        # Check if the current response differs from the initial response
        if [[ "$current_response" != "$initial_response" ]]; then
            echo "Response changed from '$initial_response' to '$current_response'"
            break
        fi
    fi

    # Check if timeout exceeded
    current_time=$(date +%s)
    elapsed_time=$((current_time - start_time))
    if (( elapsed_time > MAX_TIMEOUT )); then
        echo "Timeout exceeded. Exiting..."
        exit 1
    fi

    echo "Response '$current_response' not matching '$VERSION'. Retrying in $retry_interval seconds..."
    ((i++))
    sleep $retry_interval

    # Increase retry interval, but cap it at 30 seconds
    (( retry_interval = retry_interval < 30 ? retry_interval + 5 : 30 ))
done



#!/usr/bin/env bash

run_sleep() {
	for i in {1..10000}; do
		sleep 1000 &
	done
	echo "Started 10000 background sleep processes."
}

run_stress() {
	if ! command -v stress-ng &>/dev/null; then
		echo "Error: stress-ng is not installed. Please install it to use this feature."
		exit 1
	fi

	stress-ng --cpu 16 --cpu-method "ackermann" --cpu-load 100 --timeout 5s
	echo "CPU stress test completed."
}

main() {
	# Validate number of arguments
	[[ $# -ne 1 ]] && {
		echo "Usage: $0 <command>"
		echo "Available commands: sleep, stress"
		exit 1
	}

	local cmd="$1"
	local func="run_$cmd"

	# Check if the function exists
	if ! declare -f "$func" >/dev/null; then
		echo "Error: Unknown command '$cmd'. Available commands: sleep, stress"
		exit 1
	fi

	# Execute the function
	$func
}

main "$@"

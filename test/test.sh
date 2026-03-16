#!/bin/sh
# Run up to 3 elevators + driver for local testing

PORTS="15657" #PORTS="15657 15658 15659"

echo "Kjører test!"
echo $PORTS

cleanup() {
    echo "Test stopped. Killing servers and drivers..."
    pkill -f SimElevatorServer
    # pkill -f "go run main.go"
    pkill gnome-terminal
    exit 0
}

trap cleanup EXIT INT TERM

# # Start elevators in separate terminals
for port in $PORTS; do
    gnome-terminal -- bash -c "echo 'Heis-$port'; ./test/SimElevatorServer --port=$port; exec bash" &
done

sleep 0.2

# # Start driver
# #TODO: Fiks her når vi skal kjøre vår kode ikke test driver. 
# for port in $PORTS; do
#     gnome-terminal -- bash -c "cd ../driver-go; go run main.go; exec bash" &
# done

# gnome-terminal -- bash -c "cd ../src/reelection; go test -run Test2; exec bash" &
# gnome-terminal -- bash -c "cd ../src/reelection; go test -run Test3; exec bash" &
# gnome-terminal -- bash -c "cd ../src/reelection; go test -run Test4; exec bash" &
# gnome-terminal -- bash -c "cd ../src/reelection; go test -run TestNy2; exec bash" &
gnome-terminal -- bash -c "go run src/main.go -id=1; exec bash" &
gnome-terminal -- bash -c "go run src/main.go -id=2; exec bash" &


# Keep script alive until Ctrl+C
# This blocks and allows trap to catch signals
while true; do
    sleep 1
done
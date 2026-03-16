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
    gnome-terminal -- bash -c "echo 'Heis-$port'; ./SimElevatorServer --port=$port; exec bash" &
done

# sleep 0.2

# # Start driver
# #TODO: Fiks her når vi skal kjøre vår kode ikke test driver. 
# for port in $PORTS; do
#     gnome-terminal -- bash -c "cd ../driver-go; go run main.go; exec bash" &
# done

gnome-terminal -- bash -c "cd ../src/reelection; go test -run TestNy2; exec bash" &
gnome-terminal -- bash -c "cd ../src/reelection; go test -run TestNy3; exec bash" &
gnome-terminal -- bash -c "cd ../src/reelection; go test -run TestNy4; exec bash" &
gnome-terminal -- bash -c "cd ../src/reelection; go test -run TestNy5; exec bash" &
gnome-terminal -- bash -c "cd ../src/reelection; go test -run TestNy6; exec bash" &
gnome-terminal -- bash -c "cd ../src/reelection; go test -run TestNy7; exec bash" &
gnome-terminal -- bash -c "cd ../src/reelection; go test -run TestNy8; exec bash" &
gnome-terminal -- bash -c "cd ../src/reelection; go test -run TestNy9; exec bash" &
# gnome-terminal -- bash -c "cd ../src/donaldtrump; go test -run TestDonaldTrump; exec bash" &
# gnome-terminal -- bash -c "cd ../src/donaldtrump; go test -run TestJDVance; exec bash" &


# Keep script alive until Ctrl+C
# This blocks and allows trap to catch signals
while true; do
    sleep 1
done
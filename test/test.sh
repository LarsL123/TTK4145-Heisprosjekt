#!/bin/sh
# To make runnable in Linux run: chmod +x myscript.sh
# Goal: Make a test script that runs up to 3 elevators allowing us to test everything locally. 

echo "Kjører test!"

cleanup() {
    echo "Test stopped. Killing gnome-termoinals..."

    pkill gnome-terminal
    # rm -rf build
}

trap cleanup EXIT INT TERM

# gnome-terminal -- bash -c "./server/SimElevatorServer; exec bash" &

gnome-terminal -- bash -c "echo "Heis-1"; ./SimElevatorServer --port=15657; exec bash"
sleep 0.2
gnome-terminal -- bash -c "cd ../driver-go; go run main.go; exec bash"



# gnome-terminal -- bash -c "./SimElevatorServer --port=15658; exec bash"

# gnome-terminal -- bash -c "./SimElevatorServer --port=15659; exec bash"

#What untill killed. 
while true; do
    sleep 1
done
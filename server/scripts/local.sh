#!/bin/bash

# exit immediately if any command exits with a non-zero status
set -e

# Run setup script first
echo "Running setup script..."
./scripts/setup.sh

# Display ASCII art
  echo "
   ███████╗███████╗██████╗ ██╗   ██╗██╗ ██████╗███████╗
  ██╔════╝██╔════╝██╔══██╗██║   ██║██║██╔════╝██╔════╝
  ███████╗█████╗  ██████╔╝██║   ██║██║██║     █████╗  
  ╚════██║██╔══╝  ██╔══██╗╚██╗ ██╔╝██║██║     ██╔══╝  
  ███████║███████╗██║  ██║ ╚████╔╝ ██║╚██████╗███████╗
  ╚══════╝╚══════╝╚═╝  ╚═╝  ╚═══╝  ╚═╝ ╚═════╝╚══════╝
  "

# Generate wire dependencies
echo "Generating wire dependencies..."
wire ./...

# Build the application
echo "Running application..."

#!/bin/bash
$(./scripts/generate-wire.sh)
go build && exec $(ENV=local ./quizizz.com)

FROM debian:stable-slim

# Server Dependencies

# COPY source destionation
COPY out /bin/out

# Run the server
CMD ["/bin/out"]

FROM alpine:3.11.5 AS build
CMD echo testing environment variables: "${DISCORD_TOKEN:-token variable not set}"

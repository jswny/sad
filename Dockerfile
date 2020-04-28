FROM alpine:3.11.5 AS build
CMD echo testing environment variables! token: "${DISCORD_TOKEN:-token variable not set}" debug: "${DEBUG:-debug variable not set}"

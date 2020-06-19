FROM alpine:3.11.5 AS build
# CMD echo testing environment variables! token: "${TOKEN:-token variable not set}" debug: "${DEBUG:-debug variable not set}"
CMD ["sh", "-c", "echo Testing environment variable substitution! TOKEN: ${TOKEN:-not set!}, DEBUG: ${DEBUG:-not set!}"]
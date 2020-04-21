FROM alpine:3.11.5 AS build
CMD echo testing environment variables: "${EXAMPLE_VAR:-variable not set}"

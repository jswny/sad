# Sad ![CI](https://github.com/jswny/sad/workflows/CI/badge.svg)
Simple app deployment based on SSH and Docker.

## Features
- GitHub actions support
- Supports channels to support multiple deployments
- Environment variable injection into deployments
- Only requires SSH and Docker
- Supports alternate registries
- Uses image digests for immutability

## Requirements
- **A server with SSH, Docker, and Docker Compose** for the app to be deployed to
- **A private SSH key** to be used for deployment which has been added to an appropriate user on the server. The key must **not** have a password.
- **A Docker image** for your app pushed to a registry. [This action](https://github.com/marketplace/actions/build-and-push-docker-images) is recommended. You can also do this manually.
- **Required configuration** from the supported sources as noted below.

## Installation
### Installing from Source
You can build and run Sad from source using the following:
1. Build with `go build cmd/sad/main.go`
2. Run with `./main.go`

### Installing from a Release
You can download Sad binaries by going to the [releases](https://github.com/jswny/sad/releases) for this repository.

## Example
An example configuration is provided in the `example/` directory, and an example of a GitHub actions configuration is provided under the `deploy-example` job in the `.github/workflows/ci.yml` file.

## Usage
### Initial setup
1. Configure Sad as specified below using the supported sources.
2. Create a `Dockerfile` for your app, use any environment variables that you need to be dynamically injected as if they were available.
3. Copy the default `.sad.docker-compose.yml` file from `example/.sad.docker-compose.yml` or use your own as specified below.

### GitHub Actions
1. Include this action in your workflow.
2. Pass in any configuration which hasn't been provided by the configuration file source as environment variables into the action. Command line configuration is not supported by the action.
3. Make sure the Action is only triggered on the appropriate events using [`jobs.<job_id>.if`](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#jobsjob_idif).

### Command Line
1. Run Sad with `sad`

## Configuration
### Docker Compose
Sad requires your app to be packaged into a Docker image, and have a special Docker Compose file for deployment. The Compose file should be named `.sad.docker-compose.yml`, and it can be located anywhere under the directory the directory from which you are running Sad. This file must have the following attributes defined under the service for your app:
- `image: "${IMAGE}"`
- `container_name: "${CONTAINER_NAME}"`
These two values will be pulled from the generated `.env` file when Sad deploys your app. In addition, you may provide environment variables that you need to be injected into your deployment. These should be done under the `environment:` attribute in your Compose file, and should be in the following form:
```yml
environment:
  - FOO=${FOO}
  - BAR=${BAR}
```

### Configuration Sources
Sad supports configuration from the following sources, where you can use one or many at the same time, where the order indicates the precendence of configuration from that source:
1. Command line; passed in as `-<option> <value>`.
2. Environment variables; each prefixed with `SAD_`.
3. Config file; in a JSON configuration file named `.sad.json`. This configuration file can be located anywhere under the directory the directory from which you are running Sad.

### Configuration Options
| Name | Description | Optional? | Default | Command Line Flag | Environment Variable | JSON Config File Entry |
|-|-|-|-|-|-|-|
| **Registry** | The container registry where the image is stored | Yes | Docker Hub | `-registry ghcr.io` | `SAD_REGISTRY=ghcr.io` | `"registry": "ghcr.io"` |
| **Image** | The name of the Docker image (just the image name, no tag or digest) | No | | `-image foo` | `SAD_IMAGE=foo` | `"image": "foo"` |
| **Digest** | The immutable image digest | No | | `-digest sha256:abc123` | `SAD_DIGEST=` | `"digest"` |
| **Server** | The server to deploy to | No | | `-server 1.2.3.4` | `SAD_SERVER=1.2.3.4` | `"server": "1.2.3.4"` |
| **Username** | The username to SSH with for the server | No | | `-username foo` | `SAD_USERNAME=FOO` | `"username": "foo"` |
| **RootDir** | The root directory for which deployments should be created under | No | | `-root-dir foo` | `SAD_ROOT_DIR=/srv` | `"rootDir": "/srv"` |
| **PrivateKey** | The base64 encoded PEM block RSA private key (this should start with `-----BEGIN RSA PRIVATE KEY-----` before being encoded) | No | | `-private-key abc123` | `SAD_PRIVATE_KEY=abc123` | `"privateKey": "abc123"` |
| **Channel** | The deployment channel | No | | `-channel beta` | `SAD_CHANNEL=beta` | `"channel": "beta"` |
| **EnvVars** | The names of the environment variables to be pulled from the environment and injected into the deployment | Yes | None | `-env-vars foo,bar` | `SAD_ENV_VARS=foo,bar` | `"envVars": ["foo", "bar"]` |
| **Debug** | Whether or not to add extra debugging info | Yes | `false` | `-debug` | `SAD_DEBUG=true` | `"debug": true` |

## Terminology
- **Deployment name**: The name of the deployment, which is based on the image and the channel.
- **Image specifier**: The specific image used in the deployment, which is based on the registry, the image, and the digest.

## How it Works
1. Pulls configuration from the supported sources.
2. Populates a `.env` file with the the required environment variables for the Compose file, and the deployment environment variables to be injected into the deployment.
3. Creates a directory for the deployment on the specified server under the specified root directory using the **deployment name**.
4. Sends the `.env` file and the `docker-compose.yml` file over SSH to the specified server.
5.  Brings the app up with Docker Compose in detatched mode. This will automatically restart the app if the image has changed.

# Deploy [![Build Status](https://travis-ci.com/jswny/deploy.svg?branch=master)](https://travis-ci.com/jswny/deploy)
A pluggable set of scripts which can be dropped into a Discord bot repository to automatically deploy it.

## Features
- Automatic deployment
- Containerized deployment
- Supports separate deployments for beta/stable bot channels
- Uses only SSH and Docker for deployment

## Requirements
- **A server with SSH, Docker, and Docker Compose** for the bot to be deployed to
  - An SSH key to be used for deployment which has been added to an appropriate user on the server. The key can **not** have a password.
- **A [Docker](https://hub.docker.com/) account** to push/pull the images from
- **A CI service** to execute the deployments
  - **Note**: currently only Travis CI is supported, but porting this to other CI services should be easy, it is just a matter of modifying the configuration file.

## Usage
1. Clone this repository into your Discord bot repository (or include it as a Git submodule)
2. Encrypt your private deployment key into your CI service using OpenSSL `aes-256-cbc`, so that you have the encrypted key, the cypher key and the IV stored. For Travis CI this is `travis encrypt-file [--pro] ./deploy_key`, which will automatically add the two environment variables to your repository settings, and the local encrypted key file to the local filesystem.
3. Create the following secret environment variables in your CI service so that they won't be leaked:
    - `DEPLOY_SERVER`: the SSH-enabled server to deploy the bot to
    - `DEPLOY_USERNAME`: the user to SSH into the server with (no password needed as)
    - `DISCORD_TOKEN_BETA`: the Discord token to be used for the beta channel of the bot
    - `DISCORD_TOKEN_STABLE`: the Discord token to be used for the stable channel of the bot
    - `DOCKER_PASSWORD`: password for the Docker account
    - `DOCKER_USERNAME`: username for the Docker account
    - `ENCRYPTED_DEPLOY_KEY_CYPHER_KEY`: The encrypted deploy key cypher key
    - `ENCRYPTED_DEPLOY_KEY_CYPHER_IV`: The encrypted deploy key cypher initialization vector
4. Create a `Dockerfile` for your bot
5. Copy the default `docker-compose.yml` file. You can use your own Compose file, or add to the default one. Just make sure that the `$TAG` variable is used to grab the right image, and that any environment variables you need to pass through to your various services are included under the `environment:` key.
6. Copy the example configuration, and modify the following environment variables (all paths relative to the configuration file, optional variables should work without modification):
    - `TEST_CMD`: the command to run the tests for the bot (if any, you can always just compile it)
    - `DEPLOY_ARTIFACTS_PATH`: the path of the `artifacts/` directory contained in this respository
    - `REPOSITORY` (**optional**): the repository name
    - `BRANCH` (**optional**): the git branch
    - `BETA_BRANCH`: the git branch for the beta bot
    - `ENCRYPTED_DEPLOY_KEY_PATH`: the path to the encrypted deploy key
    - `ENCRYPTED_DEPLOY_KEY_CYPHER_KEY`: the encrypted deploy key cypher key
    - `ENCRYPTED_DEPLOY_KEY_IV`: the encrpyted deploy key initialization vector
    - `DOCKER_IMAGE_SHELL`: the shell to be used to run commands inside the Docker container, this depends on which distribution you are using for your base image. Usually `sh` will work
7. If you are not using a [multistage build](https://docs.docker.com/develop/develop-images/multistage-build/), you should remove the `--target build` from the `docker build` command. You can also just mark your only step as `build` if you don't have multiple steps.

## How it Works
1. Sets up the SSH agent inside the CI server from the provided encypted key.
2. Builds the Docker image from the provided `Dockerfile`, only to the `build` step of a multistage build (so that the buiild tools needed for testing are available).
3. Runs the provided test command inside the built image
4. Fully builds the Docker image.
5. Tags the image using the provided Docker username, repository name, and either `beta` for the beta channel bot or `latest` for the stable channel bot.
6. Pushes the image to Docker Hub.
7. Populates a `.env` file with the appropriate Discord token depending on the bot channel (beta or stable), the appropriate Docker image tag as noted above, the Docker username as provided, and the repository as provided. These are needed to correctly populate the Discord token for the bot, and the information needed for the correct bot channel in the Docker Compose file.
8. Creates a directory for the appropriate bot channel using the provided deploy root directory.
9. Uses SCP to send the `.env` file and the `docker-compose.yml` file to the remote server using the provded SSH credentials.
10. Pulls the Docker image on the remote server.
11. Brings the bot up with Docker Compose in detatched mode. This will automatically restart the bot if the image has changed.

## Example Travis CI Configuration
```yaml
language: shell

services:
  - docker

env:
  global:
    - TEST_CMD='echo simulated test!'
    - SSH_KEY_TYPES='rsa,dsa,ecdsa'
    - DEPLOY_ARTIFACTS_PATH='discord-bot-deploy/artifacts'
    - DEPLOY_ROOT_DIR='/srv'
    - REPOSITORY=$(basename "$TRAVIS_REPO_SLUG")
    - BRANCH="${TRAVIS_BRANCH}"
    - BETA_BRANCH='develop'
    - DOCKER_IMAGE_TAG=$(bash "${DEPLOY_ARTIFACTS_PATH}"/generate_docker_image_tag.sh)
    - DOCKER_IMAGE_NAME="${DOCKER_USERNAME}"/"${REPOSITORY}":"${DOCKER_IMAGE_TAG}"
    - ENCRYPTED_DEPLOY_KEY_PATH="deploy_key.enc"
    - ENCRYPTED_DEPLOY_KEY_CYPHER_KEY="${encrypted_dfdcfd5172af_key}"
    - ENCRYPTED_DEPLOY_KEY_IV="${encrypted_dfdcfd5172af_iv}"
    - DOCKER_IMAGE_SHELL="sh"

before_deploy:
  - eval "$(ssh-agent -s)"
  - bash "${DEPLOY_ARTIFACTS_PATH}"/setup_ssh.sh

script:
  - docker build --target build --tag "$REPOSITORY" .
  - docker run "$REPOSITORY" ${DOCKER_IMAGE_SHELL} -c "$TEST_CMD"

deploy:
  - provider: script
    script: bash "${DEPLOY_ARTIFACTS_PATH}"/docker_push.sh
    on:
      branch: ${BETA_BRANCH}
  
  - provider: script
    script: bash "${DEPLOY_ARTIFACTS_PATH}"/docker_push.sh
    on:
      branch: master
  
  - provider: script
    script: bash "${DEPLOY_ARTIFACTS_PATH}"/deploy.sh
    on:
      branch: ${BETA_BRANCH}

  - provider: script
    script: bash "${DEPLOY_ARTIFACTS_PATH}"/deploy.sh
    on:
      branch: master

```

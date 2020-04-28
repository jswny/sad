# Deploy [![Build Status](https://travis-ci.com/jswny/deploy.svg?branch=master)](https://travis-ci.com/jswny/deploy)
A drop-in set of scripts to deploy apps using Docker and SSH.

## Features
- Automatic deployment from CI
- Containerized deployments
- Supports separate deployments for beta/stable channels
- Uses only SSH and Docker

## Requirements
- **A server with SSH, Docker, and Docker Compose** for the app to be deployed to
  - An SSH key to be used for deployment which has been added to an appropriate user on the server. The key can **not** have a password.
- **A [Docker](https://hub.docker.com/) account** to push/pull the images from
- **A CI service** to execute the deployments
  - **Note**: currently only Travis CI is supported, but porting this to other CI services should be easy, it is just a matter of modifying the configuration file.

## Usage
1. Include this repository as a git submodule inside your repository
2. Encrypt your private deployment SSH key into using OpenSSL `aes-256-cbc`, so that you have the encrypted key, the cypher key and the IV stored. For Travis CI this is `travis encrypt-file [--pro] ./deploy_key`, which will automatically add the two environment variables to your CI repository settings, and the local encrypted key file to the local filesystem.
3. Create the following secret environment variables in your CI service so that they won't be leaked:
    - `DEPLOY_SERVER`: the SSH-enabled server to deploy the app to
    - `DEPLOY_USERNAME`: the user to SSH into the server with (no password needed as)
    - `DOCKER_PASSWORD`: password for the Docker account
    - `DOCKER_USERNAME`: username for the Docker account
    - `ENCRYPTED_DEPLOY_KEY_CYPHER_KEY`: The encrypted deploy key cypher key
    - `ENCRYPTED_DEPLOY_KEY_CYPHER_IV`: The encrypted deploy key cypher initialization vector
    - One pair of environment variables for each variable which your app requires. Each one should have the same prefix, and either `_STABLE` or `_BETA` after the prefix to indicate which channel the variable corresponds to. For example, you should set `DEBUG_STABLE=false` and `DEBUG_BETA=true` if you want the variable `$DEBUG` to be available for your app.
4. Create a `Dockerfile` for your app
5. Copy the default `docker-compose.yml` file. You can use your own Compose file, or add to the default one. Just make sure that the `$TAG` variable is used in the actual image definition to grab the right image, and that any environment variables you need to pass through to your various app services are included under the `environment` key. 
  - **Note:** to use your Docker Compose file locally, simply create a local `.env` file and fill it in with the environment variables that your app needs (the same ones you listed under the `environment` key in your Compose File). You can also directly set those environment variables from the command line before using `docker-compose`.
  - **Note:** you can fill in the `$DOCKER_USERNAME` and `$REPOSITORY` environment variables with static values if you want to decrease the amount of variables required to run the app via Docker Compose.
6. Copy the example configuration, and modify the following environment variables (all paths relative to the configuration file, optional variables should work without modification):
  - `TEST_CMD`: the command to run the tests for the app (if any, you can always just compile it)
  - `DEPLOY_ARTIFACTS_PATH`: the path of the `artifacts/` directory contained in this respository
  - `REPOSITORY` (**optional**): the repository name
  - `BRANCH` (**optional**): the git branch
  - `BETA_BRANCH`: the git branch for the beta channel (`master` will be treated as the branch for the stable channel, and all other branches will not work)
  - `DEPLOY_CHANNEL_VAR_PREFIXES`: the environment variable prefixes for the variables required by your app (see step 3 above for more information)
  - `ENCRYPTED_DEPLOY_KEY_PATH`: the path to the encrypted deploy key
  - `ENCRYPTED_DEPLOY_KEY_CYPHER_KEY`: the encrypted deploy key cypher key
  - `ENCRYPTED_DEPLOY_KEY_IV`: the encrpyted deploy key initialization vector

## How it Works
1. Sets up the SSH agent inside the CI server from the provided encypted key.
2. Builds the Docker image from the provided `Dockerfile`, only to the `build` step of a multistage build (so that the buiild tools needed for testing are available).
3. Runs the provided test command inside the built image
4. Fully builds the Docker image.
5. Tags the image using the provided Docker username, repository name, and either `beta` for the beta channel or `latest` for the stable channel.
6. Pushes the image to Docker Hub.
7. Populates a `.env` file with the appropriate environment variables required by your app depending on the deploy channel (beta or stable), the appropriate Docker image tag as noted above, the Docker username as provided, and the repository as provided.
8. Creates a directory on the remote server for the app given the current deploy channel using the provided deploy root directory.
9. Uses SCP to send the `.env` file and the `docker-compose.yml` file to the remote server using the provded SSH credentials.
10. Pulls the Docker image on the remote server.
11. Brings the app up with Docker Compose in detatched mode. This will automatically restart the app if the image has changed.

## Example Travis CI Configuration
```yaml
language: shell

services:
  - docker

env:
  global:
    - SSH_KEY_TYPES='rsa,dsa,ecdsa'
    - DEPLOY_ARTIFACTS_PATH='deploy/artifacts'
    - DEPLOY_ROOT_DIR='/srv'
    - REPOSITORY=$(basename "$TRAVIS_REPO_SLUG")
    - BRANCH="${TRAVIS_BRANCH}"
    - BETA_BRANCH='develop'
    - DEPLOY_CHANNEL_VAR_PREFIXES="TOKEN,DEBUG"
    - ENCRYPTED_DEPLOY_KEY_PATH="deploy_key.enc"
    - ENCRYPTED_DEPLOY_KEY_CYPHER_KEY="${encrypted_dfdcfd5172af_key}"
    - ENCRYPTED_DEPLOY_KEY_IV="${encrypted_dfdcfd5172af_iv}"
    - DOCKER_IMAGE_SHELL="sh"

script:
  - 'echo simulated test!'

before_deploy:
  - eval "$(ssh-agent -s)"
  - bash "${DEPLOY_ARTIFACTS_PATH}"/setup_ssh.sh

deploy:
  - provider: script
    script: bash "${DEPLOY_ARTIFACTS_PATH}"/docker_push.sh
    on:
      branch: "${BETA_BRANCH}"
  
  - provider: script
    script: bash "${DEPLOY_ARTIFACTS_PATH}"/docker_push.sh
    on:
      branch: master
  
  - provider: script
    script: bash "${DEPLOY_ARTIFACTS_PATH}"/deploy.sh
    on:
      branch: "${BETA_BRANCH}"

  - provider: script
    script: bash "${DEPLOY_ARTIFACTS_PATH}"/deploy.sh
    on:
      branch: master

```

# Discord Bot Deploy [![Build Status](https://travis-ci.com/jswny/discord-bot-deploy.svg?branch=master)](https://travis-ci.com/jswny/discord-bot-deploy)
A pluggable set of scripts which can be dropped into a Discord bot repository to automatically deploy it.

## Features
- Automatic deployment
- Containerized deployment
- Supports separate deployments for beta/stable bots
- Uses only SSH and Docker for deployment

## Requirements
- **A server with SSH, Docker, and Docker Compose** for the bot to be deployed to
- **A [Docker Hub](https://hub.docker.com/) account** to push/pull the images from
- **A CI service** to execute the deployments

## Usage
**Note**: currently only Travis CI is supported, but porting this to other CI services should be easy, it is just a matter of modifying the configuration file.

1. Clone this repository into your Discord bot repository (or include it as a Git submodule)
2. Create the following secret environment variables in your CI service so that they won't be leaked:
   1. 
3. 

## Example Travis CI Configuration
```yaml
language: shell

services:
  - docker

env:
  global:
    - TEST_CMD='echo simulated test!'
    - SSH_KEY_TYPES='rsa,dsa,ecdsa'
    - DEPLOY_ARTIFACTS_DIR='artifacts'
    - DEPLOY_ROOT_DIR='/srv'
    - REPOSITORY="$(basename "$TRAVIS_REPO_SLUG")"
    - BRANCH="${TRAVIS_BRANCH}"
    - BETA_BRANCH='develop'
    - DOCKER_IMAGE_TAG=$(bash "${DEPLOY_ARTIFACTS_DIR}"/generate_docker_image_tag.sh)
    - DOCKER_IMAGE_NAME="${DOCKER_USERNAME}"/"${REPOSITORY}":"${DOCKER_IMAGE_TAG}"
    - ENCRYPTED_DEPLOY_KEY_PATH="deploy_key.enc"
    - DEPLOY_KEY="${encrypted_dfdcfd5172af_key}"
    - DEPLOY_KEY_IV="${encrypted_dfdcfd5172af_iv}"
    - DOCKER_IMAGE_SHELL="sh"

before_deploy:
  - eval "$(ssh-agent -s)"
  - bash "${DEPLOY_ARTIFACTS_DIR}"/setup_ssh.sh

script:
  - docker build --target build --tag "$REPOSITORY" .
  - docker run "$REPOSITORY" ${DOCKER_IMAGE_SHELL} -c "$TEST_CMD"

deploy:
  - provider: script
    script: bash "${DEPLOY_ARTIFACTS_DIR}"/docker_push.sh
    skip_cleanup: true
    on:
      branch: ${BETA_BRANCH}
  
  - provider: script
    script: bash "${DEPLOY_ARTIFACTS_DIR}"/docker_push.sh
    skip_cleanup: true
    on:
      branch: master
  
  - provider: script
    script: bash "${DEPLOY_ARTIFACTS_DIR}"/deploy.sh
    skip_cleanup: true
    on:
      branch: ${BETA_BRANCH}

  - provider: script
    script: bash "${DEPLOY_ARTIFACTS_DIR}"/deploy.sh
    skip_cleanup: true
    on:
      branch: master

```

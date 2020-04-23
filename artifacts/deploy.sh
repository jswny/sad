#!/usr/bin/env bash

verify_var_set 'BRANCH'
verify_var_set 'BETA_BRANCH'

if [ "$BRANCH" = "$BETA_BRANCH" ]
then
  DEPLOY_CHANNEL='beta'
elif [ "$BRANCH" = 'master' ]
then
  DEPLOY_CHANNEL='stable'
else
  echo "[ERROR] Unsupported branch $BRANCH" 1>&2
  exit 1
fi
verify_var_set 'DEPLOY_CHANNEL'

DISCORD_TOKEN_VARIABLE=$(echo "DISCORD_TOKEN_${DEPLOY_CHANNEL}" | tr '[:lower:]' '[:upper:]')
DISCORD_TOKEN="${!DISCORD_TOKEN_VARIABLE}"

verify_var_set 'DEPLOY_ROOT_DIR'
verify_var_set 'REPOSITORY'

DEPLOY_DIR="${DEPLOY_ROOT_DIR}"/"${REPOSITORY}"-"${DEPLOY_CHANNEL}"

verify_var_set 'DEPLOY_USERNAME'
verify_var_set 'DEPLOY_SERVER'
verify_var_set 'DEPLOY_DIR'
ssh "${DEPLOY_USERNAME}"@"$DEPLOY_SERVER" mkdir "$DEPLOY_DIR"

scp docker-compose.yml "${DEPLOY_USERNAME}"@"$DEPLOY_SERVER":"$DEPLOY_DIR"

verify_var_set 'DISCORD_TOKEN'
verify_var_set 'DOCKER_IMAGE_TAG'
{
  echo "DISCORD_TOKEN=${DISCORD_TOKEN}"
  echo "TAG=${DOCKER_IMAGE_TAG}"
  echo "DOCKER_USERNAME=${DOCKER_USERNAME}"
  echo "REPOSITORY=${REPOSITORY}"
} >> ".env"
scp ".env" "${DEPLOY_USERNAME}"@"$DEPLOY_SERVER":"$DEPLOY_DIR"

ssh "${DEPLOY_USERNAME}"@"$DEPLOY_SERVER" "docker pull '${DOCKER_IMAGE_NAME}'"
ssh "${DEPLOY_USERNAME}"@"$DEPLOY_SERVER" "cd '${DEPLOY_DIR}' && docker-compose up -d"

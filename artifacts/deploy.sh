#!/usr/bin/env bash

# shellcheck source=artifacts/utils.sh
source "$(dirname "$0")"/utils.sh

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

echo "Deploying to channel $DEPLOY_CHANNEL..."

verify_var_set 'DEPLOY_ROOT_DIR'
verify_var_set 'REPOSITORY'
verify_var_set 'DEPLOY_CHANNEL'

DEPLOY_DIR="${DEPLOY_ROOT_DIR}"/"${REPOSITORY}"-"${DEPLOY_CHANNEL}"

verify_var_set 'DEPLOY_USERNAME'
verify_var_set 'DEPLOY_SERVER'
verify_var_set 'DEPLOY_DIR'
ssh "${DEPLOY_USERNAME}"@"$DEPLOY_SERVER" mkdir "$DEPLOY_DIR"

scp docker-compose.yml "${DEPLOY_USERNAME}"@"$DEPLOY_SERVER":"$DEPLOY_DIR"

verify_var_set 'DOCKER_IMAGE_TAG'
{
  echo "TAG=${DOCKER_IMAGE_TAG}"
  echo "DOCKER_USERNAME=${DOCKER_USERNAME}"
  echo "REPOSITORY=${REPOSITORY}"
} >> ".env"


verify_var_set 'DEPLOY_CHANNEL_VAR_PREFIXES'
IFS=', ' read -r -a DEPLOY_CHANNEL_VAR_PREFIXES_ARRAY <<< "$DEPLOY_CHANNEL_VAR_PREFIXES"

for VAR_PREFIX in "${DEPLOY_CHANNEL_VAR_PREFIXES_ARRAY[@]}"
do
  VAR_NAME=$(echo "${VAR_PREFIX}_${DEPLOY_CHANNEL}" | tr '[:lower:]' '[:upper:]')
  VAR_VALUE="${!VAR_NAME}"
  verify_var_set "${!VAR_NAME}"
  echo "Setting deploy variable $VAR_NAME..."
  echo "${VAR_PREFIX}=${VAR_VALUE}" >> ".env"
done

scp ".env" "${DEPLOY_USERNAME}"@"$DEPLOY_SERVER":"$DEPLOY_DIR"

ssh "${DEPLOY_USERNAME}"@"$DEPLOY_SERVER" "docker pull '${DOCKER_IMAGE_NAME}'"
ssh "${DEPLOY_USERNAME}"@"$DEPLOY_SERVER" "cd '${DEPLOY_DIR}' && docker-compose up -d"

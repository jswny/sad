#!/usr/bin/env bash

DEPLOY_KEY_VARIABLE="${DEPLOY_KEY_ENV_VAR_PREFIX}_key"
DEPLOY_IV_VARIABLE="${DEPLOY_KEY_ENV_VAR_PREFIX}_iv"

openssl aes-256-cbc -K "${!DEPLOY_KEY_VARIABLE}" -iv "${!DEPLOY_IV_VARIABLE}" -in "${DEPLOY_ARTIFACTS_DIR}"/deploy_key.enc -out "${DEPLOY_ARTIFACTS_DIR}"/deploy_key -d

chmod 600 "${DEPLOY_ARTIFACTS_DIR}"/deploy_key

ssh-add "${DEPLOY_ARTIFACTS_DIR}"/deploy_key

{ ssh-keyscan -t "$SSH_KEY_TYPES" -H "$DEPLOY_SERVER" >> "${TRAVIS_HOME}"/.ssh/known_hosts; } 2>&1

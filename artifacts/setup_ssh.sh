#!/usr/bin/env bash

openssl aes-256-cbc -K "${ENCRYPTED_DEPLOY_KEY_CYPHER_KEY}" -iv "${ENCRYPTED_DEPLOY_KEY_IV}" -in "${ENCRYPTED_DEPLOY_KEY_PATH}" -out "${DEPLOY_ARTIFACTS_PATH}"/deploy_key -d

chmod 600 "${DEPLOY_ARTIFACTS_PATH}"/deploy_key

ssh-add "${DEPLOY_ARTIFACTS_PATH}"/deploy_key

{ ssh-keyscan -t "$SSH_KEY_TYPES" -H "$DEPLOY_SERVER" >> "${TRAVIS_HOME}"/.ssh/known_hosts; } 2>&1

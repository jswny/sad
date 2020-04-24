#!/usr/bin/env bash

# shellcheck source=artifacts/utils.sh
source "$(dirname "$0")"/utils.sh

verify_var_set 'ENCRYPTED_DEPLOY_KEY_CYPHER_KEY'
verify_var_set 'ENCRYPTED_DEPLOY_KEY_IV'
verify_var_set 'ENCRYPTED_DEPLOY_KEY_PATH'
verify_var_set 'DEPLOY_ARTIFACTS_PATH'

openssl aes-256-cbc -K "${ENCRYPTED_DEPLOY_KEY_CYPHER_KEY}" -iv "${ENCRYPTED_DEPLOY_KEY_IV}" -in "${ENCRYPTED_DEPLOY_KEY_PATH}" -out "${DEPLOY_ARTIFACTS_PATH}"/deploy_key -d

chmod 600 "${DEPLOY_ARTIFACTS_PATH}"/deploy_key

ssh-add "${DEPLOY_ARTIFACTS_PATH}"/deploy_key

verify_var_set 'SSH_KEY_TYPES'
verify_var_set 'DEPLOY_SERVER'

{ ssh-keyscan -t "$SSH_KEY_TYPES" -H "$DEPLOY_SERVER" >> "${TRAVIS_HOME}"/.ssh/known_hosts; } 2>&1

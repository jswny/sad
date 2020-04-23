#!/usr/bin/env bash

if [ -z "$ENCRYPTED_DEPLOY_KEY_CYPHER_KEY" ]
then
  echo '$ENCRYPTED_DEPLOY_KEY_CYPHER_KEY is blank or unset! Exiting...'
  exit 1
fi

if [ -z "$ENCRYPTED_DEPLOY_KEY_IV" ]
then
  echo '$ENCRYPTED_DEPLOY_KEY_IV is blank or unset! Exiting...'
  exit 1
fi

if [ -z "$ENCRYPTED_DEPLOY_KEY_PATH" ]
then
  echo '$ENCRYPTED_DEPLOY_KEY_PATH is blank or unset! Exiting...'
  exit 1
fi

if [ -z "$DEPLOY_ARTIFACTS_PATH" ]
then
  echo '$ENCRYPTED_DEPLOY_KEY_PATH is blank or unset! Exiting...'
  exit 1
fi

openssl aes-256-cbc -K "${ENCRYPTED_DEPLOY_KEY_CYPHER_KEY}" -iv "${ENCRYPTED_DEPLOY_KEY_IV}" -in "${ENCRYPTED_DEPLOY_KEY_PATH}" -out "${DEPLOY_ARTIFACTS_PATH}"/deploy_key -d

chmod 600 "${DEPLOY_ARTIFACTS_PATH}"/deploy_key

ssh-add "${DEPLOY_ARTIFACTS_PATH}"/deploy_key

if [ -z "$SSH_KEY_TYPES" ]
then
  echo '$SSH_KEY_TYPES is blank or unset! Exiting...'
  exit 1
fi

if [ -z "$DEPLOY_SERVER" ]
then
  echo '$DEPLOY_SERVER is blank or unset! Exiting...'
  exit 1
fi

{ ssh-keyscan -t "$SSH_KEY_TYPES" -H "$DEPLOY_SERVER" >> "${TRAVIS_HOME}"/.ssh/known_hosts; } 2>&1

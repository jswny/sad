#!/usr/bin/env bash

# Exit on any error, undefined variable, or pipe failure 
set -eo pipefail

log() {
  local prefix_spacer="-----"
  local component="Deploy"
  local prefix="${prefix_spacer} [${component}]"
  if [ "${1}" = 'debug' ]; then
    if [ "${debug}" = 'true' ]; then
      echo "::debug::${prefix} DEBUG: ${2}"
    fi
  elif [ "${1}" = 'info' ]; then
    echo "${prefix} INFO: ${2}"
  elif [ "${1}" = 'warn' ]; then
    echo "::warning::${prefix} WARN: ${2}" >&2
  elif [ "${1}" = 'error' ]; then
    echo "::error::${prefix} ERROR: ${2}" >&2
  else
    echo "::error::${prefix} INTERNAL ERROR: invalid option \"${1}\" for log() with message \"${2}\"" >&2 "" >&2
  fi
}

verify_var_set() {
  if [ -z "${!1}" ]; then
    if [ -z "${2}" ]; then
      log 'error' "\"${1}\" is blank or unset!"
    else
      log 'error' "${2}"
    fi
    exit 1
  fi
}

check_exists_file() {
  if [ ! -f "${1}" ]; then
    if [ -e "${1}" ]; then
      log 'error' "Item at path \"${1}\" exists, but it is not a file!"
    else
      log 'error' "File \"${1}\" does not exist!"
    fi
    exit 1
  else
    log 'debug' "File \"${1}\" exists!"
  fi
}

scp_wrapper() {
  log 'info' "Sending \"$(basename "${1}")\" to deploy server..."

  if [ "${debug}" = 'true' ]; then
    scp -v "${1}" "${deploy_username}@${deploy_server}:${deploy_dir}"
  else
    scp "${1}" "${deploy_username}@${deploy_server}:${deploy_dir}"
  fi
}

ssh_wrapper() {
  log 'debug' "Running command over SSH \"${1}\""

  if [ "${debug}" = 'true' ]; then
    ssh -v "${deploy_username}@${deploy_server}" "${1}"
  else
    ssh "${deploy_username}@${deploy_server}" "${1}"
  fi
}

any_branch_identifier='ANY'

# Translate input environment variables
deploy_server="${INPUT_DEPLOY_SERVER}"
deploy_username="${INPUT_DEPLOY_USERNAME}"
deploy_root_dir="${INPUT_DEPLOY_ROOT_DIR}"
ssh_key="${INPUT_SSH_KEY}"
app_path="${INPUT_PATH}"
stable_branch="${INPUT_STABLE_BRANCH}"
beta_branch="${INPUT_BETA_BRANCH}"
debug="${INPUT_DEBUG}"
env_var_prefixes="${INPUT_ENV_VAR_PREFIXES}"

repository="${GITHUB_REPOSITORY}"
ref="${GITHUB_REF}"

verify_var_set 'ref' 'GITHUB_REF is blank or unset!'
verify_var_set 'repository' 'GITHUB_REPOSITORY is blank or unset!'

verify_var_set 'stable_branch'
verify_var_set 'beta_branch'

log 'info' 'Detecting Git and release info...'

if echo "${ref}" | grep -qE '^refs\/(tags|remote)\/'; then
  ref_type='tag/remote'
elif echo "${ref}" | grep -qE '^refs\/heads\/'; then
  ref_type='branch'
fi

verify_var_set 'ref_type' "Could not detect valid ref type from ref \"${ref}\""

ref_name=$(echo "${ref}" | sed -E 's/refs\/(heads|tags|remote)\///')

verify_var_set 'ref_name' 'Could not extract a proper supported Git reference name!'

log 'debug' "Reference type detected \"${ref_type}\" with name \"${ref_name}\""

if [ "${ref_type}" = 'tag/remote' ]; then
  log 'error' "Unsupported reference \"${ref}\" with detected reference type \"${ref_type}\""
  exit 1
elif [ "${ref_type}" = 'branch' ]; then
  if [ "${ref_name}" = "${stable_branch}" ] || [ "${stable_branch}" = "${any_branch_identifier}" ]; then
    channel='stable'
  elif [ "${ref_name}" = "${beta_branch}" ] || [ "${beta_branch}" = "${any_branch_identifier}" ]; then 
    channel='beta'
  fi
fi

verify_var_set 'channel' "Could not detect release channel from reference type \"${ref_type}\" and reference name \"${ref_name}\""

log 'info' "Detected release channel \"${channel}\""

log 'info' 'Verifying action inputs...'

home_path="/root"
repository_path="${GITHUB_WORKSPACE}"

verify_var_set 'repository_path' 'GITHUB_WORKSPACE is blank or unset!'
verify_var_set 'deploy_server'
verify_var_set 'deploy_username'
verify_var_set 'deploy_root_dir'
verify_var_set 'ssh_key'
verify_var_set 'app_path' 'path is blank or unset!'

debug="${debug:-false}"

full_app_path="${repository_path}/${app_path}"
verify_var_set 'full_app_path' 'Could not generate full app path based on provided app path!'

log 'info' 'Detecting local Docker image...'

local_image_id="$(docker images -q "${GITHUB_REPOSITORY}" 2> /dev/null)"

verify_var_set 'local_image_id' "No local Docker image detected for this repository! Please build a local image first before deploying, and ensure it is tagged with the name of this repository \"${repository}\"."

log 'debug' "Local Docker image ID(s) detected: \"${local_image_id}\""

if [ "$(echo "${local_image_id}" | grep -c '$')" -gt 1 ]; then
  log 'error' $"Detected multiple Docker image IDs for this repository! Make sure there is only one Docker image tagged with the name of this repository \"${repository}\"."
  exit 1
fi

local_image="$(docker inspect --format='{{ (index .RepoTags 0) }}' "${local_image_id}" 2> /dev/null)"

verify_var_set 'local_image' 'Could not detect the local Docker image name and tag!'

log 'debug' "Detected local Docker image name and tag: ${local_image}"

local_image_name="$(echo "${local_image}" | sed -E 's/^.*\///' | sed -E 's/:.*$//')"

verify_var_set 'local_image_name' "Could not parse local image name (without tag or username) from full image name ${local_image_name}!"

log 'debug' "Parsed local image name \"${local_image_name}\""

container_name="${local_image_name}-${channel}"
verify_var_set 'container_name' 'Could not generate container name for deployment!'

log 'info' "Generated container name for deployment ${container_name}"

log 'info' 'Adding SSH key to agent...'

eval "$(ssh-agent -s)"

ssh-add - <<< "${ssh_key}"

ssh_path="${home_path}/.ssh"
verify_var_set 'ssh_path'
mkdir -p "${ssh_path}"

log 'info' 'Adding SSH key(s) to known hosts...'

ssh-keyscan "${deploy_server}" >> "${ssh_path}/known_hosts"

check_exists_file "${ssh_path}/known_hosts"

deploy_dir="${deploy_root_dir}/${container_name}"
verify_var_set 'deploy_dir' 'Could not generate deploy directory path!'

log 'info' "Creating deploy directory \"${deploy_dir}\" on deploy server"

ssh_wrapper "mkdir -p '${deploy_dir}'"

log 'info' 'Generating ".env" file for deployment...'

env_file_name='deploy.env'
env_file_path="${full_app_path}/${env_file_name}"

{
  echo "IMAGE=${local_image}"
  echo "CONTAINER_NAME=${container_name}"
} >> "${env_file_path}"

if [ -z "${env_var_prefixes}" ]; then
  log 'info' 'No custom environment variables found to inject into the deployment. See the "env_var_prefixes" input to add some.'
else
  IFS=', ' read -r -a env_var_prefixes_array <<< "${env_var_prefixes}"
  for env_var_prefix in "${env_var_prefixes_array[@]}"; do
    env_var_name=$(echo "${env_var_prefix}_${channel}" | tr '[:lower:]' '[:upper:]')
    verify_var_set "${env_var_name}" "Environment variable \"${env_var_name}\" generated from environment variable prefixes is blank or unset!"
    env_var_value="${!env_var_name}"
    log 'debug' "Setting deploy environment variable ${env_var_name}"
    echo "${env_var_prefix}=${env_var_value}" >> "${env_file_path}"
  done
fi

scp_wrapper "${env_file_path}"
ssh_wrapper "mv '${deploy_dir}/${env_file_name}' '${deploy_dir}/.env'"

scp_wrapper "${full_app_path}/docker-compose.yml"

log 'info' "Pulling pushed image \"${local_image}\" from deploy server..."

ssh_wrapper "docker pull '${local_image}'"

log 'info' "Bringing app up on deploy server with Docker Compose..."

ssh_wrapper "cd '${deploy_dir}' && docker-compose up -d"

log 'info' 'Done!'

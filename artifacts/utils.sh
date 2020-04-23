#!/usr/bin/env bash

verify_var_set() {
  if [ -z ${!1} ]
  then
    echo "[ERROR] $1 is blank or unset! Exiting..." 1>&2
    exit 1
  fi
} 

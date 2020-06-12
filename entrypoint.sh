#!/bin/bash

echo "Variable: $1"
time=$(date)
echo "::set-output name=time::$time"

#!/bin/bash

echo "Variable: $1"
echo "FOOBAR: $FOOBAR"
time=$(date)
echo "::set-output name=time::$time"

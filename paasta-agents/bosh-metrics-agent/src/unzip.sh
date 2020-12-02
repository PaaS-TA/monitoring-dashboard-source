#!/usr/bin/env bash

for file in `ls *.zip`; do unzip "${file}"; done

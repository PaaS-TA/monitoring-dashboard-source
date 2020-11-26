#!/usr/bin/env bash

for file in `ls *.zip`; do unzip "${file}" -d "${file:0:-4}"; done

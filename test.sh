#!/bin/bash

TASK="1068"

cd project
echo "[*] rm -rf ~/.cses/"
#rm -rf ~/.cses/

echo "[*] go run ./*.go login"
go run ./*.go login

echo "[*] go run ./*.go github"
go run ./*.go github

echo "[*] go run ./*.go list"
go run ./*.go list

echo "[*] go run ./*.go show ${TASK}"
go run ./*.go show "${TASK}"

echo "[*] go run ./*.go solve ${TASK}"
go run ./*.go solve "${TASK}"

echo "[*] go run ./*.go submit ${TASK}.task.cpp"
go run ./*.go submit "${TASK}".task.cpp

cd ..

#!/bin/bash

echo "[*] rm -rf ~/.cses/"
rm -rf ~/.cses/

echo "[*] go run ./*.go login"
go run ./*.go login

echo "[*] go run ./*.go list"
go run ./*.go list

echo "[*] go run ./*.go show 1742"
go run ./*.go show 1742

echo "[*] go run ./*.go solve 1742"
go run ./*.go solve 1742

echo "[*] go run ./*.go submit 1742.task.cpp"
go run ./*.go submit 1742.task.cpp

#!/bin/zsh

go run ./server/server.go & 
echo started http://localhost:7353
ls ../**/*.(go|html|js)  | entr make remake
killall `lsof -i :7353 | grep -w "(LISTEN)" | awk '{print $2}'`

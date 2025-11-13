#! /bin/bash

session='qr-pastebin'
session_exists=$(tmux list-sessions | grep $session)

if [ "$session_exists" = "" ]
then
    tmux new-session -d -s $session

    # Main window for GIT
    tmux rename-window -t 0 'main'

    # Web window
    tmux new-window -t $session:1 -n 'web'
    tmux send-keys -t 'web' 'cd web' C-m 'npm run dev' C-m

    # API window
    tmux new-window -t $session:2 -n 'api'
    tmux send-keys -t 'api' 'cd api' C-m 'docker compose -f docker-compose.yaml up -d --build --scale api=2' C-m
fi

tmux attach-session -t $session:0
#! /bin/bash

session='qr-pastebin'
session_exists=$(tmux list-sessions | grep $session)

if [ "$session_exists" = "" ]
then
    tmux new-session -d -s $session

    # Main window for GIT
    tmux rename-window -t 0 'main'

    # API window
    tmux new-window -t $session:1 -n 'api'
    tmux send-keys -t 'api' 'cd api' C-m 'DATABASE_URL=postgres://postgres:root@localhost:5432/qr_pastebin go run .' C-m

    # Web window
    tmux new-window -t $session:2 -n 'web'
    tmux send-keys -t 'web' 'cd web' C-m 'npm run developmentserver' C-m
fi

tmux attach-session -t $session:0
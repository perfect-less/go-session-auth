#! /bin/bash

# I've got this script from https://ryan.himmelwright.net/post/scripting-tmux-workspaces/
# credit to ryan for this script.

session="gosessionauth"

# For session checking
SESSIONEXISTS=$(tmux list-sessions | grep $session)

if [ "$SESSIONEXISTS" = "" ]
then
    # Create new session
    tmux new-session -d -s $session

    # Create nvim window:1
    tmux new-window -t $session:1 -n 'Code'
    tmux send-keys -t $session:1 "nvim ." C-m

    # Create two more aux windows
    tmux new-window -t $session:2 
    tmux new-window -t $session:3
fi

# Attaching to the session
tmux attach-session -t $session:1

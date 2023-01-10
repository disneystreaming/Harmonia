#!/bin/bash
read -p "Enter oauth client id: " client_id
echo "export OAUTH_CLIENT_ID=$client_id" > localenv
stty -echo
read -p "Enter oauth client secret: " client_secret
stty echo
echo
echo "export OAUTH_CLIENT_SECRET=$client_secret" >> localenv
stty -echo
read -p "Enter git machine token: " git_token
stty echo
echo
echo "export GIT_MACHINE_TOKEN=$git_token" >> localenv
read -p "Enter tracking repo: " tracking_repo
echo "export TRACKING_REPOSITORY=$tracking_repo" >> localenv
echo "export IS_LOCAL=true" >> localenv
echo "\n\nLocal run configuration written to localenv.\nUse 'source localenv' to apply them. Don't forget to 'rm localenv' when finished!"

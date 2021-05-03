set -eu

          
mkdir -p ~/.ssh/
echo "${ACTION_DEPLOY_KEY}" > ~/.ssh/id_rsa
chmod 700 ~/.ssh
chmod 600 ~/.ssh/id_rsa
ssh-keyscan github.com >> ~/.ssh/known_hosts
git config --global user.email "elitegoblinrb@gmail.com"
git config --global user.name "Frank Sun"

npm install
sh
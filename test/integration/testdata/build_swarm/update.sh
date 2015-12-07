#!/bin/bash

set -e

if [ $# -eq 0 ] ; then
	echo "Usage: ./update.sh <docker/swarm tag or branch>"
	exit
fi

VERSION=$1

#set environemnt GITHUBPATH as environemnt variable for user defined
if [ $GITHUBPATH ]; then
  GITHUB_PATH=$GITHUBPATH
else
  GITHUB_PATH="https://github.com/docker/swarm.git"
fi

# cd to the current directory so the script can be run from anywhere.
cd `dirname $0`

# Update the certificates.
echo "Updating certificates..."
./certs/update.sh

echo "Fetching and building swarm $VERSION..."

# Create a temporary directory.
TEMP=`mktemp -d`

git clone -b $VERSION $GITHUB_PATH $TEMP
docker build -t swarm-builder $TEMP

# Create a dummy swarmbuild container so we can run a cp against it.
ID=$(docker create swarm-builder)

# Update the local binary.
docker cp $ID:/go/bin/swarm .

# Cleanup.
docker rm -f $ID
docker rmi swarm-builder

echo "Done."

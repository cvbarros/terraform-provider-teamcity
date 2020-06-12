#!/bin/bash
set -e

pushd integration_tests/

# using tar which is readily available whereas bsdtar is not even available on some linux distros
tar -xzf teamcity_data.tar.gz

docker-compose up -d

until $(curl --output /dev/null --silent --head --fail http://localhost:8112/login.html); do
    echo "Waiting for TeamCity to become available.."
    sleep 5
done

echo "TeamCity is ready!"

popd

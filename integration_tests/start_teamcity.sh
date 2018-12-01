#!/bin/bash
set -e

pushd integration_tests/

tar -xzf teamcity_data.tar.gz

docker-compose up -d

until $(curl --output /dev/null --silent --head --fail http://localhost:8112/login.html); do
    echo "Waiting for TeamCity to become available.."
    sleep 5
done

echo "TeamCity is ready!"

popd

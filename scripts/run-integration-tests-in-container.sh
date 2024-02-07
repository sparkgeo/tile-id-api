#!/bin/bash

dco="docker compose -p integrationtest"

$dco build integrationtest
$dco up api -d
$dco run integrationtest
exit_code=$?
$dco down
exit $exit_code

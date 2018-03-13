#!/bin/bash

basedir="$( cd -P "$( dirname "$0" )" && pwd )"

ssh-keygen -f "${basedir}/host_key"
ssh-keygen -f "${basedir}/user_key"

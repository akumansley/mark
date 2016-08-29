#!/bin/bash
set -euo pipefail
# This script is run in the VM each time you run `vagrant-spk dev`.  This is
# the ideal place to invoke anything which is normally part of your app's build
# process - transforming the code in your repository into the collection of files
# which can actually run the service in production
#
# Some examples:
#
#   * For a C/C++ application, calling
#       ./configure && make && make install
#   * For a Python application, creating a virtualenv and installing
#     app-specific package dependencies:
#       virtualenv /opt/app/env
#       /opt/app/env/bin/pip install -r /opt/app/requirements.txt
#   * Building static assets from .less or .sass, or bundle and minify JS
#   * Collecting various build artifacts or assets into a deployment-ready
#     directory structure

# By default, this script does nothing.  You'll have to modify it as
# appropriate for your application.
cd /opt/app
export GOPATH=/home/vagrant

# TODO uncomment these
# npm install
npm run build
go get github.com/awans/mark/cmd/mark
rm -rf /home/vagrant/src/github.com/awans/mark
ln -s /opt/app /home/vagrant/src/github.com/awans/mark
go build -o ./mark cmd/mark/main.go


exit 0

#!/bin/bash
user="rocket-notify"
path="/usr/local/bin/rocket-notify"
config="/etc/rocket-notify/config.yml"

useradd $user
curl -L -s https://api.github.com/repos/ksarink/rocket-notify/releases/latest | grep "browser_download_url" | cut -d '"' -f 4 | wget -O $path -qi -
chown $user:root $path
chmod +x $path
chmod u+s $path

mkdir -p /etc/rocket-notify/

if [ ! -f "$config" ]; then
	echo -e 'baseurl: "https://"\nuserid: ""\ntoken: ""' > $config
fi
chown $user $config
chmod 600 $config

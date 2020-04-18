# rocket-notify - Send Rocket.Chat notifications from the commandline!

This tool allows you to send quick notifications or whole messages to yourself from to commandline. Simply add this tool at the end of your calculation, and you receive a message when it finishes.

Example usages (`sleep 10` is the placeholder for your calculation):
```
sleep 10; rocket-notify                         # simple
sleep 10; rocket-notify Calculation finished!   # sends the message "Calculation finished!"
sleep 10; rocket-notify -emoji :astonished:     # use a emoji as avator 
sleep 10; rocket-notify -sender Application1    # define the name of the sender 
sleep 10; rocket-notify -sender="Application 1" # parameter style with and without equal sign or quote marks - if not seperated by space
sleep 10 | rocket-notify finished               # sends the output of the preceding command with the message "finished"

# Help:
rocket-notify --help
rocket-notify -h 
```

Please install by executing the rocket-notify-install.sh script as a privileged user. Afterwards edit the config (/etc/rocket-notify/config.yml) according to your needs.
The installer downloads the latest release into the /usr/local/bin directory, creates a user, sets ownership of the binary and the user-sticky bit. Also the ownership of the config is set to the created user and its permissions to 600. This setup allows to safely write your REST-API-token into your config in a multi user environment. With the user sticky bit, the binary can access the config but your users cannot (even though they can execute the binary and indirectly use the config)

The Rocket.Chat username is derived from the currently logged in user. Which helps when you have connected your os and your Rocket.Chat instance to an LDAP server. In the future alternatively you can add a local config to define your user. 

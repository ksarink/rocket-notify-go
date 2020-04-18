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
# monitor-zCatch

It is the monitoring and command execution service that connects to a zCatch server and creates events.

This servic ehas two purposes:

- parsing of server logs and creation of events as well as pushing them to the event specific topic at the broker
- subscribing to the "econ:port" topic and waiting for command events that need to be executed
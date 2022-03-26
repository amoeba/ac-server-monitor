# ac-server-monitor

Status and historical uptime information for Asheron's Call [private servers](https://github.com/acresources/serverslist).

https://servers.treestats.net

## How this works

Every ten minutes, we fetch the [wild west server list](https://github.com/acresources/serverslist) and check whether server is up or down by sending a login packet.
The result is stored and the application shows current and historical status for each server.

## API

There's a pretty basic API.
Feel free to build stuff with it:

- [`/api`](https://servers.treestats.net/api): List of API routes
- [`/api/servers/`](https://servers.treestats.net/api/servers): List of all servers and their statuses
- [`/api/uptime/:id`](https://servers.treestats.net/uptime/1): Recent uptime information for a single server
- [`/api/logs`](https://servers.treestats.net/logs): Recent logs from the application itself. This is mostly for debugging

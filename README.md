# ac-server-monitor

Status and historical uptime information for Asheron's Call [private servers](https://github.com/acresources/serverslist).

View it live at <https://servers.treestats.net>.

## How this works

Every ten minutes, we fetch the [community server list](https://github.com/acresources/serverslist) and check whether server is up or down by sending a login packet.
The result is stored and the application shows current and historical status for each server.

## Development Setup

### Building

To build the application and the database seeding tool:

```bash
make
```

This will create two executables:
- `monitor` - The main web application
- `seed` - Database seeding tool for development

### Database Seeding

For development and testing, you can populate the database with fake data:

```bash
./seed
```

This will:
- Create the database automatically 
- Run all migrations defined in AutoMigrate
- Generate 10 servers with realistic two-part names
- Create 2 weeks of fake uptime data for each server

The seeded data includes servers with names like "Phoenix Fortress", "Crystal Tower", etc., and realistic uptime patterns with occasional downtime.

### Running the Application

After seeding the database, you can start the web application:

```bash
./monitor --no-cron
```

The `--no-cron` flag prevents the application from trying to fetch real server data.

## API

There's a pretty basic API.
Feel free to build stuff with it:

- [`/api`](https://servers.treestats.net/api): List of API routes
- [`/api/servers/`](https://servers.treestats.net/api/servers): List of all servers and their statuses
- [`/api/uptime/:id`](https://servers.treestats.net/uptime/1): Recent uptime information for a single server

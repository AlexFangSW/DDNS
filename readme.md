# DDNS service
Updates DNS 'A' records on CloudFlare.

## Build and run
```bash
./scripts/docker_build.sh && ./scripts/docker_run.sh
```

## Example output
```bash
2024/07/20 07:21:32 start updating dns record
2024/07/20 07:21:33 current public ipv4: "PUBLIC_IP"
2024/07/20 07:21:33 current record, type: "A", name: "DOMAIN", content: "PUBLIC_IP"
2024/07/20 07:21:33 same ip, no need to change
2024/07/20 07:21:33 finish updating dsn record
```

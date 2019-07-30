# Slack wifi changer

[![Build Status](https://cloud.drone.io/api/badges/JeffreyVdb/slack-working-from-home/status.svg)](https://cloud.drone.io/JeffreyVdb/slack-working-from-home)

## Getting started

1. Get a slack legacy token
2. build container

```shell
./build_container.sh slack-wifi-changer
```

3. Create status file like `status-example.json`
4. Run container with `SLACK_TOKEN` environment variable and configuration file

```bash
podman run --rm --net=host -it -v $HOME/slack-status.json:/slack.json:Z \
    -e SLACK_STATUS_FILE=/slack.json -e SLACK_TOKEN=your-slack-token localhost/slack-wifi-changer ```

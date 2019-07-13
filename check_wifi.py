#!/usr/bin/env python3
import sys
import subprocess
import logging
import json
import random
from typing import Optional
from dataclasses import dataclass
from os import getenv

import requests

DEFAULT_SLACK_STATUS_FILE = "/etc/slack-status.json"

# These all return {"ip": "<ipv4_address>"}
IPV4_PUBLIC_IP_SERVICES = (
    "https://ip4.seeip.org/json",
    "https://api.ipify.org/?format=json",
    "https://ip4.seeip.org/json",
)


@dataclass
class SlackStatus:
    text: str
    emoji: str


class SlackException(Exception):
    pass


def get_public_ip(api_urls):
    return requests.get(random.choice(api_urls)).json()["ip"]


def setup_logging():
    logging.basicConfig(stream=sys.stdout, level=logging.INFO)


def get_current_wifi_network():
    return subprocess.check_output(["iwgetid", "-r"]).decode("utf-8").strip()


def get_slack_status_from_file(filename, wifi_name, public_ip) -> Optional[SlackStatus]:
    with open(filename, "r", encoding="utf8") as fh:
        contents = json.load(fh)

    def match_wifi(entry):
        ssids = entry.get("wifi_names")
        public_ips = entry.get("public_ips", None)
        return wifi_name in ssids and (not public_ips or public_ip in public_ips)

    try:
        entry = next(filter(match_wifi, contents))

        return SlackStatus(text=entry["status_text"], emoji=entry["status_emoji"])
    except StopIteration:
        return None


def set_slack_status(status, *, api_token):
    payload = {"profile": {"status_text": status.text, "status_emoji": status.emoji}}
    headers = {
        "Content-Type": "application/json; charset=utf-8",
        "Authorization": f"Bearer {api_token}",
    }

    response_json = requests.post(
        "https://slack.com/api/users.profile.set", json=payload, headers=headers
    ).json()

    if not response_json["ok"]:
        raise SlackException(response_json["error"])


def main():
    setup_logging()
    status_file = getenv("SLACK_STATUS_FILE", DEFAULT_SLACK_STATUS_FILE)
    api_token = getenv("SLACK_TOKEN", "")
    wifi_network = get_current_wifi_network()
    public_ip = get_public_ip(IPV4_PUBLIC_IP_SERVICES)
    status = get_slack_status_from_file(
        status_file, wifi_name=wifi_network, public_ip=public_ip
    )
    if not status:
        logging.warning("No match found for SSID %s and IP %s", wifi_network, public_ip)
        sys.exit(0)

    try:
        set_slack_status(status, api_token=api_token)
    except SlackException as ex:
        logging.error(ex)
        sys.exit(1)


if __name__ == "__main__":
    main()

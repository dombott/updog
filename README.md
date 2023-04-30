# What's updog?
Well, hopefully your application is :point_left::sunglasses::point_left:

If not, updog will tell you.

Updog is statuspages companion and best friend. It receives alerts via webhook from prometheus-alertmanager and automatically creates/updates GitHub issues to be displayed as incidents using [dombott/statuspage](https://github.com/dombott/statuspage).

# Configuration
The following flags and env vars can be used to configure updog.
```
addr := flag.String("addr", ":8080", "address to listen on")
owner := flag.String("owner", "", "github repo owner")
repo := flag.String("repo", "", "github repo name")
token := os.Getenv("GH_TOKEN")
```

`addr`: the address of the webhook. Updog handles requests to `/webhook` and `/healthz`.

`owner` and `repo`: the owner and repo name that updog operates in.

`token`: a GitHub PAT. Needs permission on the repo to read and write issues and labels. 

# How does it work?
Updog generates an issue (title, body, labels) from an alert using its metadata.

It uses the content of the labels `updog/title`, `updog/body` and `updog/labels` to fill the respective field.
Multiple labels can be added as a semicolon separated list. The label `type/incident` is added automatically (to make the incident show up in statuspage).
It is recommended to use the alertmanager templating to fill these labels with values from the prometheus metrics.

To identify an issue later, the fingerprint of the alert is added to the title during creation (statuspage will cut off this identifier when displaying the incident).
The alertmanager fingerprint is the hash value of all the labels of an alert.

If a `firing` alert comes from the alertmanager, updog will check if there is already an open issue for the alert by searching for the identifier.
If the issue already exists, updog does nothing. If the issue doesn't exist, updog will create the issue (and the labels if necessary).

If a `resolved` alert comes from the alertmanager, updog will check if there is an open issue for the alert by searching for the identifier.
If the issue exists, updog will close it.

# How do I use this?
Simply run updog next to your alertmanager and set up a route to send alerts to updogs webhook.
```
route:
  group_by: ['alertname']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 3h

  routes:
  - match:
      severity: critical
    receiver: updog

receivers:
- name: updog
  webhook_configs:
  - url: 'http://updog.default.svc.cluster.local:8080/webhook'
```

Add the labels `updog/title`, `updog/body` and `updog/labels` to your alerts and fill them via alertmanager templating.

Configure updog to use the GitHub repo of your choice and provide a GitHub PAT for auth.

If you want to use statuspage together with updog, configure both to use the same GitHub repo.

# What's updog?
Not much, how about you?

Jokes aside, updog is statuspages companion and best friend.

Updog receives alerts via webhook from prometheus-alertmanager and automatically creates/updates GitHub issues to be displayed as incidents using [dombott/statuspage](https://github.com/dombott/statuspage).

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
To identify an issue later, the hash of the generated issue is added to the title during creation (statuspage will cut off this identifier when displaying the issue).

If a `firing` alert comes from the alertmanager, updog will check if there is already an issue for the alert by searching for the hash value.
If the issue already exists, updog does nothing. If the issue doesn't exist, updog will create the issue (and the labels if necessary).

If a `resolved` alert comes from the alertmanager, updog will check if there is an issue for the alert by searching for the hash value.
If the issue exists, updog will close it.

# How do I use this?
Simply run updog next to your alertmanager and set up a route to send alerts to updogs webhook.

Configure updog to use the GitHub repo of your choice and provide a GitHub PAT for auth.

If you want to use statuspage together with updog, configure both to use the same GitHub repo.

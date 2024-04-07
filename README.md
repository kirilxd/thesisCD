# Overview

ThesisCD is simplified and minimalistic GitOps engine used to demonstrate principles and limitations of pull-based and push-based GitOps approaches.

# Usage

For pull-based:

`thesisCD gitpoll --interval=<polling interval in seconds>  --repoUrl=<link to git repo> --path=<path to track files in>`

For push-based (currently supports only Gitea webhooks):

`thesisCD gitwebhook --repoUrl=<link to git repo> --path=<path to track files in>`
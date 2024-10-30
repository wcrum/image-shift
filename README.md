# Imageswap Webhook

![meme](./imageswap-meme.jpg)

There is a trend within airgapped environments administrators have to manually tag / re-tag images to allow their cluster to pull from a local registry.

Syncronization tools like Hauler, mind the gap and others do an AMAZING job in syncing and hosting required images, but at the cluster level there is still a need to retag. People often do this via configuration via containerd, but in some cases that is not possible.

There are solutions for this problem. These solutions are un-maintained, have weird syntax (my opinion), and are not written in a language that supports easy re-compilation to FIPS modules.

This repository fixes all of that.

## Imageswap Configuration
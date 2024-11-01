# ImageShift Webhook

![meme](./images/imageswap-meme.jpg)

[![Docker Image Size](https://badgen.net/docker/size/wcrum/webhook/latest/arm64?icon=docker&label=Image%20Size)](https://hub.docker.com/r/wcrum/webhook/)

There is a trend within airgapped environments administrators have to manually tag / re-tag images to allow their cluster to pull from a local registry.

Syncronization tools like Hauler, mind the gap and others do an AMAZING job in syncing and hosting required images, but at the cluster level there is still a need to retag. People often do this via configuration via containerd, but in some cases that is not possible.

There are solutions for this problem. These solutions are un-maintained, have weird syntax (my opinion), and are not written in a language that supports easy re-compilation to FIPS modules.

This repository fixes all of that.

## ImageShift Architecture

ImageShift is a Kubernetes MutatingWebhook, the MutatingWebhook patches requests to / from the kube-apiserver. Based on your configuration, ImageShift will patch the image requested with whatever you have configured. 

What this means for you? Dont change your manifests, helm charts or refernces, keep them the same across the board. As you move your images from one domain to another, you dont have to worry about the manifests.

![arch](./images/image.png)


## Stuff to do

- [] Init container patch of mutating webhook
- [] Go lang tests
- [] Helm charts
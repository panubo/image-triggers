# Image Triggers

This go script can consume an AWS SQS queue with ECR Image events on it and trigger an external script with the image name and image tag as parameters.

## Example

```
./image-triggers -queue-name ecr -region ap-southeast-2 -- ./myscript.sh
```

## Releases

Releases are published on GitHub and Docker images are pushed to [Docker Hub](https://hub.docker.com/r/panubo/image-triggers) and [Quay.io](https://quay.io/panubo/image-triggers).

```
docker pull quay.io/panubo/image-triggers:0.0.1
```

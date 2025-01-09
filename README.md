# Image Triggers

This go application can consume an AWS SQS queue with ECR Image events on it and trigger an external script with the image name and image tag as parameters.

## Usage Example

```
./image-triggers -queue-name ecr -region ap-southeast-2 -- ./myscript.sh
```

## Releases

Releases are published on GitHub and Docker images are pushed to [Quay.io](https://quay.io/panubo/image-triggers) and [Amazon ECR Public](https://gallery.ecr.aws/panubo/image-triggers).

```
docker pull quay.io/panubo/image-triggers:0.0.4
docker pull public.ecr.aws/panubo/image-triggers:0.0.4
```

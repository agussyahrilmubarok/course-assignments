# Shelf Blog API

This codebase was created to demonstrate a fully fledged backend application built with [Dart-Shelf](https://github.com/dart-lang/shelf) including CRUD operations, authentication, routing, pagination, and more. I've gone to great lengths to adhere to the [Dart](https://dart.dev/community/) community style guides & best practices.

## How it works

This is a monolithic application structures by components. It uses:

- [Dart](https://dart.dev/) a client-optimized language for fast apps on any platform.
- [Shelf](https://pub.dev/packages/shelf/) a server app using package:shelf.
- [PostgreSQL](https://www.postgresql.org/) as the database, and uses some of it's specific features such as Arrays.
- [SQL](https://en.wikipedia.org/wiki/SQL/) instead of an ORM.
- [Docker Compose](https://docs.docker.com/compose/) to run tests and the application locally. For the testing strategy, I opted to not use mocks and the such to unit test individual functions: instead, I applied a [honeycomb testing strategy](https://www.oreilly.com/library/view/hands-on-microservices/9781789133608/7c9f1260-b0c5-4416-816f-1cad140b56dd.xhtml) to run the tests against the actual app connected with an actual database. I felt that this:
    - Gave me more confidence that the application is working correctly.
    - Allowed me to refactor more easily as I test only the public interface and not the implementation details. This [talk by Dan Abramov](https://www.deconstructconf.com/2019/dan-abramov-the-wet-codebase), for example, illustrates very well the importance of this.
- [Github Actions](https://docs.github.com/en/actions) to run tests on Pull Requests and Merges to the `main` branch.

For more detailed information in [`pubspec.yml`](./pubspec.yaml)

## Getting Started

- Install Dart SDK.
- Clone this project.
- On terminal project run `dart pub get`.
- Create a `.env` file according to the example.


```bash
## RUN THE APP
$ ./dev.sh

## RUN THE TEST
$ ./test.sh
```

## Deployment

Deploy to Google Cloud Platform (GCP)

1. Install [terraform](https://www.terraform.io/).
2. Install the [gcloud-cli](https://cloud.google.com/sdk/docs/install).
3. Login with [gcloud](https://cloud.google.com/sdk/gcloud/reference/auth/login).
4. Go to the [dev bootstrap directory](./deploy/gcp/terraform/environments/dev/bootstrap), change the `locals.tf` values to your own and run `terraform apply`.
5. Go the scripts directory and run build-and-push-container-image.sh `<YOUR_PROJECT_ID> <YOUR_ARTIFACT_REPOSITORY>`.
6. Go to the dev [directory](./deploy/gcp/terraform/environments/dev/), change the locals.tf values to your own and run terraform apply.

## Docker Image

Visit our public repository [here](https://hub.docker.com/r/superosystem/shelf-blog-api)
```bash
docker push superosystem/shelf-blog-api
```

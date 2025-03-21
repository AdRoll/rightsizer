# rightsizer

`rightsizer` allows you to compute the appropriate size of your ECS services
based on the average resource usage in a given period of time. You can use the
output of the program to update your task definitions and keep your services on
the perfect size.

## Usage

In order to use `rightsizer` you need to be logged in to your `AWS Cli` using
[aws config files][aws-config] or a credential manager like our own
[hologram][hologram]. Then you invoque the `rightsizer` command using your
target **cluster** and **service**:

```sh
$ rightsizer some-cluster some-service
containerDefinitions:
  some-database:
    cpu: 1
    memory: 500
    memoryReservation: 95
  some-nginx-thing:
    cpu: 1
    memory: 500
    memoryReservation: 95
  some-language-api:
    cpu: 1
    memory: 500
    memoryReservation: 95
```

The program will output the suggested configuration for your service based on
your actual usage. Then you can use this output to patch your service or patch
your deploy configuration.

## Getting rightsizer

The best way to get `rightsizer` is from dockerhub:

```sh
docker pull nextroll/rightsizer
```

Then you can run the container with the following command:

```sh
docker run --rm nextroll/rightsizer rightsizer some-cluster some-service
```

## Building rightsizer

If you want to build `rightsizer` from source you can use the following,

```sh
go build .
```

This will generate a binary called `rightsizer` that you can use to run the
program.

## Testing rightsizer

We use [gomock][gomock] to generate mocks for our custom clients. If you update
or create a new client, you need to regenerate the mocks. You can do this by
running:

```sh
go generate ./...
```

Then you can run the tests with the following command:

```sh
go test ./...
```

[aws-config]: https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html "AWS Config"
[hologram]: https://github.com/AdRoll/hologram "Hologram"
[gomock]: https://github.com/uber-go/mock "gomock"

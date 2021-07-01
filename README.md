# Metis Server

Neural Network Design Tool

## Developing Metis Server

For building Metis Server, You'll first need Go installed (version 1.14+ is required). Make sure you have Go properly installed, including setting up your GOPATH. Then download a pre-built binary from release page and install the protobuf compiler (version 3.14.0+ is required).

Next, clone this repository into some local directory and then just type `make build`. In a few moments, you'll have a working `metis` executable:
```
$ make build
...
$ bin/metis
```

Tests can be run by typing `make test`.

*NOTE: `make test` includes integration tests that require local applications
 such as MongoDB. To start them, type `docker-compose -f
  docker/docker-compose-ci.yml up --build -d`.*

## For contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches and the contribution workflow.

## License

```
Copyright 2021-present NAVER Corp.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

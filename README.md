# Cycling Coach Lab

This is a Go project that uses the [Echo](https://echo.labstack.com) framework for building web applications and the [Templ](https://templ.guide) package for rendering HTML templates.

## Prerequisites

- Go 1.22.0 or later
- Docker (for building and running the Docker image)
- [migrate](https://github.com/golang-migrate/migrate/tree/master?tab=readme-ov-file) (for running database migrations)
- [direnv](https://direnv.net) (for loading automatically environment variables from the .env file)


## Setup

1. Make sure the your Go setup is completed:
 Paste the following in your `$HOME/.zshrc`
 ```sh
 export GOPATH=$HOME/go
 export PATH=$GOPATH/bin:$PATH
 ```

2. Install the dependencies:
```sh
go install github.com/cosmtrek/air@latest
go install github.com/a-h/templ/cmd/templ@latest
go install go.uber.org/mock/mockgen@latest
make init
```

3. Setup you `.env` file. See [.env.template](.env.template)

## Local Development
To run the server and watch for file changes simultaneously, follow these steps:

1. Open two separate terminal tabs.
2. In the first tab, start the server on port 3000 with the following command:
```sh
make start
```
3. In the second tab, generate the styles.css file by watching for file changes:
```sh
make tailwind
```

## Run App using Docker
1. First, build the Docker image:
```sh
make docker-build
```

2. Then, run the Docker image:
```sh
make docker-run
```

This command will start the server on port 8080.



## Testing
To run the unit tests:
```sh
make test
```

To get test coverage:
```sh
make test-coverage
```

To generate mocks:
```sh
make generate-dev
```

To update snapshots:
```sh
make update-snapshots
```


## K8s development
To test the kubernetes deployment you can install [minikube](https://minikube.sigs.k8s.io/docs/start/) on your dev maschine.

1. Install minikube on your maschine
2. Init minikube
```sh
minikube start --vm-driver=docker --alsologtostderr
```
3. Test minikube
```sh
kubectl get pods -A
```

## Contributing
Please read [CONTRIBUTION.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.


### License
This project is licensed under the LICENSE - see the [LICENSE](LISCENSE) file for details

# Cycling Coach Lab

This is a Go project that uses the [Echo](https://echo.labstack.com) framework for building web applications and the [Templ](https://templ.guide) package for rendering HTML templates.

## Prerequisites

- Go 1.22.0 or later
- Docker (for building and running the Docker image)
- [migrate](https://github.com/golang-migrate/migrate/tree/master?tab=readme-ov-file)


## Setup

1. Install the dependencies:
```sh
make init
```


## Running the Project
You can run the project in two ways:


### Using Go
- This command will start the server on port 3000.
```sh
make start
```
- This command will start the server with hot reload on port 3000.

### Using Docker
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

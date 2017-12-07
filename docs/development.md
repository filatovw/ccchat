# development process

You need to install

    Go >= 1.9
    deps

Fetch all requirements

    make ensure

Install repository packages

    make install

Run tests

    make test

Build client and server for your local environment.
Static will be collected into `bindata` file, thus every time you changed template you should rebuild artefacts.

Built artefacts will be in `bin/local/`

    make build

Build aftertefacts for `linux`.

Built artefacts will be in `bin/linux`

    make build-linux

Clean environment

    make clean

Run `server`, 2 `clients`, `postgres` and database manager in docker. Useful for hand testing.

    make docker-start

Shutdown test environment

    make docker-stop
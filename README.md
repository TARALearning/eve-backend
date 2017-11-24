# eve

## LICENSE
BSD

## DESCRIPTION
The eve pkg contents are:
    * codebase of the service manager
    * codebase for eve tools
    * codebase of the core services
    * templates for code generation

## BUILD STATUS
* MASTER [![Build Status](https://jenkins.evalgo.de:8443/buildStatus/icon?job=eve%20backend%20pipeline/master)](https://jenkins.evalgo.de:8443/job/eve%20backend%20pipeline/master)

## DOWNLOAD
 [ ![Download](https://api.bintray.com/packages/evalgo/eve-backend/core/images/download.svg) ](https://bintray.com/evalgo/eve-backend/core/_latestVersion)

## TOOLS
* [eve-gen]()
* [eve-setup]()
* [eve-bintray]()
  have a look about the scope and usage for this tool at the bottom of this file

## SERVICES
* [evauth]()
* [evbolt]()
* [evlog]()
* [evschedule]()

## SERVICE MODULES/FEATURES
* debug
* cross origin
* prometheus
* evbolt(with or without authentication)
* evsecret(with or without authentication)
* evsession(with or without authentication)
* evuser(with or without authentication)
* evtoken(with or without authentication)
* evlog(with or without authentication)
* evlogin(with or without authentication)
* evschedule(with or without authentication)
* cookie
* routes/urls

## BUILD DEPENDENCIES
* [github.com/boltdb/bolt](https://github.com/boltdb/bolt)
* [github.com/gorilla/mux](https://github.com/gorilla/mux)
* [github.com/prometheus/client_golang/prometheus](https://github.com/prometheus/client_golang/tree/master/prometheus)
* [github.com/prometheus/client_golang/prometheus/promhttp](https://github.com/prometheus/client_golang/tree/master/prometheus/promhttp)
* [github.com/dchest/uniuri](https://github.com/dchest/uniuri)
* [github.com/mitchellh/go-ps](https://github.com/mitchellh/go-ps)

## BUILD
```bash
    DEPENDENCIES( github.com/boltdb/bolt github.com/gorilla/mux github.com/mitchellh/go-ps github.com/dchest/uniuri github.com/prometheus/client_golang/prometheus github.com/prometheus/client_golang/prometheus/promhttp )
    # get dependencies
    for DEP in "${DEPENDENCIES[@]}"; do
        echo "get dependecy pkg::${DEP}"
        go get -v ${DEP}
    done
    git clone ssh://git@git.evalgo.de:7999/eve/backend.git $GOPATH/src/evalgo.org/eve
    cd $GOPATH/src/evalgo.org/eve
    mkdir dist
    TOOLS=(eve-gen eve-setup eve-bintray)
    for TOOL in "${TOOLS[@]}";do
        go build -o dist/${TOOL} bin/${TOOL}/main.go
    done
    SERVICES=( evauth evbolt evlog evschedule )
    USES=""
    USES_DEFAULT="-use debug"
    for SRV in "${SERVICES[@]}";do
        echo "build ${SRV}"
        if [ "${SRV}" == "evbolt" ];then
            USES=" -use evBoltRoot=. ${USES_DEFAULT}"
        else
            USES="${USES_DEFAULT}"
        fi
        dist/eve-gen generate -service ${SRV} ${USES} -target dist/${SRV}_main.go
        gofmt -w -s dist/${SRV}_main.go
        go build -o dist/${SRV} dist/${SRV}_main.go
        #rm -fv dist/${SRV}_main.go
    done
```

## RUN
```bash
    # start the scheduler
    dist/evschedule http&
    # register all services
    curl -X POST -d "Id=evauth&Cmd=dist/evlog&Args=http" http://127.0.0.1:9091/0.0.1/eve/evschedule
    curl -X POST -d "Id=evauth&Cmd=dist/evbolt&Args=http" http://127.0.0.1:9091/0.0.1/eve/evschedule
    curl -X POST -d "Id=evauth&Cmd=dist/evauth&Args=http" http://127.0.0.1:9091/0.0.1/eve/evschedule
    # start all services
    curl -X PUT -d "Mode=Start.Processes" http://127.0.0.1:9091/0.0.1/eve/evschedule
```

## SETUP
```bash
    go run bin/eve-setup/main.go \
        http://localhost:9092/0.0.1/eve/evbolt  \
        francisc.simon@evalgo.org \
        secret \
        123456789012345678901234567890ab \
        secret.sig.key
```

## EVE-BINTRAY SCOPE
The eve-bintray tool is used for deleting old bintray files before running the jenkins pipeline job which does upload the service and tool executables to bintray

## EVE-BINTRAY USAGE
```bash
    go build -o eve-bintray $GOPATH/src/evalgo.org/eve/bin/eve-bintray/main.go
    chmod +x eve-bintray
    # all the variables used here are defined in the bintray REST API documentation
    # https://bintray.com/docs/api/
    ./eve-bintray \
        ${subject} \
        ${repo} \
        ${package} \
        ${version} \
        ${user} \
        ${api_token}
```

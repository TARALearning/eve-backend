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

| Master   |     Testing      |  Development |
|----------|:-------------:|------:|
| [![Build Status](https://jenkins.campus-one.com:8443/buildStatus/icon?job=eve%20backend%20pipeline/master)](https://jenkins.campus-one.com:8443/job/eve%20backend%20pipeline/master)|  [![Build Status](https://jenkins.campus-one.com:8443/buildStatus/icon?job=eve%20backend%20pipeline/test)](https://jenkins.campus-one.com:8443/job/eve%20backend%20pipeline/test) | [![Build Status](https://jenkins.campus-one.com:8443/buildStatus/icon?job=eve%20backend%20pipeline/dev)](https://jenkins.campus-one.com:8443/job/eve%20backend%20pipeline/dev) |

## RELEASES/VERSION(S)
| Master   |     Testing      |  Development |
|----------|:-------------:|------:|
| 0.0.3    |     0.0.3     |    RC 0.0.4 |

## DOWNLOAD
 [ ![Latest Version](https://api.bintray.com/packages/evalgo/eve-backend/core/images/download.svg) ](https://bintray.com/evalgo/eve-backend/core/_latestVersion)

## CODE coverage
[![codecov](https://codecov.io/gh/evalgo/eve-backend/branch/master/graph/badge.svg)](https://codecov.io/gh/evalgo/eve-backend)

## TOOLS
* [eve]()

## SERVICES
* [evauth]()
* [evbolt]()
* [evlog]()
* [evschedule]()
* [evweb]()

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

## START SERVICES
```bash
    go run bin/eve/main.go \
        setup \
        evschedule
```
## SETUP EVAUTH
```bash
    go run bin/eve/main.go \
        setup \
        evauth \
        http://127.0.0.1:9092/{VERSION}/eve/evbolt  \
        francisc.simon@evalgo.org \
        secret \
        123456789012345678901234567890ab \
        secret.sig.key
```

## BINTRAY
```bash
    # all the variables used here are defined in the bintray REST API documentation
    # https://bintray.com/docs/api/
    go run bin/eve/main.go 
        delete \
        bintray \
        -subject ${subject} \
        -repo ${repo} \
        -rpackage ${package} \
        -version ${version} \
        -username ${user} \
        -password ${api_token} \
        -url https://api.bintray.com
```

#!/usr/bin/env sh

TYPE=""
NO_LINT="0"
NO_TEST="0"
SKIP_TESTS="${SKIP_TESTS}"
RELEASE=${RELEASE_BUILD:-"0"}
RELEASE_TYPE="unknown"
VERSION=""
COMMIT_HASH=""
BUILD_UNIXTIME="${BUILD_UNIXTIME}"
BUILD_DATE="${BUILD_DATE}"
PACKAGE_FILENAME=""
DOCKER_TAG=""

echo_red() {
    printf '\033[31m%s\033[0m\n' "$1"
}

check_dependency() {
    for cmd in $1
    do
        if ! which "$cmd" > /dev/null; then
            echo_red "Error: \"$cmd\" is required."
            exit 127
        fi
    done
}

show_help() {
    cat <<-EOF
ezBookkeeping build script

Usage:
    build.sh type [options]

Types:
    backend                 Build backend binary file
    frontend                Build frontend files
    package                 Build package archive
    docker                  Build docker image

Options:
    -r, --release           Build release (The script will use environment variable "RELEASE_BUILD" to detect whether this is release building by default)
    -o, --output <filename> Package file name (For "package" type only)
    -t, --tag               Docker tag (For "docker" type only)
    --no-lint               Do not execute lint check before building
    --no-test               Do not execute unit testing before building (You can use environment variable "SKIP_TESTS" to skip specified tests)
    -h, --help              Show help
EOF
}

parse_args() {
    if [ "$1" = "backend" ] || [ "$1" = "frontend" ] || [ "$1" = "package" ] || [ "$1" = "docker" ]; then
        TYPE="$1"
        shift 1
    fi

    while [ ${#} -gt 0 ]; do
        case "${1}" in
            --release | -r)
                RELEASE="1"
                ;;
            --output | -o)
                PACKAGE_FILENAME="$2"
                shift
                ;;
            --tag | -t)
                DOCKER_TAG="$2"
                shift
                ;;
            --no-lint)
                NO_LINT="1"
                ;;
            --no-test)
                NO_TEST="1"
                ;;
            --help | -h)
                show_help
                exit 0
                ;;
            *)
                echo_red "Invalid argument: $1"
                show_help
                exit 2
                ;;
        esac

        shift 1
    done

    if [ "$RELEASE" = "0" ]; then
        RELEASE_TYPE="snapshot"
    else
        RELEASE_TYPE="release"
    fi
}

check_type_dependencies() {
    if [ "$TYPE" = "" ]; then
        echo_red "Error: No specified type"
        show_help
        exit 2
    fi

    check_dependency "git"

    if [ "$TYPE" = "backend" ]; then
        check_dependency "go gcc"
    elif [ "$TYPE" = "frontend" ]; then
        check_dependency "node npm"
    elif [ "$TYPE" = "package" ]; then
        check_dependency "go gcc node npm tar"
    elif [ "$TYPE" = "docker" ]; then
        check_dependency "docker"
    fi
}

set_build_parameters() {
    VERSION="$(grep '"version": ' package.json | awk -F ':' '{print $2}' | tr -d ' ' | tr -d ',' | tr -d '"')"
    COMMIT_HASH="$(git rev-parse --short=7 HEAD)"

    if [ -z "$BUILD_UNIXTIME" ]; then
        BUILD_UNIXTIME="$(date '+%s')"
    fi

    if [ -z "$BUILD_DATE" ]; then
        BUILD_DATE="$(date '+%Y%m%d')"
    fi
}

build_backend() {
    echo "Pulling backend dependencies..."
    go get .

    if [ "$NO_LINT" = "0" ]; then
        echo "Executing backend lint checking..."
        go vet -v ./...

        if [ "$?" != "0" ]; then
            echo_red "Error: Failed to pass lint checking"
            exit 1
        fi
    fi

    if [ "$NO_TEST" = "0" ]; then
        echo "Executing backend unit testing..."
        go clean -cache

        if [ -z "$SKIP_TESTS" ]; then
            go test ./... -v
        else
            echo "(Skip unit test \"$SKIP_TESTS\")"
            go test ./... -v -skip "$SKIP_TESTS"
        fi

        if [ "$?" != "0" ]; then
            echo_red "Error: Failed to pass unit testing"
            exit 1
        fi
    fi

    backend_build_extra_arguments="-X main.Version=$VERSION"
    backend_build_extra_arguments="$backend_build_extra_arguments -X main.CommitHash=$COMMIT_HASH"

    if [ "$RELEASE" = "0" ]; then
        backend_build_extra_arguments="$backend_build_extra_arguments -X main.BuildUnixTime=$BUILD_UNIXTIME"
    fi

    echo "Building backend binary file ($RELEASE_TYPE)..."

    CGO_ENABLED=1 go build -a -v -trimpath -ldflags "-w -s -linkmode external -extldflags '-static' $backend_build_extra_arguments" -o ezbookkeeping ezbookkeeping.go
    chmod +x ezbookkeeping
}

build_frontend() {
    echo "Pulling frontend dependencies..."
    npm install

    if [ "$NO_LINT" = "0" ]; then
        echo "Executing frontend lint checking..."
        npm run lint

        if [ "$?" != "0" ]; then
            echo_red "Error: Failed to pass lint checking"
            exit 1
        fi
    fi

    if [ "$NO_TEST" = "0" ]; then
        echo "Executing frontend unit testing..."

        npm run test

        if [ "$?" != "0" ]; then
            echo_red "Error: Failed to pass unit testing"
            exit 1
        fi
    fi

    echo "Building frontend files ($RELEASE_TYPE)..."

    if [ "$RELEASE" = "0" ]; then
        buildUnixTime=$BUILD_UNIXTIME npm run build
    else
        npm run build
    fi
}

build_package() {
    package_file_name="$VERSION";

    if [ "$RELEASE" = "0" ]; then
        package_file_name="$package_file_name-$BUILD_DATE"
    fi

    package_file_name="ezbookkeeping-$package_file_name-$(arch).tar.gz"

    if [ -n "$PACKAGE_FILENAME" ]; then
        package_file_name="$PACKAGE_FILENAME"
    fi

    echo "Building package archive \"$package_file_name\" ($RELEASE_TYPE)..."

    build_backend
    build_frontend

    rm -rf package
    mkdir package
    mkdir package/data
    mkdir package/storage
    mkdir package/log
    cp ezbookkeeping package/
    cp -R dist package/public
    cp -R conf package/conf
    cp -R templates package/templates
    cp LICENSE package/

    cd package || { echo_red "Error: Build Failed"; exit 1; }
    tar cvzf "../$package_file_name" .
    cd - || return
}

build_docker() {
    docker_tag="$VERSION"

    if [ "$RELEASE" = "0" ]; then
        docker_tag="SNAPSHOT-$BUILD_DATE";
    fi

    docker_tag="ezbookkeeping:$docker_tag"

    if [ -n "$DOCKER_TAG" ]; then
        docker_tag="$DOCKER_TAG"
    fi

    echo "Building docker image \"$docker_tag\" ($RELEASE_TYPE)..."

    docker build . -t "$docker_tag" --build-arg RELEASE_BUILD=$RELEASE
}

main() {
    if [ -z "$1" ]; then
        show_help
        exit 0
    fi

    parse_args "$@"
    check_type_dependencies "$TYPE"
    set_build_parameters

    if [ "$TYPE" = "backend" ]; then
        build_backend
    elif [ "$TYPE" = "frontend" ]; then
        build_frontend
    elif [ "$TYPE" = "package" ]; then
        build_package
    elif [ "$TYPE" = "docker" ]; then
        build_docker
    fi
}

main "$@"

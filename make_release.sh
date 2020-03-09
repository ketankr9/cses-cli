PROJECT_NAME=cses-cli
RELEASE_VERSION="$1"

if [ -z "$RELEASE_VERSION" ]
then
      RELEASE_VERSION="latest"
fi

RDIR="$(pwd)/release/${RELEASE_VERSION}"

mkdir -p "$RDIR"

cd project
CDIR=`pwd`

for goarch in ""amd64 386""; do
  for goos in ""linux windows darwin""; do
    NAME="${PROJECT_NAME}_${RELEASE_VERSION}_${goos}_${goarch}"
    GOPATH="${CDIR}" CGO_ENABLED=0 GOOS=${goos} GOARCH=${goarch} go build -o ${RDIR}/${NAME}
    echo ${NAME}
  done
done
cd ..

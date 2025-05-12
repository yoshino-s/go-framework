set -e -u -f -o pipefail

DIGESTS="$(ko build ./demo --bare --sbom spdx)"

echo "DIGESTS=$DIGESTS"
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

if cmp -s $1 $2; then
    printf "${GREEN}TEST OK${NC}\n";
else
    printf "${RED}TEST ERROR${NC}\n";
fi

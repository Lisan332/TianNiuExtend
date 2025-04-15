#!/bin/bash
# Script to run TianNiu platform tests

# Set environment variables for testing
export CI=true
export MYSQL_PROD_USERNAME=tianniu_user
export MYSQL_PROD_PASSWORD=tianniu_password
export MYSQL_PALO_PROD_USERNAME=palo_prod_user
export MYSQL_PALO_PROD_PASSWORD=palo_prod_password
export MYSQL_PALO_DEV_USERNAME=palo_dev_user
export MYSQL_PALO_DEV_PASSWORD=palo_dev_password

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Function to run tests
run_tests() {
    local test_file=$1
    local test_name=$2
    
    echo -e "${YELLOW}Running $test_name tests...${NC}"
    
    # Run the test
    go test -v $test_file
    
    # Check the result
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}$test_name tests passed!${NC}"
        return 0
    else
        echo -e "${RED}$test_name tests failed!${NC}"
        return 1
    fi
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed. Please install Go before running tests.${NC}"
    exit 1
fi

# Check if MySQL is running
if ! command -v mysql &> /dev/null || ! mysqladmin ping -h localhost -u root -proot --silent; then
    echo -e "${YELLOW}Warning: MySQL is not running or credentials are incorrect. Database tests will be skipped.${NC}"
    export CI=false
fi

# Navigate to the tests directory
cd "$(dirname "$0")"

# Install dependencies
echo -e "${YELLOW}Installing test dependencies...${NC}"
go get github.com/stretchr/testify/assert
go get github.com/go-sql-driver/mysql
go get gopkg.in/yaml.v2

# Run all tests
echo -e "${YELLOW}Running all tests...${NC}"

# Run deployment tests
run_tests ./deployment_test.go "Deployment"
deployment_result=$?

# Run container tests
run_tests ./container_test.go "Container"
container_result=$?

# Run database tests
run_tests ./database_test.go "Database"
database_result=$?

# Print summary
echo -e "\n${YELLOW}Test Summary:${NC}"
[ $deployment_result -eq 0 ] && echo -e "${GREEN}✓ Deployment tests passed${NC}" || echo -e "${RED}✗ Deployment tests failed${NC}"
[ $container_result -eq 0 ] && echo -e "${GREEN}✓ Container tests passed${NC}" || echo -e "${RED}✗ Container tests failed${NC}"
[ $database_result -eq 0 ] && echo -e "${GREEN}✓ Database tests passed${NC}" || echo -e "${RED}✗ Database tests failed${NC}"

# Exit with error if any test failed
if [ $deployment_result -ne 0 ] || [ $container_result -ne 0 ] || [ $database_result -ne 0 ]; then
    echo -e "\n${RED}Some tests failed!${NC}"
    exit 1
else
    echo -e "\n${GREEN}All tests passed!${NC}"
    exit 0
fi

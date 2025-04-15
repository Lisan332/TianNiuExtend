#!/bin/bash
# TianNiu Platform API - Resource Management Example
# This script demonstrates how to use the TianNiu Platform API to manage resources.

# Configuration
API_BASE_URL="https://tianniuprod.baidu.com/api/v1"
CONFIG_PATH="../../config/tianniu-config.yaml"

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo "Error: jq is required but not installed. Please install jq."
    exit 1
fi

# Check if yq is installed
if ! command -v yq &> /dev/null; then
    echo "Error: yq is required but not installed. Please install yq."
    exit 1
fi

# Load configuration
load_config() {
    if [ ! -f "$CONFIG_PATH" ]; then
        echo "Error: Configuration file not found at $CONFIG_PATH"
        exit 1
    fi
    
    # Get default environment
    DEFAULT_ENV=$(yq eval '.environments[] | select(.default == true) | .name' "$CONFIG_PATH")
    if [ -z "$DEFAULT_ENV" ]; then
        echo "Error: No default environment found in configuration"
        exit 1
    fi
    
    # Get API endpoint for default environment
    API_ENDPOINT=$(yq eval '.environments[] | select(.name == "'"$DEFAULT_ENV"'") | .api_endpoint' "$CONFIG_PATH")
    if [ -z "$API_ENDPOINT" ]; then
        echo "Error: API endpoint not found for environment $DEFAULT_ENV"
        exit 1
    fi
    
    # Get API key environment variable name
    API_KEY_ENV=$(yq eval '.environments[] | select(.name == "'"$DEFAULT_ENV"'") | .auth.api_key_env' "$CONFIG_PATH")
    if [ -z "$API_KEY_ENV" ]; then
        echo "Error: API key environment variable not found for environment $DEFAULT_ENV"
        exit 1
    fi
    
    # Get API key from environment variable
    API_KEY=${!API_KEY_ENV}
    if [ -z "$API_KEY" ]; then
        echo "Error: API key not found in environment variable $API_KEY_ENV"
        exit 1
    fi
    
    echo "Using environment: $DEFAULT_ENV"
    echo "API endpoint: $API_ENDPOINT"
}

# Make API request
api_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    
    local url="${API_ENDPOINT}${endpoint}"
    local headers=(-H "Authorization: Bearer $API_KEY" -H "Accept: application/json")
    
    if [ "$method" = "POST" ] || [ "$method" = "PUT" ]; then
        headers+=(-H "Content-Type: application/json")
    fi
    
    if [ -n "$data" ]; then
        curl -s -X "$method" "${headers[@]}" -d "$data" "$url"
    else
        curl -s -X "$method" "${headers[@]}" "$url"
    fi
}

# Get resource quotas
get_quotas() {
    local namespace=$1
    local endpoint="/resources/quotas"
    
    if [ -n "$namespace" ]; then
        endpoint="$endpoint?namespace=$namespace"
    fi
    
    local response=$(api_request "GET" "$endpoint")
    echo "$response" | jq .
}

# Update resource quotas
update_quotas() {
    local namespace=$1
    local cpu_limit=$2
    local memory_limit=$3
    
    if [ -z "$namespace" ] || [ -z "$cpu_limit" ] || [ -z "$memory_limit" ]; then
        echo "Error: namespace, CPU limit, and memory limit are required"
        exit 1
    fi
    
    local data='{
        "quotas": [
            {
                "resource_type": "cpu",
                "limit": '"$cpu_limit"'
            },
            {
                "resource_type": "memory",
                "limit": '"$memory_limit"'
            }
        ]
    }'
    
    local response=$(api_request "PUT" "/resources/quotas/$namespace" "$data")
    echo "$response" | jq .
}

# Get nodes
get_nodes() {
    local status=$1
    local role=$2
    local endpoint="/resources/nodes"
    local params=()
    
    if [ -n "$status" ]; then
        params+=("status=$status")
    fi
    
    if [ -n "$role" ]; then
        params+=("role=$role")
    fi
    
    if [ ${#params[@]} -gt 0 ]; then
        endpoint="$endpoint?$(IFS=\&; echo "${params[*]}")"
    fi
    
    local response=$(api_request "GET" "$endpoint")
    echo "$response" | jq .
}

# Get node details
get_node() {
    local node_id=$1
    
    if [ -z "$node_id" ]; then
        echo "Error: node ID is required"
        exit 1
    fi
    
    local response=$(api_request "GET" "/resources/nodes/$node_id")
    echo "$response" | jq .
}

# Cordon node
cordon_node() {
    local node_id=$1
    
    if [ -z "$node_id" ]; then
        echo "Error: node ID is required"
        exit 1
    fi
    
    local response=$(api_request "POST" "/resources/nodes/$node_id/cordon")
    echo "$response" | jq .
}

# Uncordon node
uncordon_node() {
    local node_id=$1
    
    if [ -z "$node_id" ]; then
        echo "Error: node ID is required"
        exit 1
    fi
    
    local response=$(api_request "POST" "/resources/nodes/$node_id/uncordon")
    echo "$response" | jq .
}

# Drain node
drain_node() {
    local node_id=$1
    local grace_period=$2
    local force=$3
    
    if [ -z "$node_id" ]; then
        echo "Error: node ID is required"
        exit 1
    fi
    
    local endpoint="/resources/nodes/$node_id/drain"
    local params=()
    
    if [ -n "$grace_period" ]; then
        params+=("grace_period=$grace_period")
    fi
    
    if [ "$force" = "true" ]; then
        params+=("force=true")
    fi
    
    if [ ${#params[@]} -gt 0 ]; then
        endpoint="$endpoint?$(IFS=\&; echo "${params[*]}")"
    fi
    
    local response=$(api_request "POST" "$endpoint")
    echo "$response" | jq .
}

# Get resource usage
get_resource_usage() {
    local namespace=$1
    local period=$2
    local endpoint="/resources/usage"
    local params=()
    
    if [ -n "$namespace" ]; then
        endpoint="/resources/usage/$namespace"
    fi
    
    if [ -n "$period" ]; then
        params+=("period=$period")
    fi
    
    if [ ${#params[@]} -gt 0 ]; then
        endpoint="$endpoint?$(IFS=\&; echo "${params[*]}")"
    fi
    
    local response=$(api_request "GET" "$endpoint")
    echo "$response" | jq .
}

# Get resource recommendations
get_recommendations() {
    local namespace=$1
    local resource_type=$2
    local endpoint="/resources/recommendations"
    local params=()
    
    if [ -n "$namespace" ]; then
        params+=("namespace=$namespace")
    fi
    
    if [ -n "$resource_type" ]; then
        params+=("resource_type=$resource_type")
    fi
    
    if [ ${#params[@]} -gt 0 ]; then
        endpoint="$endpoint?$(IFS=\&; echo "${params[*]}")"
    fi
    
    local response=$(api_request "GET" "$endpoint")
    echo "$response" | jq .
}

# Apply resource recommendation
apply_recommendation() {
    local recommendation_id=$1
    
    if [ -z "$recommendation_id" ]; then
        echo "Error: recommendation ID is required"
        exit 1
    fi
    
    local response=$(api_request "POST" "/resources/recommendations/$recommendation_id/apply")
    echo "$response" | jq .
}

# Print usage
print_usage() {
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  quotas [namespace]                                Get resource quotas"
    echo "  update-quotas <namespace> <cpu_limit> <mem_limit> Update resource quotas"
    echo "  nodes [status] [role]                             Get nodes"
    echo "  node <node_id>                                    Get node details"
    echo "  cordon <node_id>                                  Cordon node"
    echo "  uncordon <node_id>                                Uncordon node"
    echo "  drain <node_id> [grace_period] [force]            Drain node"
    echo "  usage [namespace] [period]                        Get resource usage"
    echo "  recommendations [namespace] [resource_type]        Get resource recommendations"
    echo "  apply-recommendation <recommendation_id>          Apply resource recommendation"
    echo ""
    echo "Options:"
    echo "  namespace       Namespace name"
    echo "  cpu_limit       CPU limit (cores)"
    echo "  mem_limit       Memory limit (GB)"
    echo "  node_id         Node ID"
    echo "  status          Node status (ready, not_ready, cordoned)"
    echo "  role            Node role (master, worker)"
    echo "  grace_period    Grace period for drain operation (seconds)"
    echo "  force           Force drain operation (true/false)"
    echo "  period          Usage period (hour, day, week, month)"
    echo "  resource_type   Resource type (cpu, memory, storage, network)"
    echo "  recommendation_id Recommendation ID"
}

# Main function
main() {
    if [ $# -lt 1 ]; then
        print_usage
        exit 1
    fi
    
    load_config
    
    local command=$1
    shift
    
    case "$command" in
        quotas)
            get_quotas "$1"
            ;;
        update-quotas)
            update_quotas "$1" "$2" "$3"
            ;;
        nodes)
            get_nodes "$1" "$2"
            ;;
        node)
            get_node "$1"
            ;;
        cordon)
            cordon_node "$1"
            ;;
        uncordon)
            uncordon_node "$1"
            ;;
        drain)
            drain_node "$1" "$2" "$3"
            ;;
        usage)
            get_resource_usage "$1" "$2"
            ;;
        recommendations)
            get_recommendations "$1" "$2"
            ;;
        apply-recommendation)
            apply_recommendation "$1"
            ;;
        *)
            echo "Error: Unknown command '$command'"
            print_usage
            exit 1
            ;;
    esac
}

main "$@"

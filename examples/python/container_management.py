#!/usr/bin/env python3
"""
TianNiu Platform API - Container Management Example
This script demonstrates how to use the TianNiu Platform API to manage containers.
"""

import os
import sys
import json
import requests
import argparse
import yaml
from datetime import datetime

# Configuration
API_BASE_URL = "https://tianniuprod.baidu.com/api/v1"
CONFIG_PATH = "../../config/tianniu-config.yaml"

def load_config():
    """Load TianNiu configuration from YAML file."""
    try:
        with open(CONFIG_PATH, 'r') as f:
            config = yaml.safe_load(f)
        return config
    except Exception as e:
        print(f"Error loading config: {e}")
        sys.exit(1)

def get_api_key():
    """Get API key from environment variable."""
    api_key = os.environ.get("TIANNIU_API_KEY")
    if not api_key:
        print("Error: TIANNIU_API_KEY environment variable not set.")
        sys.exit(1)
    return api_key

def get_headers():
    """Get HTTP headers for API requests."""
    return {
        "Authorization": f"Bearer {get_api_key()}",
        "Content-Type": "application/json",
        "Accept": "application/json"
    }

def list_containers(args):
    """List all containers or filter by status."""
    params = {}
    if args.status:
        params["status"] = args.status
    if args.limit:
        params["limit"] = args.limit
    
    response = requests.get(
        f"{API_BASE_URL}/containers",
        headers=get_headers(),
        params=params
    )
    
    if response.status_code == 200:
        containers = response.json()
        print(json.dumps(containers, indent=2))
        print(f"Total containers: {containers['total']}")
    else:
        print(f"Error: {response.status_code}")
        print(response.text)

def get_container(args):
    """Get details of a specific container."""
    response = requests.get(
        f"{API_BASE_URL}/containers/{args.container_id}",
        headers=get_headers()
    )
    
    if response.status_code == 200:
        container = response.json()
        print(json.dumps(container, indent=2))
    else:
        print(f"Error: {response.status_code}")
        print(response.text)

def create_container(args):
    """Create a new container."""
    # Load container configuration from file if provided
    if args.config_file:
        try:
            with open(args.config_file, 'r') as f:
                container_config = json.load(f)
        except Exception as e:
            print(f"Error loading container config: {e}")
            sys.exit(1)
    else:
        # Use default configuration
        container_config = {
            "name": args.name,
            "image": args.image,
            "labels": {
                "app": args.name,
                "environment": args.environment or "development",
                "created_by": "tianniu-cli"
            },
            "ports": [
                {
                    "internal": 8080,
                    "external": 9000,
                    "protocol": "tcp"
                }
            ],
            "environment_variables": [
                {
                    "name": "LOG_LEVEL",
                    "value": "info"
                }
            ]
        }
    
    response = requests.post(
        f"{API_BASE_URL}/containers",
        headers=get_headers(),
        json=container_config
    )
    
    if response.status_code == 200:
        result = response.json()
        print("Container created successfully:")
        print(json.dumps(result, indent=2))
    else:
        print(f"Error: {response.status_code}")
        print(response.text)

def start_container(args):
    """Start a container."""
    response = requests.post(
        f"{API_BASE_URL}/containers/{args.container_id}/start",
        headers=get_headers()
    )
    
    if response.status_code == 200:
        result = response.json()
        print("Container started successfully:")
        print(json.dumps(result, indent=2))
    else:
        print(f"Error: {response.status_code}")
        print(response.text)

def stop_container(args):
    """Stop a container."""
    params = {}
    if args.timeout:
        params["timeout"] = args.timeout
        
    response = requests.post(
        f"{API_BASE_URL}/containers/{args.container_id}/stop",
        headers=get_headers(),
        params=params
    )
    
    if response.status_code == 200:
        result = response.json()
        print("Container stopped successfully:")
        print(json.dumps(result, indent=2))
    else:
        print(f"Error: {response.status_code}")
        print(response.text)

def restart_container(args):
    """Restart a container."""
    params = {}
    if args.timeout:
        params["timeout"] = args.timeout
        
    response = requests.post(
        f"{API_BASE_URL}/containers/{args.container_id}/restart",
        headers=get_headers(),
        params=params
    )
    
    if response.status_code == 200:
        result = response.json()
        print("Container restarted successfully:")
        print(json.dumps(result, indent=2))
    else:
        print(f"Error: {response.status_code}")
        print(response.text)

def delete_container(args):
    """Delete a container."""
    params = {}
    if args.force:
        params["force"] = "true"
    if args.remove_volumes:
        params["remove_volumes"] = "true"
        
    response = requests.delete(
        f"{API_BASE_URL}/containers/{args.container_id}",
        headers=get_headers(),
        params=params
    )
    
    if response.status_code == 200:
        result = response.json()
        print("Container deleted successfully:")
        print(json.dumps(result, indent=2))
    else:
        print(f"Error: {response.status_code}")
        print(response.text)

def get_container_logs(args):
    """Get logs from a container."""
    params = {}
    if args.tail:
        params["tail"] = args.tail
    if args.since:
        params["since"] = args.since
    if args.until:
        params["until"] = args.until
    if args.follow:
        params["follow"] = "true"
        
    response = requests.get(
        f"{API_BASE_URL}/containers/{args.container_id}/logs",
        headers=get_headers(),
        params=params
    )
    
    if response.status_code == 200:
        logs = response.json()
        for log_entry in logs["logs"]:
            timestamp = datetime.fromisoformat(log_entry["timestamp"].replace("Z", "+00:00"))
            stream = log_entry["stream"]
            message = log_entry["message"]
            print(f"[{timestamp}] [{stream}] {message}")
    else:
        print(f"Error: {response.status_code}")
        print(response.text)

def exec_container(args):
    """Execute a command in a container."""
    command = args.command.split()
    
    payload = {
        "command": command,
        "attach_stdout": True,
        "attach_stderr": True
    }
    
    response = requests.post(
        f"{API_BASE_URL}/containers/{args.container_id}/exec",
        headers=get_headers(),
        json=payload
    )
    
    if response.status_code == 200:
        result = response.json()
        print(f"Exit code: {result['exit_code']}")
        if result["stdout"]:
            print("STDOUT:")
            print(result["stdout"])
        if result["stderr"]:
            print("STDERR:")
            print(result["stderr"])
    else:
        print(f"Error: {response.status_code}")
        print(response.text)

def main():
    """Main function to parse arguments and execute commands."""
    parser = argparse.ArgumentParser(description="TianNiu Container Management CLI")
    subparsers = parser.add_subparsers(dest="command", help="Command to execute")
    
    # List containers
    list_parser = subparsers.add_parser("list", help="List containers")
    list_parser.add_argument("--status", help="Filter by status (running, stopped, paused)")
    list_parser.add_argument("--limit", type=int, help="Limit number of results")
    
    # Get container details
    get_parser = subparsers.add_parser("get", help="Get container details")
    get_parser.add_argument("container_id", help="Container ID")
    
    # Create container
    create_parser = subparsers.add_parser("create", help="Create a new container")
    create_parser.add_argument("--name", required=True, help="Container name")
    create_parser.add_argument("--image", required=True, help="Container image")
    create_parser.add_argument("--environment", help="Environment (development, staging, production)")
    create_parser.add_argument("--config-file", help="Path to container configuration JSON file")
    
    # Start container
    start_parser = subparsers.add_parser("start", help="Start a container")
    start_parser.add_argument("container_id", help="Container ID")
    
    # Stop container
    stop_parser = subparsers.add_parser("stop", help="Stop a container")
    stop_parser.add_argument("container_id", help="Container ID")
    stop_parser.add_argument("--timeout", type=int, help="Timeout in seconds")
    
    # Restart container
    restart_parser = subparsers.add_parser("restart", help="Restart a container")
    restart_parser.add_argument("container_id", help="Container ID")
    restart_parser.add_argument("--timeout", type=int, help="Timeout in seconds")
    
    # Delete container
    delete_parser = subparsers.add_parser("delete", help="Delete a container")
    delete_parser.add_argument("container_id", help="Container ID")
    delete_parser.add_argument("--force", action="store_true", help="Force deletion")
    delete_parser.add_argument("--remove-volumes", action="store_true", help="Remove associated volumes")
    
    # Get container logs
    logs_parser = subparsers.add_parser("logs", help="Get container logs")
    logs_parser.add_argument("container_id", help="Container ID")
    logs_parser.add_argument("--tail", type=int, help="Number of lines to show from the end")
    logs_parser.add_argument("--since", help="Show logs since timestamp (ISO 8601)")
    logs_parser.add_argument("--until", help="Show logs until timestamp (ISO 8601)")
    logs_parser.add_argument("--follow", action="store_true", help="Follow log output")
    
    # Execute command in container
    exec_parser = subparsers.add_parser("exec", help="Execute command in container")
    exec_parser.add_argument("container_id", help="Container ID")
    exec_parser.add_argument("command", help="Command to execute")
    
    args = parser.parse_args()
    
    # Load configuration
    config = load_config()
    
    # Execute command
    if args.command == "list":
        list_containers(args)
    elif args.command == "get":
        get_container(args)
    elif args.command == "create":
        create_container(args)
    elif args.command == "start":
        start_container(args)
    elif args.command == "stop":
        stop_container(args)
    elif args.command == "restart":
        restart_container(args)
    elif args.command == "delete":
        delete_container(args)
    elif args.command == "logs":
        get_container_logs(args)
    elif args.command == "exec":
        exec_container(args)
    else:
        parser.print_help()

if __name__ == "__main__":
    main()

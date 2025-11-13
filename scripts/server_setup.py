#!/bin/python3

import os
from pathlib import Path
import platform
import shlex
import subprocess
import argparse

GO_VER = "1.23.4"

def systemd():
    print("systemd integration coming soon!")

def install_go():
    """Download, extract, and install the golang tarball"""
    # install golang tarball
    subprocess.run(["curl", "-L", "-O", "https://go.dev/dl/go1.23.4.linux-amd64.tar.gz"], check=True)
    subprocess.run(["rm", "-rf", "/usr/local/go"], check=True)
    subprocess.run(["tar", "-C", "/usr/local", "-xzf", "go1.23.4.linux-amd64.tar.gz"], check=True)
    
    with open("/etc/profile", "a") as f:
        f.write("\n# Add Go to PATH\nexport PATH=$PATH:/usr/local/go/bin\n")
    
    # add go to setup script's PATH
    current_path = os.environ.get('PATH', '')
    if '/usr/local/go/bin' not in current_path:
        os.environ['PATH'] = f"{current_path}:/usr/local/go/bin"

def clone_repo():
    """Latest commit on main"""
    subprocess.run(["git", "clone", "https://github.com/TLop503/LogCrunch.git"])

def handle_certs(cert_path, key_path):
    """Handle certificate and key file setup"""
    # if either cert_path or key_path is empty, generate self-signed certificates
    if not cert_path or not key_path:
        print("Certificate or key path not provided. Generating self-signed certificates...")
        cert_path, key_path = generate_certs()
    else:
        # verify files exist
        if not os.path.isfile(cert_path):
            print(f"Error: Certificate file not found at {cert_path}")
            exit(1)
        
        if not os.path.isfile(key_path):
            print(f"Error: Key file not found at {key_path}")
            exit(1)
    
    return cert_path, key_path

def generate_certs():
    """Generate self-signed certificates in ~/logcrunch_crypto directory"""
    crypto_dir = os.path.expanduser("~/logcrunch_crypto")
    cert_path = os.path.join(crypto_dir, "server.crt")
    key_path = os.path.join(crypto_dir, "server.key")
    
    print(f"Generating self-signed certificates in {crypto_dir}")
    os.makedirs(crypto_dir, exist_ok=True)
    original_cwd = os.getcwd()
    os.chdir(crypto_dir)
    
    try:
        # generate crt/key
        subprocess.run([
            "openssl", "req", "-x509", "-newkey", "rsa:4096", "-keyout", "server.key",
            "-out", "server.crt", "-days", "365", "-nodes", "-subj", "/CN=localhost"
        ], check=True)
    finally:
        # return to original directory
        os.chdir(original_cwd)
    
    print(f"Crypto generated successfully:")
    return cert_path, key_path

def compile_server():
    """Compile the LogCrunch SIEM server"""
    print("Compiling SIEM Server")
    
    server_output_path = os.path.expanduser("~/logcrunch_server")
    logcrunch_dir = None
    # find mod
    # optional TODO: confirm behavior so this isn't necessary
    if os.path.exists("./go.mod"):
        logcrunch_dir = "."
    elif os.path.exists("../go.mod"):
        logcrunch_dir = ".."
    elif os.path.exists("./LogCrunch/go.mod"):
        logcrunch_dir = "./LogCrunch"
    else:
        print("Error: Cannot find LogCrunch go.mod file")
        print("Current directory:", os.getcwd())
        print("Looked for go.mod in current directory and parent directories")
        exit(1)
    
    # change to the module directory
    original_cwd = os.getcwd()
    os.chdir(logcrunch_dir)
    
    try:
        print(f"Changed to module directory: {os.getcwd()}")
        
        # grab those deps w/ go mod tidy
        print("Running go mod tidy...")
        tidy_result = subprocess.run(["go", "mod", "tidy"], capture_output=True, text=True)
        if tidy_result.returncode != 0:
            print("Warning: go mod tidy failed")
            print("STDERR:", tidy_result.stderr)
        
        # build server using module-aware compilation
        print("Building server...")
        compilation_result = subprocess.run([
            "go", "build", "-o", server_output_path, "./server"
        ], capture_output=True, text=True)
        
        # dump debug info if compilation failed
        if compilation_result.returncode != 0:
            print("Error: Go compilation failed!")
            print("Return code:", compilation_result.returncode)
            print("STDOUT:", compilation_result.stdout)
            print("STDERR:", compilation_result.stderr)
            print("Current working directory:", os.getcwd())
            print("Contents of current directory:")
            try:
                for item in os.listdir("."):
                    print(f"  {item}")
            except Exception as e:
                print(f"  Could not list directory: {e}")
            
            if os.path.exists("./server"):
                print("Server directory exists, contents:")
                try:
                    for item in os.listdir("./server"):
                        print(f"  {item}")
                except Exception as e:
                    print(f"  Could not list server directory: {e}")
            else:
                print("Error: ./server directory does not exist!")
            
            exit(1)
    
    finally:
        # return to original dir
        os.chdir(original_cwd)
    
    # check if the bin was created
    if not os.path.isfile(server_output_path):
        print(f"Error: Expected output file not found at {server_output_path}")
        print("Compilation appeared to succeed but no output file was created")
        exit(1)
    
    print("Server built successfully!")
    subprocess.run(["chmod", "+x", server_output_path], check=True)
    
    return server_output_path

def compile_agent():
    """Compile the LogCrunch SIEM agent"""
    print("Compiling SIEM Agent")
    
    agent_output_path = os.path.expanduser("~/logcrunch_agent")
    logcrunch_dir = None
    # find mod (reuse logic from compile_server)
    if os.path.exists("./go.mod"):
        logcrunch_dir = "."
    elif os.path.exists("../go.mod"):
        logcrunch_dir = ".."
    elif os.path.exists("./LogCrunch/go.mod"):
        logcrunch_dir = "./LogCrunch"
    else:
        print("Error: Cannot find LogCrunch go.mod file")
        print("Current directory:", os.getcwd())
        print("Looked for go.mod in current directory and parent directories")
        exit(1)
    
    # change to the module directory
    original_cwd = os.getcwd()
    os.chdir(logcrunch_dir)
    
    try:
        print(f"Building agent from directory: {os.getcwd()}")
        
        # build agent using module-aware compilation
        print("Building agent...")
        compilation_result = subprocess.run([
            "go", "build", "-o", agent_output_path, "./agent"
        ], capture_output=True, text=True)
        
        # dump debug info if compilation failed
        if compilation_result.returncode != 0:
            print("Error: Agent compilation failed!")
            print("Return code:", compilation_result.returncode)
            print("STDOUT:", compilation_result.stdout)
            print("STDERR:", compilation_result.stderr)
            print("Current working directory:", os.getcwd())
            
            if os.path.exists("./agent"):
                print("Agent directory exists, contents:")
                try:
                    for item in os.listdir("./agent"):
                        print(f"  {item}")
                except Exception as e:
                    print(f"  Could not list agent directory: {e}")
            else:
                print("Error: ./agent directory does not exist!")
            
            exit(1)
    
    finally:
        # return to original dir
        os.chdir(original_cwd)
    
    # check if the bin was created
    if not os.path.isfile(agent_output_path):
        print(f"Error: Expected agent output file not found at {agent_output_path}")
        print("Agent compilation appeared to succeed but no output file was created")
        exit(1)
    
    print("Agent built successfully!")
    subprocess.run(["chmod", "+x", agent_output_path], check=True)
    
    return agent_output_path

def setup_agent_config():
    """Set up agent configuration file"""
    config_dir = os.path.expanduser("~/logcrunch_config")
    config_path = os.path.join(config_dir, "agent_config.yaml")
    
    print(f"Setting up agent configuration in {config_dir}")
    os.makedirs(config_dir, exist_ok=True)
    
    # create a basic agent config file
    config_content = """---
Targets:
  - name: Auth
    path: /var/log/auth.log
    severity: low
    custom: false
    module: syslog
  - name: Syslog
    path: /var/log/syslog
    severity: medium
    custom: false
    module: syslog
"""
    
    with open(config_path, "w") as f:
        f.write(config_content)
    
    print(f"Agent configuration created at {config_path}")
    return config_path

def linux(host, cert_path, key_path, port):

    # check if we are in repo; if not, clone!
    current_dir = Path(__file__).resolve().parent
    if current_dir != "LogCrunch":
        clone_repo()

    # check if golang installed
    go_path = subprocess.run(["which", "go"], capture_output=True).stdout
    if go_path != b'':
        go_ver = subprocess.run(["go", "version"], capture_output=True).stdout
        print(f"LogCrunch is built and tested on Go {GO_VER}")
        print(f"You currently have {go_ver}")
        i = input("Would you like to clean-install this version of go? y/n: ")
        if i == "y":
            install_go()
    else:
        print("Installing latest Go...")
        install_go()
        go_path = subprocess.run(["which", "go"], capture_output=True).stdout
        if go_path == b'':
            print("Error! Go installation failed. Aborting...")
            exit()
    
    # Compile both server and agent
    server_output_path = compile_server()
    agent_output_path = compile_agent()
    
    # Setup agent configuration
    agent_config_path = setup_agent_config()
    
    cert_path, key_path = handle_certs(cert_path, key_path)

    # future systemd unit scaffolding
    """
    if (
        os.path.isdir('/usr/lib/systemd/system/') or 
        os.path.isdir('/run/systemd/system') or
        os.path.isdir('/etc/systemd/system/')
    ):
        systemd()
    """

    # Start the server
    print(f"Starting SIEM server on {host}:{port}...")
    start_siem_cmd = f"{server_output_path} {host} {port} {cert_path} {key_path}"
    start_siem_cmd = shlex.split(start_siem_cmd)
    server_process = subprocess.Popen(start_siem_cmd, start_new_session=True)
    
    # Wait a moment for server to start
    print("Waiting for server to start...")
    import time
    time.sleep(3)
    
    # Start the agent
    print(f"Starting SIEM agent connecting to {host}:{port}...")
    start_agent_cmd = f"{agent_output_path} {host} {port} {agent_config_path} n"
    start_agent_cmd = shlex.split(start_agent_cmd)
    agent_process = subprocess.Popen(start_agent_cmd, start_new_session=True)
    
    print("LogCrunch SIEM server and agent started successfully!")
    print(f"Server PID: 
    {print("Note! Certs are not verified when using the wizard, please start agent manually to fix this.")s    start_agent_cmd = f"{agent_output_path} {host} {port} {agent_config_path} n"
er running on {host}:{port}")
    print(f"Agent configuration: {agent_config_path}")
    
    return()

def main():
    # Parse command line arguments
    parser = argparse.ArgumentParser(description='LogCrunch Server Setup Mage')
    parser.add_argument('host', help='Host to bind the server to (e.g., localhost, 0.0.0.0)')
    parser.add_argument('port', type=int, help='Port number for the server')
    parser.add_argument('cert_path', nargs='?', default='', help='Path to the certificate (.crt) file (optional - will auto-generate if not provided)')
    parser.add_argument('key_path', nargs='?', default='', help='Path to the key (.key) file (optional - will auto-generate if not provided)')
    
    
    args = parser.parse_args()
    
    print("Welcome to the LogCrunch Server Setup Mage!") # like a wizard, but simpler!
    print("Please make sure your system has both tar, git, curl, and openssl are properly installed.")
    
    if args.cert_path and args.key_path:
        print(f"Using certificate: {args.cert_path}")
        print(f"Using key: {args.key_path}")
    else:
        print("No certificate/key provided - will auto-generate self-signed certificates in ~/logcrunch_crypto")
    
    print(f"Server will run on {args.host}:{args.port}")
    
    i = input("Note! this mage will automatically install missing dependencies. Press enter to continue, or x to abort!")
    if i == "x":
        return
    if platform.system() == 'Linux':
        linux(args.host, args.cert_path, args.key_path, args.port)
    elif platform.system() == 'Windows':
        print("The LogCrunch server only runs on Linux. Windows agents are coming soon!")
        exit()
    elif platform.system() == 'Darwin':
        print("Mac support is not currently in development, please reach out with questions.")
        exit()

if __name__ == "__main__":
    main()
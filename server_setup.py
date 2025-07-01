#!/bin/python3

import os
from pathlib import Path
import platform
import shlex
import subprocess
import argparse

GO_VER = "1.23.4"
WORKING_DIR = subprocess.run(["pwd"], capture_output=True).stdout

def systemd():
    print("systemd integration coming soon!")

def install_go():
    # install golang tarball
    subprocess.run(["curl", "-O", "https://go.dev/dl/go1.23.4.linux-amd64.tar.gz"])
    subprocess.run(["rm", "-rf", "/usr/local/go"])
    subprocess.run(["tar", "-C", "/usr/local", "xzf", "go1.23.4.linux-amd64.tar.gz"])
    subprocess.run("echo", "PATH=$PATH:/usr/local/go/bin", ">>", "$HOME/.profile")

def clone_repo():
    subprocess.run(["git", "clone", "https://github.com/TLop503/LogCrunch.git"])

def handle_certs(cert_path, key_path):
    """Handle certificate and key file setup"""
    # If either cert_path or key_path is empty, generate self-signed certificates
    if not cert_path or not key_path:
        print("Certificate or key path not provided. Generating self-signed certificates...")
        cert_path, key_path = generate_certs()
    else:
        # Verify certificate file exists
        if not os.path.isfile(cert_path):
            print(f"Error: Certificate file not found at {cert_path}")
            exit(1)
        
        # Verify key file exists
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
    
    # Create directory if it doesn't exist
    os.makedirs(crypto_dir, exist_ok=True)
    
    # Generate crt/key
    subprocess.run([
        "openssl", "req", "-x509", "-newkey", "rsa:4096", "-keyout", "server.key",
        "-out", "server.crt", "-days", "365", "-nodes", "-subj", "/CN=localhost"
    ], check=True)
    
    print(f"Certificates generated successfully:")
    print(f"  Certificate: {cert_path}")
    print(f"  Private key: {key_path}")
    
    return cert_path, key_path

def linux(cert_path, key_path, port):

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
        if go_path != b'':
            print("Error! Go installation failed. Aborting...")
            exit()
    
    print("Compiling SIEM Server")
    print(f"Server will be configured to run on port {port}")
    compilation_output = subprocess.run(["go", "build", "-o", "~/logcrunch_server", "./server/siem_intake_server.go"], capture_output=True).stdout
    if not os.path.isfile("~/logcrunch_server"):
        print("Erorr, compilation failed with output:")
        print(compilation_output)
    print("Server built!")
    subprocess.run("chmod", "+x", "~/logcrunch_server")
    
    cert_path, key_path = handle_certs(cert_path, key_path)


    if (
        os.path.isdir('/usr/lib/systemd/system/') or 
        os.path.isdir('/run/systemd/system') or
        os.path.isdir('/etc/systemd/system/')
    ):
        systemd()
    else :
        start_siem_cmd = f"~/logcrunch_server localhost {port} {cert_path} {key_path}"
        start_siem_cmd = shlex.split(start_siem_cmd)
        p = subprocess.Popen(start_siem_cmd, start_new_session=True)


def main():
    # Parse command line arguments
    parser = argparse.ArgumentParser(description='LogCrunch Server Setup Mage')
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
    
    print(f"Server will run on port: {args.port}")
    
    i = input("Note! this mage will automatically install missing dependencies. Press enter to continue, or x to abort!")
    if i == "x":
        return
    if platform.system() == 'Linux':
        linux(args.cert_path, args.key_path, args.port)
    elif platform.system() == 'Windows':
        print("The LogCrunch server only runs on Linux. Windows agents are coming soon!")
        exit()
    elif platform.system() == 'Darwin':
        print("Mac support is not currently in development, please reach out with questions.")
        exit()

if __name__ == "__main__":
    main()
#!/bin/python3

import os
import platform
import subprocess

GO_VER = "1.23.4"

def systemd():
    print("systemd integration coming soon!")

def install_go():
    # install golang tarball
    subprocess.run(["curl", "-O", "https://go.dev/dl/go1.23.4.linux-amd64.tar.gz"])
    subprocess.run(["rm", "-rf", "/usr/local/go"])
    subprocess.run(["tar", "-C", "/usr/local", "xzf", "go1.23.4.linux-amd64.tar.gz"])
    subprocess.run("echo", "PATH=$PATH:/usr/local/go/bin", ">>", "$HOME/.profile")


def linux():
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
    
    


    if (
        os.path.isdir('/usr/lib/systemd/system/') or 
        os.path.isdir('/run/systemd/system') or
        os.path.isdir('/etc/systemd/system/')
    ):
        systemd()


def main():
    print("Welcome to the LogCrunch Server Setup Mage!") # like a wizard, but simpler!
    print("Please make sure your system has both tar and curl properly installed.")
    i = input("Note! this mage will automatically install missing dependencies. Press enter to continue, or x to abort!")
    if i == "x":
        return
    if platform.system() == 'Linux':
        linux()
    elif platform.system() == 'Windows':
        print("The LogCrunch server only runs on Linux. Windows agents are coming soon!")
        exit()
    elif platform.system() == 'Darwin':
        print("Mac support is not currently in development, please reach out with questions.")
        exit()

if __name__ == "__main__":
    main()
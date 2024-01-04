# Cosmovisor

## What's Cosmovisor

`cosmovisor` is a process manager for Cosmos SDK application binaries that monitors the governance module for incoming chain upgrade proposals. If it sees a proposal that gets approved, cosmovisor can automatically download the new binary, stop the current binary, switch from the old binary to the new one, and finally restart the node with the new binary.

## Setup
## Installation

1. Downloading the Binary:
    - You can download Cosmovisor from the https://github.com/cosmos/cosmos-sdk/releases/tag/cosmovisor%2Fv1.5.0.
2. Using `go install`:
    - To install the latest version of cosmovisor, run the following command:
    `go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest`

    - To install a specific version, you can specify the version:
    `go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@v1.5.0`

## Setting up environmental variables

`cosmovisor` relies on the following environmental variables to work properly:
- `DAEMON_HOME` is the location where upgrade binaries should be kept (e.g. $HOME/.sided).
- `DAEMON_NAME` is the name of the binary itself (eg. sided).
- `DAEMON_ALLOW_DOWNLOAD_BINARIES` (optional, default = false) if set to true will enable auto-downloading of new binaries (for security reasons, this is intended for full nodes rather than validators).
- `DAEMON_RESTART_AFTER_UPGRADE` (optional, default = true) if set to true it will restart the sub-process with the same command line arguments and flags (but new binary) after a successful upgrade. By default, cosmovisor dies afterwards and allows the supervisor to restart it if needed. Note that this will not auto-restart the child if there was an error.
- `DAEMON_POLL_INTERVAL` (optional, default = 300ms) is the interval length for polling the upgrade plan file. The value can either be a number (in milliseconds) or a duration (e.g. 1s).
- `UNSAFE_SKIP_BACKUP` (optional, default = false), if set to true, upgrades directly without performing a backup. Otherwise (false) backs up the data before trying the upgrade. The default value of false is useful and recommended in case of failures and when a backup needed to rollback. We recommend using the default backup option `UNSAFE_SKIP_BACKUP=false`.

To set these variables correctly, we recommend editing the ~/.profile file so that they are loaded when you log into your machine. You can use the following command:

`nano ~/.profile`

We recommend that you configure the following values:

```
export DAEMON_HOME=$HOME/.sided
export DAEMON_NAME=sided
export DAEMON_ALLOW_DOWNLOAD_BINARIES=true
export DAEMON_RESTART_AFTER_UPGRADE=true
export UNSAFE_SKIP_BACKUP=true
```

After making the necessary changes, press `Ctrl+O` to save the file and then `Ctrl+X` to exit the editor. To apply the changes, reload the `~/.profile` file by running:
`source ~/.profile`
You can confirm the values that have been set by executing the following command:
`echo $DAEMON_NAME`
If the output of this command is sided then you're all set to proceed.

## Setting up folder structure
Cosmovisor expects a certain folder structure:

```
.
├── current -> genesis or upgrades/<name>
├── genesis
│   └── bin
│       └── $DAEMON_NAME
└── upgrades
    └── <name>
        └── bin
            └── $DAEMON_NAME
```

You don't need to be concerned about the current folder as it's merely a symbolic link utilized by Cosmovisor. However, the other directories will require configuration, but this process is straightforward:

`mkdir -p $DAEMON_HOME/cosmovisor/genesis/bin && mkdir -p $DAEMON_HOME/cosmovisor/upgrades`

## Copy sided to the genesis folder.

`cp $HOME/go/bin/sided $DAEMON_HOME/cosmovisor/genesis/bin`

You can verify this by running the following command (you should see the version of `sided`):

`cosmovisor version`

## Starting your node

`cosmovisor run start`

## Setup Service

First, create the service file:
`sudo nano /etc/systemd/system/sided.service`

Insert the following content into the sided.service file.

```
[Unit]
Description=SIDE Blockchain Daemon (cosmovisor)
After=network-online.target
​
[Service]
User=ubuntu
ExecStart=/usr/local/bin/cosmovisor run start
Restart=always
RestartSec=3
LimitNOFILE=4096
Environment="DAEMON_NAME=sided"
Environment="DAEMON_HOME=$HOME/.sided"
Environment="DAEMON_ALLOW_DOWNLOAD_BINARIES=false"
Environment="DAEMON_RESTART_AFTER_UPGRADE=true"
Environment="DAEMON_LOG_BUFFER_SIZE=512"
​
[Install]
WantedBy=multi-user.target
```

```
Note: We assume you are using the "ubuntu" user. If you prefer to use a different user, please replace "ubuntu" with your desired username.
```

Secondly, copy `cosmovisor` to `/usr/local/bin`

`sudo cp ~/go/bin/cosmovisor /usr/local/bin`

Finally, enable the service and start it.

```
sudo -S systemctl daemon-reload
sudo -S systemctl enable sided
# check config one last time before starting!
sudo systemctl start sided
```
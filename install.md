Guidance for setting this project up as service on a raspberry pi
- SSH into your raspberry pi
- Clone this repository to your user's home directory
- Compile the go binary `go build -o airshare`
- Create a symbolic link from this directory to /usr/local/bin: `sudo ln -s /home/<USER>/airshare/airshare /usr/local/bin/airshare`
- Create a Systemd Service file: `sudo vim /etc/systemd/system/airshare.service`
```txt
[Unit]
Description=Air Share
After=network.target

[Service]
ExecStart=/usr/local/bin/airshare
WorkingDirectory=/home/<USER>/airshare
StandardOutput=inherit
StandardError=inherit
Restart=always
User=pi

[Install]
WantedBy=multi-user.target
```
- Reload systemd service `sudo systemctl daemon-reload`
- Enable new service `sudo systemctl enable airshare.service`
- Start service `sudo systemctl start airshare.service`
- Confirm service is running `sudo systemctl status airshare.service`
- Check logs if errors: `journalctl -u airshare.service -f`


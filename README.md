# Share & Charge CPO BackOffice API

For any questions, please ask :)

## Usage Guide

The API is available at: ............


~~~~
[GIN-debug] GET    /api/v1/                  --> main.HandleIndex (2 handlers)
[GIN-debug] POST   /api/v1/cpo/create        --> main.HandleCpoCreate (2 handlers)
[GIN-debug] GET    /api/v1/cpo/info          --> main.HandleCpoInfo (2 handlers)
[GIN-debug] GET    /api/v1/wallet/info       --> main.HandleWalletInfo (2 handlers)
[GIN-debug] DELETE /api/v1/s3cr3tReinitf32fdsfsdf98yu32jlkjfsd89yaf98j320j --> main.HandleReinit (2 handlers)
~~~~


## Install Guide

1. Get an Ubuntu Instance
2. Install Golang. Configure Golang's GOROOT, GOPATH.
3. Under your GOPATH (ex: /home/you/go/)

create the directory ~/go/src/github.com/motionwerkGmbH/

into that directory run: git clone git@github.com:motionwerkGmbH/cpo-backend-api.git (remember to have this command work, you need to add your ssh key into github)

4. the share & charge config files are under configs/sc_configs. Also there you'll find a script called copy.sh that will copy this configs to ~/.sharecharge folder!
5. chmod +x copy.sh then ./copy.sh
6. Install all the dependencies of this app with: go get ./...  (it will take ~1 min)

## Running the API Server

Under the cpo-backend-api folder

~~~~
go run *.go
~~~~


## FAQ

1. I want to run it in the background with logging

You can use: nohup ./myexecutable &

or

Supervisor. Here's a config file:

~~~~
[program:myapp]
command=/home/{{ pillar['username'] }}/bin/myapp
autostart=true
autorestart=true
startretries=10
user={{ pillar['username'] }}
directory=/srv/www/myapp/
environment=MYAPP_SETTINGS="/srv/www/myapp/prod.toml"
redirect_stderr=true
stdout_logfile=/var/log/supervisord/myapp.stdout.log
stdout_logfile_maxbytes=50MB
stdout_logfile_backups=10
~~~~


#### Licence Mozilla Public License Version 2.0

why this license ? see https://christoph-conrads.name/why-i-chose-the-mozilla-public-license-2-0/
# Share & Charge CPO BackOffice API

For any questions, please ask :)

## REST endpoints list :crystal_ball:

I suggest you to use POSTMAN. Please ask me (Andy) for the latest POSTMAN collections


## Install Guide :sun_with_face:

### Attention :fire: if you are trying to follow this steps and get stuck at something, it's very important that you update this readme with the fix, so other developers will not encounter the same problem.


1. Get an Ubuntu Instance
2. Install Golang. Configure Golang's GOROOT, GOPATH.

~~~~
cd /tmp
wget -q https://storage.googleapis.com/golang/getgo/installer_linux
chmod +x installer_linux 
./installer_linux 
source $HOME/.bash_profile

echo 'export GOPATH=$HOME/go' >> ~/.bashrc 
echo 'export PATH=${PATH}:${GOPATH}/bin' >> ~/.bashrc 
source ~/.bashrc 

go get github.com/golang/example/hello
test it: ~/go/src/github.com/golang/example/hello$ go run hello.go
~~~~

3. Under your GOPATH (ex: /home/you/go/)

create the directory 

~~~~
go get -u github.com/gin-gonic/gin
cd ~
make the dirs: ~/go/src/github.com/motionwerkGmbH
~~~~

into that directory run: git clone git@github.com:motionwerkGmbH/cpo-backend-api.git (remember to have this command work, you need to add your ssh key into github)

4. the share & charge config files are under configs/sc_configs. Also there you'll find a script called copy.sh that will copy this configs to ~/.sharecharge folder!
5. chmod +x copy.sh then ./copy.sh
6. Create a config file from the example: ``` cp api_config.example.json api_config.json
7. Install all the dependencies of this app with: go get ./...  (it will take ~1 min)


#### Pdf Generation

Google-Chrome-Stable that will run in headless mode needs to be installed. 

```
sudo apt-get install -y libappindicator1 fonts-liberation dbus-x11 xfonts-base xfonts-100dpi xfonts-75dpi xfonts-cyrillic xfonts-scalable
cd /tmp
wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
sudo dpkg -i --force-depends google-chrome-stable_current_amd64.deb  (don't mind the errors)
sudo apt-get -f install
google-chrome-stable -version
```

Test if it works, upload a html file into /tmp and run:

```
google-chrome-stable --headless --disable-gpu --print-to-pdf your_file.html
```



## Configure Share & Charge Core Client

TODO:// update this section of the readme

~~~~
cd sharecharge-api
git branch
npm install
npm run install
npm run start
~~~~

#### Running the API Server

Under the cpo-backend-api folder

~~~~
go run *.go
~~~~


## FAQ :question:

1. I want to run it in the background

~~~~
sudo touch /var/log/backend.log 
sudo touch /var/log/coreclient.log 
sudo chown ubuntu:ubuntu /var/log/backend.log
sudo chown ubuntu:ubuntu /var/log/coreclient.log
~~~~

Supervisor. Here's a config file:

~~~~
[program:backendapi]
user=ubuntu
numprocs=1
command=/home/ubuntu/go/src/github.com/motionwerkGmbH/cpo-backend-api/backend
directory=/home/ubuntu/go/src/github.com/motionwerkGmbH/cpo-backend-api/
autostart=true
autorestart=true
redirect_stderr=true
stdout_logfile=/var/log/backend.log
stdout_logfile_maxbytes=10MB
stdout_logfile_backups=1
~~~~

~~~~
[program:coreclient]
user=ubuntu
numprocs=1
command=npm run start
directory=/home/ubuntu/motionwerk/sharecharge-core-client/
autostart=true
autorestart=true
startsecs=10
startretries=3
redirect_stderr=true
stdout_logfile=/var/log/coreclient.log
stdout_logfile_maxbytes=10MB
stdout_logfile_backups=1
~~~~


### Serving the CPO Frontend though nginx

~~~~
sudo apt-get install nginx
sudo nano /etc/nginx/sites-enabled/default

add >  root /home/ubuntu/motionwerk/cpo_frontend/dist/;
~~~~

#### Licence Mozilla Public License Version 2.0

why this license ? see https://christoph-conrads.name/why-i-chose-the-mozilla-public-license-2-0/

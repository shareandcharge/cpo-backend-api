For generating the invoices you should have installed google-chrome-stable

sudo apt-get install -y libappindicator1 fonts-liberation dbus-x11 xfonts-base xfonts-100dpi xfonts-75dpi xfonts-cyrillic xfonts-scalable
cd /tmp
wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
sudo dpkg -i --force-depends google-chrome-stable_current_amd64.deb  (don't mind the errors)
sudo apt-get -f install
google-chrome-stable -version
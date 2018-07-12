#!/bin/bash

# This will replace the content of .sharecharge folder!

# let'em have colors
end="\033[0m"
red="\033[0;31m"
green="\033[0;32m"

echo -e "${red}Removing ~/.sharecharge folder${end}"

# This is the core of this script!
rm -rf ~/.sharecharge
mkdir ~/.sharecharge
cp -ru * ~/.sharecharge

echo -e "${green}New Configs copied ok to ~/.sharecharge folder${end}"

#!/bin/sh
echo "Press any key to run the script..."
read -n 1 -s
# Build the image
echo "## BUILDING DOCKER IMAGE ##"
docker build -t forum .
echo -e "---------------------\n" 

echo "Press any key to launch the project..."
read -n 1 -s
echo -e "---------------------\n" 
echo "forum RUNNING"
echo "You're welcome to navigate through the website, register, login, create posts and like/dislike"
echo "---------------------" 
echo "Terminal will now display the user-server interaction"
echo "To Terminate the process, press CTRL+C and Docker Image will be stopped"
echo -e "---------------------\n" 
docker run --name forum-container -p 443:443 forum


echo "Press any key delete image and container..."
read -n 1 -s
# Stop and remove image and containers
echo -e "stopping...\n"
docker rmi -f forum

echo -e "deleting..\n"
ID=$(docker ps -a | awk '{print $1}' | head -n 2 | tail -n 1)
docker rm $ID

echo "## IMAGES LIST ##"
docker images
echo "---------------------"

echo "## CONTAINERS LIST ##"
docker ps -a
echo -e "---------------------\n"

echo "All Done, thank you (:"
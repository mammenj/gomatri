# go matrimony
Its a simple matrimony site, which alows user to register and post their ad. Also they will be able to view ads.

### This is a poc for HTMX

### For android build, using termx
env GOOS=android GOARCH=arm64 go build

ngrok config add-authtoken 



ngrok http --domain=kingfish-smart-formally.ngrok-free.app 8080


cd ~/Documents/ngrok
./ngrok http --domain=kingfish-smart-formally.ngrok-free.app 8080


### This is all the files for deployment

cp gomatri ~/Documents/gomatri/
cp matri.db ~/Documents/gomatri/
cp .env ~/Documents/gomatri/
cp authz_model.conf ~/Documents/gomatri/
cp authz_policy.csv ~/Documents/gomatri/
cp -R templates ~/Documents/gomatri
cp -R static ~/Documents/gomatri

### Find pid of port running 
lsof -i :8080

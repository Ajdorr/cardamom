set -e
sudo echo

# Git update
cd ~/repos/Cardamom/
git switch master; git fetch; git reset --hard origin/master

# Tag the version
version=$(cat ~/repos/Cardamom/version.txt)
git tag versions/$version -f
git push --tags -f

# Upgrade web
cd ~/repos/Cardamom/web
rm -r build/
yarn build
echo "Upgrading front end web application"
sudo rm -rf /var/www/html/app.cardamom.cooking/*
sudo cp -r build/* /var/www/html/app.cardamom.cooking/
sudo chown -R www-data /var/www/html/app.cardamom.cooking/
sudo chmod 700 /var/www/html/app.cardamom.cooking/

# Upgrade API
cd ~/repos/Cardamom/core

# Database Migrate
echo "Migrating database"
go build bin/migrate/migrate.go
mv migrate ~/core/

# Create the server binary
echo "Upgrading API"
go build bin/server/server.go
sudo mv server /usr/local/bin/cardamom-core
sudo systemctl restart cardamom-core.service
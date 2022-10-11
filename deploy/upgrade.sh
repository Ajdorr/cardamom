set -e
sudo echo

cd $(dirname $0)
config_file="$(pwd)/config.env"

# Git update
cd ~/repos/Cardamom/
git switch master; git fetch; git reset --hard origin/master

source $config_file
echo "Upgrade Web: $upgrade_web"
echo "Upgrade Database: $upgrade_database"
echo "Upgrade API: $upgrade_api"

# Tag the version
version=$(cat ~/repos/Cardamom/version.txt)
git tag versions/$version -f
git push --tags -f

if [[ $upgrade_web == "yes" ]]; then
  # Upgrade web
  cd ~/repos/Cardamom/web
  rm -r build/
  yarn build
  echo "Upgrading front end web application"
  sudo rm -rf /var/www/html/app.cardamom.cooking/*
  sudo cp -r build/* /var/www/html/app.cardamom.cooking/
  sudo chown -R www-data /var/www/html/app.cardamom.cooking/
  sudo chmod 700 /var/www/html/app.cardamom.cooking/
fi

# Upgrade API
cd ~/repos/Cardamom/core

if [[ $upgrade_db == "yes" ]]; then
  # Database Migrate
  echo "Migrating database"
  go run bin/migration/migrate/migrate.go
fi

if [[ $upgrade_api == "yes" ]]; then
  # Create the server binary
  echo "Upgrading API"
  go build bin/server/server.go
  sudo mv server /usr/local/bin/cardamom-core
  sudo systemctl restart cardamom-core.service
fi
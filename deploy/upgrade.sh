set -e
sudo echo

cd $(dirname $0)

# Git update
cd ~/repos/Cardamom/
git switch master; git fetch; git reset --hard origin/master

version=$(cat ~/repos/Cardamom/version.txt)
source ~/repos/Cardamom/deploy/config.env
echo "Upgrade Web: $upgrade_web"
echo "Upgrade Database: $upgrade_database"
echo "Upgrade API: $upgrade_api"

read -p "Deploy version $version [y/N]: "
if [[ !($REPLY =~ ^[Yy]$) ]]; then
  echo "Aborting upgrade."
  exit 1
fi

# Tag the version
git tag versions/$version -f
git push --tags -f

if [[ $upgrade_web == "yes" ]]; then
  # Upgrade web
  cd ~/repos/Cardamom/web
  rm -r build/ || echo "No build folder found"
  yarn build
  echo "Upgrading front end web application"
  sudo rm -rf /var/www/html/app.cardamom.cooking/*
  sudo cp -r build/* /var/www/html/app.cardamom.cooking/
  sudo chown -R www-data /var/www/html/app.cardamom.cooking/
  sudo chmod 700 /var/www/html/app.cardamom.cooking/
fi

# Upgrade API
cd ~/repos/Cardamom/core

if [[ $upgrade_database == "yes" ]]; then
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
echo "Enter ssh public key:"
read pub_key

ssh-keygen -t rsa -f ~/.ssh/id_rsa -N ''
echo $pub_key > ~/.ssh/authorized_keys
chmod 700 ~/.ssh
chmod 600 ~/.ssh/authorized_keys

# Install go
wget https://go.dev/dl/go1.19.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.19.1.linux-amd64.tar.gz

# Install npm, nodejs and postgres
sudo apt install -y nodejs npm postgresql postgresql nginx certbot python3-certbot-nginx
sudo npm install -g n
sudo n stable
sudo npm install --global yarn

cd /tmp # this is to prevent a warning message about switching to directory when running the 'sudo -u postgres' command
# start postgres service
db_passwd=$(cat /dev/urandom | tr -dc '[:alnum:]' | fold -w 128 | head -n 1)
sudo sed -i -r 's/(local\s*all\s*all\s*)peer/\1scram-sha-256/g' /etc/postgresql/14/main/pg_hba.conf
sudo systemctl start postgresql.service
sudo -u postgres createdb core
sudo -u postgres createuser jaguar
sudo -u postgres psql -c "alter user jaguar with encrypted password '$db_passwd'"
sudo -u postgres psql -c "grant all privileges on database core to jaguar"

# Generate some environment variables
jwt_secret=$(cat /dev/urandom | tr -dc '[:alnum:]' | fold -w 128 | head -n 1)
password_salt=$(cat /dev/urandom | tr -dc '[:alnum:]' | fold -w 128 | head -n 1)
admin_passwd=$(cat /dev/urandom | tr -dc '[:alnum:]' | fold -w 16 | head -n 1)
# Generate environment file
echo 'ENV="develop"' >> ~/core/.env
echo 'DOMAIN="ajdorr.dev"' >> ~/core/.env
echo "DB_PASSWORD=\"$db_passwd\"" >> ~/core/.env
echo "JWT_TOKEN_SECRET=\"$jwt_secret\"" >> ~/core/.env
echo "PASSWORD_SALT=\"$password_salt\"" >> ~/core/.env
echo "ADMIN_USER_EMAIL=\"jarrodstone.92@gmail.com\"" >> ~/core/.env
echo "ADMIN_USER_PASSWORD=\"$admin_passwd\"" >> ~/core/.env
echo 'OAUTH_GITHUB_CLIENT_ID=""' >> ~/core/.env
echo 'OAUTH_GITHUB_CLIENT_SECRET=""' >> ~/core/.env
echo 'OAUTH_FACEBOOK_CLIENT_ID=""' >> ~/core/.env
echo 'OAUTH_FACEBOOK_CLIENT_SECRET=""' >> ~/core/.env
echo 'OAUTH_MICROSOFT_CLIENT_ID=""' >> ~/core/.env
echo 'OAUTH_MICROSOFT_CLIENT_SECRET=""' >> ~/core/.env
echo 'OAUTH_GOOGLE_CREDS_FILEPATH="/home/cardamom/creds/google-creds.json"' >> ~/core/.env

# Set up nginx service
# To list available applications: sudo ufw app list
sudo ufw allow 'OpenSSH'
sudo ufw allow 'Nginx HTTP'
sudo ufw allow 'Nginx HTTPS'
sudo mkdir -p /var/www/cardamom.cooking
sudo chmod 755 /var/www/cardamom.cooking/
# Add cardamom user to www-data
sudo adduser cardamom www-data

# Lets Encrypt
domain_nginx=$(cat << 'EOF'
server {
    listen 80 default_server;
    listen [::]:80 default_server;
    root /var/www/html;
    server_name app.cardamom.cooking;
}
EOF
)
echo "$domain_nginx" | sudo tee /etc/nginx/conf.d/app.cardamom.cooking.conf > /dev/null
sudo nginx -t && sudo nginx -s reload
echo "Registering with Let's Encrypt"
sudo certbot --nginx -d app.cardamom.cooking
# exmaple config should be saved to: /etc/nginx/conf.d/app.cardamom.cooking.conf
# certificate should be saved to: /etc/letsencrypt/live/app.cardamom.cooking/fullchain.pem
# key should be saved to: /etc/letsencrypt/live/app.cardamom.cooking/privkey.pem
# https://www.nginx.com/blog/using-free-ssltls-certificates-from-lets-encrypt-with-nginx/

# Create config files, setup tool
# https://www.digitalocean.com/community/tools/nginx
ls -s /etc
sudo mkdir /var/www/html/app.cardamom.cooking/
sudo chown www-data:www-data /var/www/html/app.cardamom.cooking/

# enable site
sudo ln -s /etc/nginx/sites-available/app.cardamom.cooking /etc/nginx/sites-enabled/
sudo nginx -t && sudo nginx -s reload

# Generate ssh keys
ssh-keygen -t rsa -f ~/.ssh/id_rsa -N ''
echo "Authorize the following ssh key in Github"
cat ~/.ssh/id_rsa.pub


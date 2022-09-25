# Download cardamom
mkdir -p ~/repos
cd ~/repos
git clone git@github.com:Ajdorr/Cardamom

cd ~/repos/Cardamom/core
go build bin/server/server.go
sudo mv server /usr/local/bin/cardamom-core

go build bin/migrate/migrate.go
mv migrate ~/core/

api_service_daemon=$(cat << 'EOF'
[Unit]
Description=Cardamom Core API service
After=network.target
StartLimitIntervalSec=0

[Service]
Type=simple
Restart=always
RestartSec=1
User=cardamom
WorkingDirectory=/home/cardamom/core
ExecStart=/usr/local/bin/cardamom-core
StandardOutput=file:/var/log/cardamom-core/stdout.log
StandardError=file:/var/log/cardamom-core/stderr.log

[Install]
WantedBy=multi-user.target
EOF
)
echo "$api_service_daemon" | sudo tee /etc/systemd/system/cardamom-core.service > /dev/null
sudo systemctl enable cardamom-core.service
sudo systemctl start cardamom-core.service

echo "You will need to add oauth credentials to ~/core/.env"
echo "Github: https://github.com/settings/developers"
echo "Google: https://console.cloud.google.com/apis/credentials?project=jaguar-cardamom-production"
echo "Facebook: https://developers.facebook.com/apps/"
echo "Microsoft: https://portal.azure.com/#view/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade"
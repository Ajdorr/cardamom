# Update and upgrade
sudo apt update && sudo apt upgrade -y

# Create user
adduser cardamom --gecos 'cardamom,,,,'
usermod -aG sudo cardamom

# Harden SSH
# https://www.digitalocean.com/community/tutorials/how-to-harden-openssh-on-ubuntu-20-04
sed -i -r 's/PermitRootLogin yes/PermitRootLogin no/g' /etc/ssh/sshd_config
sed -i -r 's/#?MaxAuthTries \d+/MaxAuthTries 3/g' /etc/ssh/sshd_config
sed -i -r 's/#?PasswordAuthentication yes/PasswordAuthentication no/g' /etc/ssh/sshd_config
sed -i -r 's/X11Forwarding yes/X11Forwarding no/g' /etc/ssh/sshd_config
systemctl reload sshd.service

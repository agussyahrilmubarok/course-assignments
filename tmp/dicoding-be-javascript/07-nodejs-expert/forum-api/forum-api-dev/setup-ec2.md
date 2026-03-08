# Security

## Memasang Nginx pada EC2 Instance

```bash
# Connect to EC2
ssh -i "<berkas.pem>" <public dns ec2 instance>

# Update, Install, and Check NGINX
sudo apt update
sudo apt-get install nginx -y
sudo systemctl status nginx
```

## Reverse Proxy Server Configuration

```bash
# default configuration
cat /etc/nginx/sites-available/default

# custom configuration
sudo nano /etc/nginx/sites-available/default
## edit
    location / {
        proxy_pass http://localhost:5000; # your app's port
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    
        # …
        # First attempt to serve request as file, then
        # as directory, then fall back to displaying a 404.
        # try_files $uri $uri/ =404;
    }

# restart nginx
sudo systemctl restart nginx
```

## Limit Access

```bash
# custom configuration
sudo nano /etc/nginx/sites-available/default
## edit
    limit_req_zone $binary_remote_addr zone=one:10m rate=30r/m;
    location / {
        ...
        limit_req zone=one;
        
        ...
        # First attempt to serve request as file, then
        # as directory, then fall back to displaying a 404.
        # try_files $uri $uri/ =404;
    }

# restart nginx
sudo systemctl restart nginx
```

## Register [Sub]Domain on Sever NGINX
```bash
# Get subdomain
curl -X POST -H "Content-type: application/json" -d "{ \"ip\": \"<public IP EC2 instance>\" }" "https://sub.dcdg.xyz/dns/records"

# custom configuration
sudo nano /etc/nginx/sites-available/default
## edit
server_name weak-mirrors-sit-quietly.a276.dcdg.xyz www.weak-mirrors-sit-quietly.a276.dcdg.xyz;

```

## Install TLS Certificate

```bash
# Install certbot
sudo apt-get update
sudo apt-get install python3-certbot-nginx -y

# Make ceritificate
sudo certbot --nginx -d yourdomain.com -d yourdomain.com

## Note add your email
```
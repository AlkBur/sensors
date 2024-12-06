# Датчики температыры для proxmox

установим пограмму для считывания сенсоров
```
apt-get install lm-sensors
```

делаем файл исполняемым 
chmod +x /root/sensor
и вешаем на cron
```crontab -e```
```
* * * * * /root/sensor >/dev/null 2>&1
* * * * * (sleep 30 ; /root/sensor) >/dev/null 2>&1
```

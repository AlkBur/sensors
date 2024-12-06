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

Добавм xYFCNHJQRE сенморов в Home Assistant
```
mqtt:
  sensor:
    - name: proxmox_cpu_temp
      state_topic: "homeassistant/sensor/proxmox_system/cpu_temp"
      unit_of_measurement: "°C"
      value_template: "{{ value | round(0) }}"
    - name: proxmox_nvme_temp
      state_topic: "homeassistant/sensor/proxmox_system/nvme_temp"
      unit_of_measurement: "°C"
      value_template: "{{ value | round(0) }}"
```

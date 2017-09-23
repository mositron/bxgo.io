# BXGo.io

![Screenshot](https://github.com/positronth/bxgo.io/raw/master/theme/default/img/screenshot.png)  



ตัวอย่างการเซ็ท proxy pass ผ่าน nginx เพื่อให้สามารถเข้าผ่านโดเมนหรือ port 80 ได้
```
server {
  listen  80;
  server_name  bxgo.io;
  location / {
    proxy_pass http://127.0.0.1:8000;
    ...
  }
}
```


Positron  
Sarawut Chongrakchit [Facebook](https://www.facebook.com/positron.th)

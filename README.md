# BXGo.io
โปรแกรมแสดงราคาเหรียญ BTC ETH DASH XRP OMG และซื้อขายบนเว็บ bx.in.th อัตโนมัติ   

โปรแกรมนี้พัฒนาด้วยภาษา GO ซึ่งสามารถ compile เพื่อนำไปใช้งานได้ทั้ง Windows, macOS และ Linux หรือจะเรียกใช้งานผ่าน Docker ได้ที่ positron/bxgo.io(เร็วๆนี้)

## Download
ดาวน์โหลด โปรแกรมพร้อมใช้ หรือ source code ได้ที่หน้า [https://github.com/positronth/bxgo.io/releases](https://github.com/positronth/bxgo.io/releases)

## ตัวอย่าง
[https://bxgo.io/](https://bxgo.io/)

![Screenshot](https://github.com/positronth/bxgo.io/raw/master/theme/default/img/screenshot.png)  


## config.ini
ในการติดตั้งครั้งแรก จะต้องเพิ่มค่าคอนฟิคให้โปรแกรมคือ  
key + secret ---- ได้จากการสร้างคีย์ใหม่ใน (ยังไม่รองรับ 2FA for API อย่าลืมปิดการใช้งานด้วย)  https://bx.in.th/account  
pass ---- รหัสอะไรก็ได้ สำหรับการยืนยันการสั่งลบการซื้อขาย  

```
# พอร์ทสำหรับหน้าเว็บเพื่อแสดงการทำงาน
# พอร์ทสำหรับหน้าเว็บเพื่อแสดงการทำงาน
port = 8000

# Api Key / อ้างอิงจาก https://bx.in.th/account
key =
secret =
# 2FA for API (ปล่อยว่างหากปิด 2FA for API ไว้)
twofa =

# รหัสผ่านสำหรับสั่งลบคำสั่งซื้อขาย
pass =

# token สำหรับรายงานการซื้อขายไปในกลุ่ม Line (ปล่อยว่างได้)
# รับ token ได้จาก https://notify-bot.line.me/my/
# หลังรับ token อย่าลืมเชิญ LINE Notify เข้ากลุ่ม
line =

# ธีมเว็บสำหรับการแสดงผล
theme = default

# 1 = THB/BTC
# เปิด/ปิดการซื้อขายอัตโนมัติ,  true/false
1_enable = false
# งบในการซื้อแต่ละครั้ง
1_budget = 500.0
# ราคาสูงสุดที่ต้องการซื้อ (ป้องกันดอยยยยย สูงๆ)
1_max_price = 118000.0
# จำนวนยอดสั่งซื้อสูงสุด
1_max_order = 3
# % ของราคาที่จะขาย ลบด้วยราคาซื้อหรือราคาปัจจุบัน (ราคาขาย-ราคาปัจจุบัน)
1_margin = 6
# % ของราคาขายในแต่ละรอบ (ห่างกันขั้นต่ำ)
1_cycle = 4

# 21 = THB/ETH
21_enable = false
21_budget = 500.0
21_max_price = 8500.0
21_max_order = 3
21_margin = 6
21_cycle = 4

# 22 = THB/DAS
22_enable = false
22_budget = 500.0
22_max_price = 11000.0
22_max_order = 3
22_margin = 6
22_cycle = 4

# 25 = THB/XRP
25_enable = false
25_budget = 500.0
25_max_price = 5.5
25_max_order = 3
25_margin = 6
25_cycle = 4

# 26 = THB/OMG
26_enable = false
26_budget = 500.0
26_max_price = 285.0
26_max_order = 3
26_margin = 6
26_cycle = 4
```

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

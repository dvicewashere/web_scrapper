Go Web Scraper

Bu proje, Go (Golang) dili kullanılarak verilen bir web sitesinin

HTML içeriğini,

ekran görüntüsünü,

sayfadaki linkleri

çekip yerel dosyalara kaydeden bir web scraper uygulamasıdır.


Özellikler

+Hedef URL’yi komut satırından alır

+HTTP bağlantı durumunu kontrol eder

+Sayfanın HTML içeriğini kaydeder

+Sayfanın ekran görüntüsünü alır

+Sayfadaki tüm linkleri listeler



Gereksinimler

+Go (Golang)

+Google Chrome

+chromedp kütüphanesi

Kullanım
```bash
go run main.go https://örnekurl.com
```



Program çalıştığında, zaman damgalı bir klasör oluşturur ve aşağıdaki dosyaları kaydeder:

`site_data.html` → Sayfanın HTML içeriği

`output.txt` → HTML içeriği (metin)

`screenshot.png` → Sayfanın ekran görüntüsü

`links.txt` → Sayfada bulunan linkler


// //////////////////////////////////////////////////////// gerekli paketleri içe aktar////////////////////////////////////////////////////////////////////////////
package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Hedef Site URL'si giriniz (Örn: https://harunseker.com): ")
		return
	}

	targetURL := os.Args[1]

	//////////////////////////////////////////////////////////// URL geçerliliğini kontrol et////////////////////////////////////////////////////////////////////////////
	u, err := url.Parse(targetURL)
	if err != nil {
		fmt.Println("Geçersiz URL")
		return
	}

	//////////////////////////////////////////////////////////// Zaman Damgası//////////////////////////////////////////////////////////////////////////////////////////////

	siteName := strings.Replace(u.Host, "www.", "", 1)    // www. kısmını kaldır
	timestamp := time.Now().Format("2006-01-02_15-04-05") // Zaman damgası oluştur

	//////////////////////////////////////////////////////////// klasör Oluştur/////////////////////////////////////////////////////////////////////////////////////////////

	folderName := "dvice_scrapper_files/" + siteName + "_" + timestamp
	err = os.MkdirAll(folderName, 0755)
	if err != nil {
		fmt.Println("Klasör oluşturma hatası:", err)
		return
	}

	//////////////////////////////////////////////////////////// Terminalde Başlık,bağlantı durumu mesajı/////////////////////////////////////////////////////////////////////////////////////////////

	fmt.Println("--- AŞAMA 1: HTTP İlişkisi ---")
	fmt.Println("Scrapper başlatılıyor...")
	fmt.Printf("[*] Bağlantı kuruluyor: %s\n", targetURL)

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var htmlContent string
	var screenshot []byte
	var links []string
	var statusCode int64

	err = chromedp.Run(ctx,
		chromedp.EmulateViewport(1920, 1080),
		chromedp.Navigate(targetURL),
		chromedp.Sleep(2*time.Second),
		chromedp.Evaluate(`window.performance.getEntriesByType('navigation')[0].responseStatus || 200`, &statusCode), // HTTP durum kodunu al
	)

	if err != nil {
		fmt.Println("[!] Bağlantı hatası:", err) // Bağlantı hatası mesajı
		return
	}

	/////////////// HTTP durum kodunu kontrol et ve göster ///////////// /* İleride GUI yapılırsa http.cat sitesinden görsel durum kodları eklenebilir */ ///////////////////
	if statusCode == 200 {
		fmt.Println("[*] Bağlantı başarılı! (200 OK)")
	} else if statusCode == 404 {
		fmt.Println("[!] Bağlantı başarısız! (404 Sayfa Bulunamadı)")
	} else if statusCode == 403 {
		fmt.Println("[!] Bağlantı başarısız! (403 Erişim Engellendi)")
	} else if statusCode == 500 {
		fmt.Println("[!] Bağlantı başarısız! (500 Sunucu Hatası)")
	} else if statusCode == 503 {
		fmt.Println("[!] Bağlantı başarısız! (503 Servis Kullanılamıyor)")
	} else {
		fmt.Printf("[!] HTTP Durum Kodu: %d\n", statusCode)
	}

	fmt.Printf("[1] Sayfanın HTML içeriği 'site_data.html' dosyasına kaydedildi.\n") // HTML içeriği kaydedildi mesajı

	fmt.Println("\n--- AŞAMA 2: Ekran Görüntüsü Alma ---")
	fmt.Println("[*] Chrome başlatılıyor, lütfen bekleyiniz...")
	fmt.Println("[+] Sayfa görüntüsü işleniyor...")

	err = chromedp.Run(ctx,
		chromedp.OuterHTML("html", &htmlContent),
		chromedp.Evaluate(`
			Array.from(document.querySelectorAll("a[href]"))
			.map(a => a.href)
		`, &links),
		chromedp.FullScreenshot(&screenshot, 90), // Ekran görüntüsü alma
	)

	if err != nil {
		fmt.Println("[-] Veri çekme hatası:", err) // Veri çekme hatası
		return
	}

	fmt.Println("[+] Ekran görüntüsü başarıyla 'screenshot.png' dosyasına kaydedildi.")

	//////////////////////////////////////////////////// Dosyaları kaydetme- Herhangi hata esnasında hatanın nerden kaynaklı olduğu belirlenir/////////////////////////////////////////

	if err := os.WriteFile(folderName+"/output.txt", []byte(htmlContent), 0644); err != nil {
		fmt.Println("Dosya yazma hatası (output.txt):", err)
	}
	if err := os.WriteFile(folderName+"/site_data.html", []byte(htmlContent), 0644); err != nil {
		fmt.Println("Dosya yazma hatası (site_data.html):", err)
	}
	if err := os.WriteFile(folderName+"/screenshot.png", screenshot, 0644); err != nil {
		fmt.Println("Dosya yazma hatası (screenshot.png):", err)
	}
	uniqLinks := unique(links)
	if err := os.WriteFile(folderName+"/links.txt", []byte(strings.Join(uniqLinks, "\n")), 0644); err != nil {
		fmt.Println("Dosya yazma hatası (links.txt):", err)
	}

	fmt.Println("\n[+] Tüm görevler tamamlandı.")                                      // Görev tamamlandı mesajı
	fmt.Printf("[+] Dosyalar '%s' klasörüne kaydedildi.\n", folderName)                // Klasör bilgisi
	fmt.Printf("[+] Toplam %d adet link bulundu ve kaydedildi.\n", len(unique(links))) // Bulunan link sayısı
}

func unique(input []string) []string {
	m := make(map[string]bool)
	var result []string
	for _, v := range input {
		if !m[v] {
			m[v] = true
			result = append(result, v)
		}
	}
	return result
}

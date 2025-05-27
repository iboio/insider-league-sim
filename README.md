# Insider League Simulator ⚽

Go ile geliştirilmiş, round-robin fikstür algoritmasına dayalı bir futbol ligi simülasyon uygulaması. Kullanıcılar kendi liglerini oluşturabilir, maçları oynatabilir, şampiyonluk tahminleri alabilir ve maç sonuçlarını düzenleyebilirler. Tüm veriler MySQL veritabanında saklanır; uygulama, bellek üzerinde herhangi bir state tutmaz.

> Backend+Frontend tam entegre çalışır. Frontend UI tasarımı AI destekli üretilmiştir, API düzenlemeleri ve entegrasyon ise elle yazılmıştır.

---

## 🎯 Amaç

Bu proje, Go dilinde context tabanlı state yönetimi, dependency injection (DI) yapısı ve modüler servis mimarisi gibi kavramları derinlemesine deneyimlemek amacıyla geliştirilmiştir. Veritabanı tasarımı bazı yönlerden optimize edilebilir olsa da, tüm işlevsellik eksiksiz çalışmaktadır ve proje, daha iyi bir AppContext mimarisi ile geliştirilmeye devam edecektir.

---

## 🚀 Özellikler

- [x] Lig oluşturma (kendi takım isimlerini seçerek)
- [x] Round-robin fikstür üretimi
- [x] Maç oynatma ve sonuç güncelleme
- [x] Şampiyonluk tahminleri (predict) hesaplama
- [x] Mevcut ligleri görüntüleme ve devam ettirme
- [x] Tüm state'ler MySQL veritabanında JSON formatında saklanır
- [x] Toplam 10 API endpoint: 5 GET, 3 POST, 1 PUT, 1 DELETE
- [x] %75+ unit test coverage (AI destekli test yazımı)
- [x] Dockerfile ve docker-compose.yml ile konteynerleştirme
- [x] Postman collection ve environment dosyaları ile API testi

---

## 🗂️ Proje Yapısı

```
insider-league-sim/
├── backend/
│   ├── internal/
│   │   ├── appcontext/
│   │   ├── generator/
│   │   ├── league/
│   │   ├── predictor/
│   │   ├── repository/
│   │   └── simulation/
│   ├── api/
│   ├── config/
│   ├── main.go
│   └── go.mod/go.sum
├── frontend/
├── postman_collection/
├── Dockerfile
├── docker-compose.yml
└── README.md
```

---

## 🗃️ Veritabanı Yapısı (MySQL)

Tüm veriler `league_sim` adlı veritabanında saklanır.

### `league` tablosu
- `id`, `name`, `leagueId`, `createdAt`

### `active_league` tablosu
- `leagueId`, `teams`, `playedFixtures`, `upcomingFixtures`, `currentWeek`, `standings`, `onActiveLeague`

### `match_results` tablosu
- `homeTeam`, `awayTeam`, `homeGoals`, `awayGoals`, `winnerName`, `matchWeek`

---

## 🔮 Tahmin (Predict) Algoritması

Takımların moral, güç, stamina, savunma gibi istatistikleri dikkate alınarak hesaplanır.  
%40 geçmiş performans + %60 güncel stat çarpanları kullanılır.

---

## 🧪 Testler

- Unit test coverage: %75+
- Test yazımında AI destekli kod üretiminden yararlanılmıştır.
- Predict, League ve Simulation servisleri mocklanarak test edilmiştir.

---

## 📬 API Bilgisi

Toplam 10 endpoint:
- **GET**: 5 adet
- **POST**: 3 adet
- **PUT**: 1 adet
- **DELETE**: 1 adet

Frontend bileşenleri kendi API çağrılarını yapar.  
Postman collection ve env dosyası sayesinde Create League sonrası dönen lig ID'si otomatik olarak test ortamına set edilir.

---

## 🖥️ Frontend

Frontend AI destekli tasarlanmış, API ile entegre edilmiştir.  
Lig kurma, maç oynatma ve tahmin alma işlemleri UI üzerinden yapılabilir.

---

## 🌐 Canlı Demo

Uygulama Raspberry Pi üzerine self-hosted bir şekilde deploy edilmiştir.  
Coolify kullanılarak deploy edilmiş, Cloudflare Tunneling ile yönlendirilmiştir.  
HTTP olarak herkesin erişimine açıktır.

🔗 Uygulamayı deneyimlemek için: [http://iboio.kilicstation.com](http://iboio.kilicstation.com)

---

## 🧠 Kişisel Yorumlar

- ✅ En sevdiğim yön: AppContext yapısı, NestJS geçmişim sayesinde tanıdık ve esnek.
- ⚠️ En çok zorlandığım nokta: Context interface yapısını ilk kez bu kadar kapsamlı uygulamak.
- ❌ Zayıf yön: Veritabanı modeli; JSON alanlar sadeleştirilebilir ve normalize edilebilirdi.

---

## 🛠️ Geliştirme Planları

- Daha modüler AppContext sistemi
- Veri tabanı için daha fazla ilişkisel ve normalleştirilmiş modelleme
- Redis destekli opsiyonel cache katmanı
- Geliştirme "dev" branch'i üzerinden yapılacak, ana branch'e merge edilmeden önce kod gözden geçirme yapılacak.

---

## 📦 Kurulum

```bash
git clone https://github.com/iboio/insider-league-sim.git
cd insider-league-sim
docker-compose up --build
```

> `.env` dosyasını `backend/` klasörü altında kendinize göre düzenlemelisiniz.

---

## 📄 Lisans

MIT

---

Daha fazla bilgi için → [GitHub: iboio/insider-league-sim](https://github.com/iboio/insider-league-sim)
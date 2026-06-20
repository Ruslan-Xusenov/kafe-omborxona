'use client';
import React from 'react';

export default function GuidePage() {
  return (
    <>
      <div className="topbar">
        <h1>📖 Tizimdan Foydalanish Qo'llanmasi</h1>
      </div>
      <div className="page-content fade-in" style={{ maxWidth: 900, margin: '0 auto' }}>
        
        <div className="guide-card">
          <div className="guide-icon">📊</div>
          <div className="guide-text">
            <h2>1. Asosiy oyna (Dashboard)</h2>
            <p>Bu oyna kafening yuragi hisoblanadi. Siz bu yerda kunlik qilingan barcha hisob-kitoblarni bitta ekranda ko'ra olasiz.</p>
            <img src="/guide-dashboard.png" alt="Dashboard" style={{ width: '100%', borderRadius: 8, marginBottom: 16, border: '1px solid var(--border-color)' }} />
            <ul>
              <li><span className="badge badge-green">Kirim / Chiqim</span> Yuqoridagi bloklarda bugungi kunlik sof savdo va qilingan xarajatlar ko'rinadi.</li>
              <li><span className="badge badge-red">Ogohlantirishlar</span> O'ng tomondagi qizil oynada <b>qoldig'i kam qolgan</b> yoki <b>muddati o'tayotgan</b> mahsulotlar ro'yxati chiqadi. Ularga qarab nima sotib olish kerakligini bilib olasiz.</li>
              <li><button className="btn btn-primary btn-sm">📲 Telegramga Yuborish</button> Shu tugmani bossangiz, joriy hisobot telefoningizdagi botga boradi.</li>
            </ul>
          </div>
        </div>

        <div className="guide-card">
          <div className="guide-icon">🏷️</div>
          <div className="guide-text">
            <h2>2. Mahsulotlar Bo'limi</h2>
            <p>Omboringizdagi barcha narsalar shu yerda ro'yxatga olinadi. Har bir mahsulotni yaratishda uning o'lchov birligiga (Litr, Kg, Dona) e'tibor bering.</p>
            <img src="/guide-products.png" alt="Mahsulotlar" style={{ width: '100%', borderRadius: 8, marginBottom: 16, border: '1px solid var(--border-color)' }} />
            <ul>
              <li>Yangi mahsulot qo'shish uchun yuqoridagi <b>"+ Yangi mahsulot"</b> tugmasi bosiladi.</li>
              <li>Siz <b>"Dastlabki miqdor"</b> ga biron raqam yozsangiz, u avtomat omborga qo'shiladi.</li>
              <li><b>"Minimal qoldiq"</b> ga masalan <i>"5"</i> deb yozsangiz, qoldiq 5 tadan kamayib ketganda tizim qizil rangda ogohlantirish beradi.</li>
            </ul>
          </div>
        </div>

        <div className="guide-card">
          <div className="guide-icon">🔄</div>
          <div className="guide-text">
            <h2>3. Tranzaksiyalar (Kirim-Chiqim)</h2>
            <p>Omborga nima kirdi, nima sotildi va nima aynib qolib tashlab yuborildi - barchasi shu yerda qayd etiladi.</p>
            <img src="/guide-transactions.png" alt="Tranzaksiyalar" style={{ width: '100%', borderRadius: 8, marginBottom: 16, border: '1px solid var(--border-color)' }} />
            <div className="transaction-types">
              <div className="t-type"><span className="badge badge-purchase">Kirim (Xarid)</span><br/><small>Bozordan mahsulot kelganda</small></div>
              <div className="t-type"><span className="badge badge-sale">Sotuv (Chiqim)</span><br/><small>Mijozga sotilganda (Ombordan ayiriladi)</small></div>
              <div className="t-type"><span className="badge badge-return">Vozvrat</span><br/><small>Sifatsiz mahsulot orqaga qaytarilganda</small></div>
              <div className="t-type"><span className="badge badge-writeoff">Spisaniye</span><br/><small>Aynib qolgan yoki sinib qolgan narsalar</small></div>
            </div>
            <p style={{marginTop: 15}}>Hisobotni Excelda olish uchun <b>"📥 Excelga Yuklash"</b> tugmasini bosing.</p>
          </div>
        </div>

        <div className="guide-card">
          <div className="guide-icon">🍳</div>
          <div className="guide-text">
            <h2>4. Retseptlar (Kalkulyatsiya)</h2>
            <p>Eng muhim bo'lim! Bitta tayyor kofe sotganingizda ombordan nimalar yechib olinishini shu yerda sozlab qo'yasiz.</p>
            <img src="/guide-recipes.png" alt="Retseptlar" style={{ width: '100%', borderRadius: 8, marginBottom: 16, border: '1px solid var(--border-color)' }} />
            <ul>
              <li>Masalan: <b>"Kapuchino"</b> nomli retsept ochasiz.</li>
              <li>Uning ichiga <i>12 gramm kofe, 0.1 litr sut, 1 dona stakan, 1 dona qopqoq</i> qo'shib qo'yasiz.</li>
              <li>Endi siz bitta "Kapuchino" sotganingizda tizim avtomatik ravishda yuqoridagi 4 ta narsani ombordan keraklicha ayirib tashlaydi!</li>
            </ul>
          </div>
        </div>

        <div className="guide-card">
          <div className="guide-icon">🤝</div>
          <div className="guide-text">
            <h2>5. Qarzlar Boshqaruvi</h2>
            <p>Ta'minotchilardan (Molochniy, Go'shtchi va hokazo) qarzga narsa olganingizda shu yerda yozib borasiz.</p>
            <ul>
              <li>Tizim qaysi ta'minotchiga qancha qarzingiz qolganini qizil rangda ko'rsatib turadi.</li>
              <li>Qarzni uzganingizda <b>"To'lash"</b> tugmasini bosib, to'langan summani yozasiz. Shunda qoldiq kamayadi.</li>
            </ul>
          </div>
        </div>

      </div>

      <style dangerouslySetInnerHTML={{__html: `
        .guide-card {
          background: var(--surface-light);
          border-radius: 16px;
          padding: 24px;
          margin-bottom: 24px;
          display: flex;
          gap: 24px;
          box-shadow: 0 4px 6px -1px rgba(0,0,0,0.1);
          border: 1px solid var(--border-color);
        }
        .guide-icon {
          font-size: 48px;
          background: var(--bg-main);
          width: 80px;
          height: 80px;
          display: flex;
          align-items: center;
          justify-content: center;
          border-radius: 20px;
          flex-shrink: 0;
        }
        .guide-text h2 {
          margin: 0 0 12px 0;
          font-size: 1.25rem;
          color: var(--text-main);
        }
        .guide-text p {
          color: var(--text-muted);
          line-height: 1.6;
          margin-bottom: 16px;
        }
        .guide-text ul {
          padding-left: 20px;
          color: var(--text-main);
          line-height: 1.8;
        }
        .guide-text li {
          margin-bottom: 8px;
        }
        .transaction-types {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
          gap: 12px;
          background: var(--bg-main);
          padding: 16px;
          border-radius: 12px;
        }
        .t-type small {
          color: var(--text-muted);
          display: block;
          margin-top: 4px;
        }
      `}} />
    </>
  );
}

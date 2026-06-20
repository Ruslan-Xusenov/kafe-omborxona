const puppeteer = require('puppeteer');
const { spawn } = require('child_process');

(async () => {
  console.log("Server ishga tushirilmoqda...");
  const server = spawn('node', ['../start.js'], { detached: true });
  await new Promise(r => setTimeout(r, 10000)); // Server yonguncha kutamiz

  const browser = await puppeteer.launch({
    headless: "new",
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  });
  const page = await browser.newPage();
  await page.setViewport({ width: 1280, height: 800 });

  console.log("Sahifaga kirilmoqda...");
  await page.goto('http://localhost:3000/login');
  await new Promise(r => setTimeout(r, 2000));

  console.log("Login qilinmoqda...");
  await page.type('input[type="text"]', 'admin');
  await page.type('input[type="password"]', 'admin');
  await page.click('button[type="submit"]');

  await new Promise(r => setTimeout(r, 4000));
  
  console.log("Dashboard yuklandi.");
  await page.screenshot({ path: 'public/guide-dashboard.png' });

  await page.goto('http://localhost:3000/products');
  await new Promise(r => setTimeout(r, 3000));
  await page.screenshot({ path: 'public/guide-products.png' });

  await page.goto('http://localhost:3000/transactions/new');
  await new Promise(r => setTimeout(r, 3000));
  await page.screenshot({ path: 'public/guide-transactions.png' });

  await page.goto('http://localhost:3000/recipes');
  await new Promise(r => setTimeout(r, 3000));
  await page.screenshot({ path: 'public/guide-recipes.png' });

  await browser.close();
  process.kill(-server.pid); // Kill the server group
  console.log("Tugadi!");
  process.exit(0);
})();

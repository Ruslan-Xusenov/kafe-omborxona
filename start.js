const { spawn } = require('child_process');
const path = require('path');
const os = require('os');

console.log('🚀 Kafe Omborxona tizimini ishga tushirish (Development mode)...\n');

// Pkll port 8080 agar band bo'lsa (faqat Linux/Mac uchun)
if (os.platform() !== 'win32') {
    const { execSync } = require('child_process');
    try {
        execSync('fuser -k 8080/tcp 2>/dev/null');
        execSync('fuser -k 3000/tcp 2>/dev/null');
    } catch (e) {
        // e'tibor bermaymiz
    }
}

// 1. Backendni ishga tushirish (Go)
const backend = spawn('go', ['run', './cmd/api/'], {
  cwd: path.join(__dirname, 'backend'),
  stdio: 'inherit',
  shell: true
});

// 2. Frontendni ishga tushirish (Next.js)
const frontend = spawn('npm', ['run', 'dev'], {
  cwd: path.join(__dirname, 'frontend'),
  stdio: 'inherit',
  shell: true
});

backend.on('error', (err) => console.error('🔴 Backend xatosi:', err));
frontend.on('error', (err) => console.error('🔴 Frontend xatosi:', err));

process.on('SIGINT', () => {
  console.log('\n🛑 Tizim to\'xtatilmoqda...');
  backend.kill('SIGINT');
  frontend.kill('SIGINT');
  process.exit(0);
});

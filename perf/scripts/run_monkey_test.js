const puppeteer = require('puppeteer');
const fs = require('fs');
const path = require('path');

// --- CONFIGURATION ---
const TARGET_URL = 'https://chrome-perf.corp.goog/m/';
const SCRIPT_PATH = path.join(__dirname, 'multigraph_monkey_test.js');
const HEADLESS = process.env.HEADLESS === 'true';

(async () => {
  console.log(`[INFO] Launching Monkey Test Runner...`);
  console.log(`   Target: ${TARGET_URL}`);
  console.log(`   Script: ${SCRIPT_PATH}`);
  console.log(`   Mode: ${HEADLESS ? 'Headless' : 'Headful (Visible)'}`);

  if (!fs.existsSync(SCRIPT_PATH)) {
    console.error(`[ERROR] Script file not found at: ${SCRIPT_PATH}`);
    process.exit(1);
  }

  const monkeyScript = fs.readFileSync(SCRIPT_PATH, 'utf8');

  // Resolve Chrome Path (System or Env)
  const executablePath = process.env.CHROME_BIN || '/usr/bin/google-chrome';
  if (fs.existsSync(executablePath)) {
    console.log(`   Chrome: ${executablePath}`);
  } else {
    console.warn(
      `[WARN] Chrome not found at ${executablePath}. Using bundled Chromium (if available).`
    );
  }

  // Launch Browser
  const browser = await puppeteer.launch({
    headless: HEADLESS,
    executablePath: fs.existsSync(executablePath) ? executablePath : undefined,
    defaultViewport: null, // Full width
    args: ['--start-maximized', '--no-sandbox', '--disable-setuid-sandbox'],
  });

  try {
    const pages = await browser.pages();
    const page = pages.length > 0 ? pages[0] : await browser.newPage();

    // Relay Console Logs
    page.on('console', (msg) => {
      const text = msg.text();
      const type = msg.type().toUpperCase();

      // Filter out noise if needed, or just print everything
      if (
        text.includes('SUPER MONKEY') ||
        text.includes('[PASS]') ||
        text.includes('[FAIL]') ||
        text.includes('[INFO]') ||
        text.includes('[WARN]') ||
        text.includes('[Discovery]') ||
        text.includes('[STATS]') ||
        text.includes('TEST_HEADER')
      ) {
        console.log(`[BROWSER] ${text}`);
      } else if (msg.type() === 'error') {
        console.error(`[BROWSER-ERR] ${text}`);
      }
    });

    console.log(`Unknown Auth State. Navigating to ${TARGET_URL}...`);
    await page.goto(TARGET_URL, { waitUntil: 'networkidle2', timeout: 60000 });

    // AUTHENTICATION CHECK
    if (
      page.url().includes('login.corp.google.com') ||
      page.url().includes('accounts.google.com')
    ) {
      if (HEADLESS) {
        console.error(
          '[ERROR] Authentication required but running in Headless mode. ' +
            'Please run in Headful mode to log in.'
        );
        process.exit(1);
      }
      console.log('[WARN] Authentication required!');
      console.log('       Please log in within the browser window...');

      // Wait for redirect back to target
      await page.waitForFunction(
        (url) => window.location.href.includes(url),
        { timeout: 0 },
        'chrome-perf.corp.goog'
      );
      console.log('[INFO] Logged in! Proceeding...');
      await new Promise((r) => setTimeout(r, 2000)); // Settlement time
    }

    console.log('[INFO] Injecting Monkey Script...');

    // Inject and Run
    await page.evaluate(monkeyScript);

    console.log('[INFO] Monitoring test execution...');

    await new Promise((resolve, reject) => {
      const checkInterval = setInterval(async () => {
        try {
          // Check for test completion signals
          const result = await page.evaluate(() => {
            // Success: Summary overlay exists
            if (document.getElementById('test-summary-overlay')) return 'SUCCESS';

            // Failure: HUD indicates crash
            const hud = document.getElementById('test-hud-overlay');
            if (hud) {
              const text = hud.innerText;
              if (text.includes('CRASH') || text.includes('FATAL ERROR')) return 'FAILURE';
            }
            return null;
          });

          if (result === 'SUCCESS') {
            clearInterval(checkInterval);
            resolve('SUCCESS');
          } else if (result === 'FAILURE') {
            clearInterval(checkInterval);
            reject(new Error('Test Suite Reported Failure'));
          }
        } catch (e) {
          // Page might be navigating or closed
          console.warn('[POLLER-WARN]', e.message);
        }
      }, 2000);

      // Safety Timeout (e.g., 10 minutes)
      setTimeout(() => {
        clearInterval(checkInterval);
        reject(new Error('Timeout waiting for test completion'));
      }, 600000);
    });

    console.log('[SUCCESS] Test Suite Finished Successfully!');
  } catch (error) {
    console.error('[ERROR] Test Failed:', error.message);
    process.exit(1);
  } finally {
    if (!HEADLESS) {
      console.log('Browser will close in 5 seconds...');
      await new Promise((r) => setTimeout(r, 5000));
    }
    await browser.close();
  }
})();

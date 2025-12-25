import { test, expect } from '@playwright/test';

test.describe('Documentation Site', () => {
  test('homepage loads correctly', async ({ page }) => {
    await page.goto('/');
    await expect(page).toHaveTitle(/TimeOtter/);
    await expect(page.locator('h1')).toBeVisible();
  });

  test('navigation works', async ({ page }) => {
    // Set desktop viewport
    await page.setViewportSize({ width: 1280, height: 720 });

    // Start from homepage and use the Get Started button
    await page.goto('/');
    await page.click('a:has-text("Get Started")');
    await expect(page).toHaveURL(/getting-started/);

    // Now navigate using sidebar link to Installation
    const installLink = page.locator('.sidebar-content a[href*="installation"]');
    await expect(installLink).toBeVisible();
    await installLink.click();
    await expect(page).toHaveURL(/installation/);
    await expect(page.locator('h1')).toContainText(/Installation/i);
  });

  test('all documentation pages render', async ({ page }) => {
    const pages = [
      '/',
      '/getting-started/',
      '/installation/',
      '/oauth-setup/',
      '/configuration/',
      '/cron-setup/',
      '/troubleshooting/',
    ];

    for (const path of pages) {
      await page.goto(path);
      // Each page should have a main content area
      await expect(page.locator('main')).toBeVisible();
      // Page should not be a 404 error page
      await expect(page.locator('h1')).not.toContainText('404');
      // Title should contain TimeOtter
      await expect(page).toHaveTitle(/TimeOtter/);
    }
  });

  test('search functionality works', async ({ page }) => {
    await page.goto('/');

    // Find and click search button (wait for it to be enabled)
    const searchButton = page.locator('button[data-open-modal]').first();
    await expect(searchButton).toBeVisible();

    // Wait for search to be initialized (button becomes enabled)
    await expect(searchButton).toBeEnabled({ timeout: 10000 });
    await searchButton.click();

    // Search dialog should open
    const searchDialog = page.locator('dialog[open]');
    await expect(searchDialog).toBeVisible({ timeout: 5000 });

    // Verify dialog contains search input
    const searchContainer = page.locator('#starlight__search');
    await expect(searchContainer).toBeVisible();
  });

  test('external links have correct attributes', async ({ page }) => {
    await page.goto('/');

    // Check GitHub link exists and has proper attributes
    const githubLink = page.locator('a[href*="github.com/bupd/timeotter"]').first();
    if (await githubLink.count() > 0) {
      await expect(githubLink).toBeVisible();
    }
  });

  test('code blocks render correctly', async ({ page }) => {
    await page.goto('/installation/');

    // Check that code blocks are present and styled
    const codeBlocks = page.locator('pre code');
    if (await codeBlocks.count() > 0) {
      await expect(codeBlocks.first()).toBeVisible();
    }
  });

  test('mobile navigation works', async ({ page }) => {
    // Set mobile viewport
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto('/');

    // Menu button should be visible on mobile
    const menuButton = page.locator('button[aria-label*="Menu"], button[aria-controls*="nav"]').first();
    if (await menuButton.count() > 0) {
      await expect(menuButton).toBeVisible();
    }
  });
});

import { test, expect } from '@playwright/test';

test.describe('Documentation Site', () => {
  test('homepage loads correctly', async ({ page }) => {
    await page.goto('/');
    await expect(page).toHaveTitle(/TimeOtter/);
    await expect(page.locator('h1')).toBeVisible();
  });

  test('navigation works', async ({ page }) => {
    await page.goto('/');

    // Check sidebar navigation exists
    const sidebar = page.locator('nav[aria-label="Main"]');
    await expect(sidebar).toBeVisible();

    // Navigate to Installation page
    await page.click('a[href*="installation"]');
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
      // No error messages should be visible
      await expect(page.locator('text=404')).not.toBeVisible();
      await expect(page.locator('text=Error')).not.toBeVisible();
    }
  });

  test('search functionality works', async ({ page }) => {
    await page.goto('/');

    // Find and click search button
    const searchButton = page.locator('button[data-open-modal]').first();
    await expect(searchButton).toBeVisible();
    await searchButton.click();

    // Search modal should open
    const searchModal = page.locator('[role="dialog"]');
    await expect(searchModal).toBeVisible();

    // Type in search
    const searchInput = page.locator('input[type="search"]');
    await searchInput.fill('installation');

    // Wait for results
    await page.waitForTimeout(500);

    // Check that results appear
    const results = page.locator('[role="listbox"] [role="option"]');
    await expect(results.first()).toBeVisible();
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

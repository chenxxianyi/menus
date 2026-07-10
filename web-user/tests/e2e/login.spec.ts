import { expect, test } from '@playwright/test'

test('a user can log in and persist an authenticated session', async ({ page }) => {
  await page.route('**/api/auth/login', async (route) => {
    await route.fulfill({
      contentType: 'application/json',
      body: JSON.stringify({ code: 0, message: 'success', data: { token: 'test-token', user: { id: 1, username: 'tester' } } }),
    })
  })
  await page.route('**/api/user/preferences/status', async (route) => {
    await route.fulfill({ contentType: 'application/json', body: JSON.stringify({ code: 0, message: 'success', data: { completed: true } }) })
  })
  await page.route('**/api/user/info', async (route) => {
    await route.fulfill({
      contentType: 'application/json',
      body: JSON.stringify({ code: 0, message: 'success', data: { id: 1, username: 'tester', nickname: '测试用户' } }),
    })
  })
  await page.route('**/api/home*', async (route) => {
    await route.fulfill({ contentType: 'application/json', body: JSON.stringify({ code: 0, message: 'success', data: { hot_recipes: [] } }) })
  })

  await page.goto('/login')
  await page.getByPlaceholder('用户名').fill('tester')
  await page.getByPlaceholder('密码').fill('password123')
  const [loginResponse] = await Promise.all([
    page.waitForResponse('**/api/auth/login'),
    page.getByRole('button', { name: '登录' }).click(),
  ])
  expect(loginResponse.ok()).toBeTruthy()

  await expect.poll(() => page.evaluate(() => localStorage.getItem('token'))).toBe('test-token')
})

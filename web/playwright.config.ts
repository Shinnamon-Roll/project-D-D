import { defineConfig, devices } from '@playwright/test'

const url = 'http://127.0.0.1:18080'

export default defineConfig({
  testDir: './e2e',
  timeout: 60_000,
  use: { baseURL: url, trace: 'retain-on-failure' },
  projects: [{ name: 'chromium', use: { ...devices['Desktop Chrome'] } }],
  webServer: {
    command: 'cd .. && make build && ./bin/dnd -addr 127.0.0.1:18080',
    url,
    reuseExistingServer: false,
    timeout: 180_000,
  },
})

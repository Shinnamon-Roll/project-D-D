import { defineConfig, devices } from '@playwright/test'

const url = 'http://127.0.0.1:18080'

export default defineConfig({
  testDir: './e2e',
  timeout: 120_000, // one test walks the whole game; ~30s idle, ~60s on a loaded machine
  use: { baseURL: url, trace: 'retain-on-failure' },
  projects: [{ name: 'chromium', use: { ...devices['Desktop Chrome'] } }],
  webServer: {
    // fresh notes dir per run: the spec expects exactly the notes it writes
    command: 'cd .. && make build && ./bin/dnd -addr 127.0.0.1:18080 -data "$(mktemp -d)"',
    url,
    reuseExistingServer: false,
    timeout: 180_000,
  },
})

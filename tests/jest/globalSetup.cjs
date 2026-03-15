const fs = require('fs');
const path = require('path');
const { spawn } = require('child_process');
const axios = require('axios');
require('dotenv').config({ path: path.join(__dirname, '../../.env') });

const statePath = path.join(__dirname, '.server-state.json');

async function waitForServer(baseURL, timeoutMs = 30000) {
  const start = Date.now();

  while (Date.now() - start < timeoutMs) {
    try {
      await axios.get(`${baseURL}/swagger/index.html`, { timeout: 2000 });
      return;
    } catch (error) {
      const status = error?.response?.status;
      if (status && status < 500) {
        return;
      }
    }

    await new Promise((resolve) => setTimeout(resolve, 500));
  }

  throw new Error(`Server did not become ready within ${timeoutMs}ms`);
}

module.exports = async () => {
  const manualServer = process.env.BACKEND_GO_MANUAL_SERVER === 'true';

  if (manualServer) {
    return;
  }

  const rootDir = path.join(__dirname, '../..');
  const port = process.env.BACKEND_TEST_PORT || process.env.PORT || '3001';
  const baseURL = process.env.BACKEND_TEST_BASE_URL || `http://localhost:${port}`;

  const child = spawn('go', ['run', '.'], {
    cwd: rootDir,
    env: {
      ...process.env,
      PORT: String(port),
      SEED: process.env.SEED || 'true',
    },
    detached: true,
    stdio: 'ignore',
  });

  child.unref();

  fs.writeFileSync(
    statePath,
    JSON.stringify({ pid: child.pid, baseURL }, null, 2),
    'utf-8',
  );

  await waitForServer(baseURL);
};

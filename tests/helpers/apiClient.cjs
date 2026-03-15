const fs = require('fs');
const path = require('path');
const axios = require('axios');
require('dotenv').config({ path: path.join(__dirname, '../../.env') });

const statePath = path.join(__dirname, '../jest/.server-state.json');

function getBaseURL() {
  if (process.env.BACKEND_TEST_BASE_URL) return process.env.BACKEND_TEST_BASE_URL;

  if (fs.existsSync(statePath)) {
    const state = JSON.parse(fs.readFileSync(statePath, 'utf-8'));
    if (state?.baseURL) return state.baseURL;
  }

  const port = process.env.BACKEND_TEST_PORT || process.env.PORT || '3001';
  return `http://localhost:${port}`;
}

function createApiClient(token) {
  return axios.create({
    baseURL: getBaseURL(),
    timeout: 15000,
    headers: token
      ? {
          Authorization: `Bearer ${token}`,
        }
      : undefined,
  });
}

module.exports = { createApiClient };

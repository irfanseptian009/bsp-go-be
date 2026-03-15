const fs = require('fs');
const path = require('path');

const statePath = path.join(__dirname, '.server-state.json');

module.exports = async () => {
  if (!fs.existsSync(statePath)) return;

  try {
    const state = JSON.parse(fs.readFileSync(statePath, 'utf-8'));
    if (state?.pid) {
      try {
        process.kill(-state.pid, 'SIGTERM');
      } catch (error) {
        process.kill(state.pid, 'SIGTERM');
      }
    }
  } catch (error) {
    // ignore teardown errors
  } finally {
    fs.rmSync(statePath, { force: true });
  }
};

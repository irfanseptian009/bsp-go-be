module.exports = {
  testEnvironment: 'node',
  testMatch: ['<rootDir>/tests/**/*.test.js'],
  globalSetup: '<rootDir>/tests/jest/globalSetup.cjs',
  globalTeardown: '<rootDir>/tests/jest/globalTeardown.cjs',
  setupFilesAfterEnv: ['<rootDir>/tests/jest/setup.cjs'],
  verbose: true,
};

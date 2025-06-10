module.exports = {
  verbose: true,
  testEnvironment: 'jsdom',
  setupFilesAfterEnv: ['<rootDir>/test/setupTests.ts'],
  testMatch: ['**/test/**/*.test.(js|jsx|ts|tsx)'],
  roots: ['<rootDir>/src', '<rootDir>/test'],
  transform: {
    '^.+\\.(ts|tsx)$': 'ts-jest',
    '^.+\\.(js|jsx)$': 'babel-jest',
  },
  moduleNameMapper: {
    '^@core/(.*)$': '<rootDir>/src/core/$1',
    '^@auth/(.*)$': '<rootDir>/src/auth/$1',
    '^@pages/(.*)$': '<rootDir>/src/pages/$1',
    '^@component/(.*)$': '<rootDir>/src/component/$1',

    '\\.(css|less|scss|sass)$': 'identity-obj-proxy',
    '\\.(gif|ttf|eot|svg|png)$': '<rootDir>/__mocks__/fileMock.js',
  },
  moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx', 'json', 'node'],
  collectCoverage: true,
  collectCoverageFrom: [
    'src/**/*.{js,jsx,ts,tsx}', 
    '!src/index.tsx',
    '!src/react-app-env.d.ts',
    '!src/reportWebVitals.ts',
  ],
  testEnvironmentOptions: {
    url: 'http://app.cardamom.cooking/',
  },
};

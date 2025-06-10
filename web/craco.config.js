const path = require('path');

module.exports = {
  webpack: {
    alias: {
      "@core": path.resolve(__dirname, 'src/core'),
      "@auth": path.resolve(__dirname, "src/auth"),
      "@pages": path.resolve(__dirname, "src/pages"),
      "@component": path.resolve(__dirname, "src/component"),
    },
  },
};

module.exports = {
  apps: [
    {
      name: 'taskly-js',
      script: './server.js',
      instances: 'max',
      exec_mode: 'cluster',
      watch: false,
      env: {
        NODE_ENV: 'production',
        PORT: 8080,
      },
      env_development: {
        NODE_ENV: 'development',
        PORT: 8080,
      },
    },
  ],
};

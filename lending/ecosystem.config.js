module.exports = {
  apps: [
    {
      name: 'xreplyagent-lending',
      script: 'npm',
      args: 'start',
      cwd: '/Users/dogukangundogan/Desktop/Dev/XReplyAgent/lending',
      instances: 1,
      autorestart: true,
      watch: false,
      max_memory_restart: '1G',
      env: {
        NODE_ENV: 'production',
        PORT: 3001
      },
      env_development: {
        NODE_ENV: 'development',
        PORT: 3001
      },
      env_production: {
        NODE_ENV: 'production',
        PORT: 3001
      },
      error_file: './logs/err.log',
      out_file: './logs/out.log',
      log_file: './logs/combined.log',
      time: true,
      // Build before starting in production
      pre_deploy_local: 'npm run build',
      // Health check
      health_check_http: {
        enable: true,
        url: 'http://localhost:3001',
        timeout: 30000,
        interval: 10000
      }
    }
  ]
};
module.exports = {
  apps: [{
    name: "yoyofun-api",
    cwd: "/var/www/yoyofun/backend",
    script: "/var/www/yoyofun/backend/bin/yoyofun-api",
    interpreter: "none",
    autorestart: true,
    max_restarts: 10,
    restart_delay: 3000,
    env: {
      APP_ENV: "production",
      HOST: "127.0.0.1",
      PORT: "9001"
    }
  }]
}

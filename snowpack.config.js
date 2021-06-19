const proxy = require('http2-proxy');
const finalhandler = require('finalhandler')

const defaultWebHandler = (err, req, res) => {
  if (err) {
    console.error('proxy error', err)
    finalhandler(req, res)(err)
  }
}

const proxyHandler = (req, res) => {
  const host = req.headers.host.replace(/:\d+/, "")
  return proxy.web(req, res, {
    hostname: host,
    port: 8081
  }, defaultWebHandler)
}

/** @type {import("snowpack").SnowpackUserConfig } */
module.exports = {
  mount: {
    public: '/',
    src: '/dist',
  },
  plugins: [
    '@snowpack/plugin-svelte',
    '@snowpack/plugin-dotenv',
  ],
  routes: [
    {
      src: '/api/.*',
      dest: proxyHandler
    }
  ],
  optimize: {
    // "bundle": true,
  },
  packageOptions: {
    /* ... */
  },
  devOptions: {
    /* ... */
  },
  buildOptions: {
  },
};

const Bundler = require('parcel-bundler');
const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');

const app = express();

app.use(createProxyMiddleware('/api', {
    target: 'http://127.0.0.1:7777'
}));

const bundler = new Bundler('src/index.html');
app.use(bundler.middleware());

const port = Number(process.env.PORT || 1234);
console.log(`Application available at: http://127.0.0.1:${port}`);
app.listen(port);

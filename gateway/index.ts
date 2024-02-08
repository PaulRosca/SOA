import http from "http";
import express from "express";
import cors from "cors";
import helmet from "helmet";
import rateLimit from "express-rate-limit";
import morgan from "morgan";
import axios from "axios";
import cookies from "cookie-parser";
import { createProxyMiddleware, fixRequestBody } from "http-proxy-middleware";

const limiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 1500,
});

const server = express();

server.use(limiter);
server.use(helmet());
server.use(cors());
server.use(morgan('combined'));
server.use(cookies());
server.use(express.json());
server.use(express.urlencoded({ extended: true }));

const authProxy = createProxyMiddleware({
  target: 'http://gateway.openfaas:8080',
  pathRewrite: { '^/auth': '/function' },
  onProxyReq: fixRequestBody
});

const catalogProxy = createProxyMiddleware({
  target: 'http://catalog-service.default.svc.cluster.local.',
  pathRewrite: { '^/catalog': '/' },
  onProxyReq: fixRequestBody
});

const ordersProxy = createProxyMiddleware({
  target: 'http://orders-service.default.svc.cluster.local.',
  pathRewrite: { '^/orders': '/' },
  onProxyReq: fixRequestBody
});


server.use('/auth', authProxy);
server.use(async (req: express.Request, res: express.Response, next: express.NextFunction) => {
  if (!req.body) {
    req.body = {};
  }
  try {
    const resp = await axios.post('http://gateway.openfaas:8080/function/extract', req.cookies.token, { withCredentials: true });
    if (!req.body) {
      req.body = {};
    }
    req.body.user = resp.data;
    console.log("BODY", req.body, req.path);
  } catch (err) {
    console.log(err);
  }
  return next();
});
server.use('/catalog', catalogProxy);
server.use('/orders', ordersProxy);

const PORT = 9000;
server.listen(PORT, () => {
  console.log(`Server running on PORT: ${PORT}`);
});

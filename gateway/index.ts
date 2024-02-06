import express from "express";
import cors from "cors";
import helmet from "helmet";
import rateLimit from "express-rate-limit";
import morgan from "morgan";
import { createProxyMiddleware } from "http-proxy-middleware";

const limiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 1500,
});

const server = express();

server.use(limiter);
server.use(helmet());
server.use(cors());
server.use(morgan('combined'));

const authProxy = createProxyMiddleware({
  target: 'http://gateway.openfaas:8080',
  pathRewrite: {'^/auth' : '/function'}
});

const catalogProxy = createProxyMiddleware({
  target: 'http://catalog-service.default.svc.cluster.local.',
  pathRewrite: {'^/catalog' : '/'}
});


server.use('/auth', authProxy);
server.use('/catalog', catalogProxy);

const PORT = 9000;
server.listen(PORT, () => {
  console.log(`Server running on PORT: ${PORT}`);
});

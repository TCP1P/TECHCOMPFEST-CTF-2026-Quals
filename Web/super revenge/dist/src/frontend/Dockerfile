# Builder: install deps and create production dist
FROM node:20.18.1-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Runner: ship only what the server needs
FROM node:20.18.1-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/dist ./dist
EXPOSE 2937
ENV HOST=0.0.0.0
ENV PORT=2937

CMD ["node", "./dist/server/entry.mjs"]

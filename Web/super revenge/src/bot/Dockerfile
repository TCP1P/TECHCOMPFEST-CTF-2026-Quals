# Base image
FROM oven/bun:1 AS base
WORKDIR /app

# Install dependencies in dev mode
FROM base AS install
WORKDIR /app
COPY ./package.json ./bun.lockb ./
RUN bun install --frozen-lockfile

# Prerelease build (could do TS build here if needed)
FROM base AS prerelease
WORKDIR /app
COPY --from=install /app/node_modules ./node_modules
COPY ./ ./

RUN bunx playwright@1.47 install-deps

ENV NODE_ENV=production

# Final image
FROM base AS release

RUN apt-get update && apt-get install -y --no-install-recommends \
    libx11-6 libxext6 libxcomposite1 libxdamage1 libxfixes3 libxrandr2 libxcb1 \
    libgtk-3-0 libgbm1 libasound2 libatk1.0-0 libatk-bridge2.0-0 libcups2 libdrm2 \
    libexpat1 libnspr4 libnss3 libxss1 libxshmfence1 libxkbcommon0 \
    libgconf-2-4 libglib2.0-0 libdbus-1-3 libatspi2.0-0 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=prerelease /app ./
RUN bun install --frozen-lockfile --production

# RUN bunx playwright install-deps

USER bun

RUN bunx playwright install chromium
RUN bunx playwright install firefox

EXPOSE 1337

ENTRYPOINT ["bun", "run", "server.ts"]
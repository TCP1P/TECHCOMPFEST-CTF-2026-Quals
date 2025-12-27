FROM oven/bun:1.3.3-debian AS base
WORKDIR /app

FROM base AS install
COPY package.json bun.lockb ./
RUN bun install --frozen-lockfile

FROM base AS prerelease
COPY --from=install /app/node_modules ./node_modules
COPY . ./


FROM base AS release
WORKDIR /app
COPY --from=prerelease /app ./

ENV NODE_ENV=production
RUN bun install --frozen-lockfile --production

# Install OS deps as root, then install browsers as bun user into a shared path
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright
RUN bunx playwright@1.47.0 install-deps \
 && mkdir -p /ms-playwright \
 && chmod -R 777 /ms-playwright

USER bun
RUN bunx playwright@1.47.0 install chromium \
 && bunx playwright@1.47.0 install firefox

EXPOSE 1337
ENTRYPOINT ["bun", "run", "server.ts"]

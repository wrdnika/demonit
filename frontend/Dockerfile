# syntax=docker/dockerfile:1

FROM node:22-alpine AS build
WORKDIR /app

COPY package.json package-lock.json ./
# npm install (not ci) avoids host/container lockfile drift across npm majors
RUN npm install

COPY . .

# Public API URL is overridden at runtime via NUXT_PUBLIC_*;
# bake a sensible default for local docker-compose.
ARG NUXT_PUBLIC_API_BASE=http://localhost:8080
ARG NUXT_PUBLIC_POLL_INTERVAL_MS=15000
ENV NUXT_PUBLIC_API_BASE=$NUXT_PUBLIC_API_BASE \
    NUXT_PUBLIC_POLL_INTERVAL_MS=$NUXT_PUBLIC_POLL_INTERVAL_MS

RUN npm run build

FROM node:22-alpine AS runtime
WORKDIR /app

ENV NODE_ENV=production \
    HOST=0.0.0.0 \
    PORT=3000 \
    NITRO_HOST=0.0.0.0 \
    NITRO_PORT=3000

COPY --from=build /app/.output ./.output

USER node
EXPOSE 3000

CMD ["node", ".output/server/index.mjs"]

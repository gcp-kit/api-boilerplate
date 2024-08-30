FROM node:18-alpine AS base

FROM base AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /app

COPY docker/swagger_cli/package.json ./

RUN npm i

FROM base AS runner
WORKDIR /app

COPY --from=deps /app/package.json ./package.json
COPY --from=deps /app/node_modules ./node_modules

CMD ["npm", "run", "bundle"]

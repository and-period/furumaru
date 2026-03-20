import { readdirSync, readFileSync, writeFileSync } from 'node:fs'
import { join } from 'node:path'

const pnpmDir = join(process.cwd(), 'node_modules', '.pnpm')
const targetSuffix = join('node_modules', '@nuxt', 'nitro-server', 'dist', 'index.mjs')
const before = 'if (options.cors(event.path)) {'
const after = 'if (options?.cors?.(event.path)) {'
const corsBefore = 'if (handleCors(event, nuxt.options.devServer.cors)) return null;'
const corsAfter = 'if (nuxt.options.devServer?.cors && handleCors(event, nuxt.options.devServer.cors)) return null;'

function patchNuxtNitroServer() {
  const entries = readdirSync(pnpmDir, { withFileTypes: true })
  const candidates = entries
    .filter(entry => entry.isDirectory() && entry.name.startsWith('@nuxt+nitro-server@'))
    .map(entry => join(pnpmDir, entry.name, targetSuffix))

  for (const filePath of candidates) {
    let content
    try {
      content = readFileSync(filePath, 'utf8')
    }
    catch {
      continue
    }

    const patched = content
      .replace(before, after)
      .replace(corsBefore, corsAfter)

    if (patched !== content) {
      writeFileSync(filePath, patched)
    }
  }
}

patchNuxtNitroServer()

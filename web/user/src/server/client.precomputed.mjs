import { existsSync, readdirSync } from 'node:fs'
import { readFile } from 'node:fs/promises'
import { join } from 'node:path'

function createResource(file, manifest) {
  const entry = manifest[file] || { file }
  const ext = (entry.file || '').split('.').pop()
  const isScript = ext === 'js' || ext === 'mjs'
  const isStyle = ext === 'css'

  return {
    ...entry,
    file: entry.file || file,
    module: isScript,
    resourceType: isScript ? 'script' : isStyle ? 'style' : undefined,
  }
}

function getModuleDependencies(id, manifest, cache) {
  if (cache[id]) {
    return cache[id]
  }

  const dependencies = {
    scripts: {},
    styles: {},
    preload: {},
    prefetch: {},
  }
  cache[id] = dependencies

  const meta = manifest[id]
  if (!meta) {
    return dependencies
  }

  if (meta.file) {
    const resource = createResource(id, manifest)
    dependencies.preload[id] = resource
    if (meta.isEntry || meta.sideEffects) {
      dependencies.scripts[id] = resource
    }
  }

  for (const css of meta.css || []) {
    const resource = createResource(css, manifest)
    dependencies.styles[css] = resource
    dependencies.preload[css] = resource
    dependencies.prefetch[css] = resource
  }

  for (const asset of meta.assets || []) {
    const resource = createResource(asset, manifest)
    dependencies.preload[asset] = resource
    dependencies.prefetch[asset] = resource
  }

  for (const depId of meta.imports || []) {
    const depDeps = getModuleDependencies(depId, manifest, cache)
    Object.assign(dependencies.styles, depDeps.styles)
    Object.assign(dependencies.preload, depDeps.preload)
    Object.assign(dependencies.prefetch, depDeps.prefetch)
    if (meta.isEntry || depDeps._hasScript) {
      Object.assign(dependencies.scripts, depDeps.scripts)
    }
  }

  dependencies._hasScript = Object.keys(dependencies.scripts).length > 0
  return dependencies
}

function normalizeDependencies(dependencies) {
  const prefetch = {}

  for (const [id, dep] of Object.entries(dependencies.prefetch)) {
    if (dep.prefetch) {
      prefetch[id] = dep
    }
  }

  dependencies.prefetch = prefetch

  for (const id of Object.keys(dependencies.preload)) {
    delete dependencies.prefetch[id]
  }

  for (const id of Object.keys(dependencies.styles)) {
    delete dependencies.preload[id]
    delete dependencies.prefetch[id]
  }

  delete dependencies._hasScript
  return dependencies
}

async function loadManifest() {
  const manifestPaths = [
    join(process.cwd(), 'public/_nuxt/manifest.json'),
    join(process.cwd(), '.output/public/_nuxt/manifest.json'),
    join(process.cwd(), 'node_modules/.cache/nuxt/.nuxt/dist/client/manifest.json'),
    join(process.cwd(), '.nuxt/dist/client/manifest.json'),
  ]

  for (const manifestPath of manifestPaths) {
    try {
      if (!existsSync(manifestPath)) {
        continue
      }

      const content = await readFile(manifestPath, 'utf8')
      const manifest = JSON.parse(content)
      if (manifest && typeof manifest === 'object') {
        return manifest
      }
    }
    catch {
      // Ignore and try the next path.
    }
  }

  return {}
}

function loadCssFallbackPrecomputed() {
  const assetDirs = [
    join(process.cwd(), 'public/_nuxt'),
    join(process.cwd(), '.output/public/_nuxt'),
    join(process.cwd(), '.nuxt/dist/client/_nuxt'),
  ]

  for (const assetDir of assetDirs) {
    if (!existsSync(assetDir)) {
      continue
    }

    const scripts = {}
    const styles = {}
    for (const file of readdirSync(assetDir)) {
      if (file.endsWith('.js')) {
        scripts[file] = {
          file,
          module: true,
          resourceType: 'script',
        }
        continue
      }

      if (!file.endsWith('.css')) {
        continue
      }
      styles[file] = {
        file,
        resourceType: 'style',
      }
    }

    if (Object.keys(styles).length > 0) {
      return {
        dependencies: {
          __entry__: {
            scripts,
            styles,
            preload: {},
            prefetch: {},
          },
        },
        entrypoints: ['__entry__'],
      }
    }
  }

  return {
    dependencies: {},
    entrypoints: [],
  }
}

export default async function clientPrecomputed() {
  const manifest = await loadManifest()
  const entrypoints = Object.entries(manifest)
    .filter(([, meta]) => meta && meta.isEntry)
    .map(([id]) => id)

  if (Object.keys(manifest).length === 0 || entrypoints.length === 0) {
    return loadCssFallbackPrecomputed()
  }

  const cache = {}
  const dependencies = {}

  for (const id of Object.keys(manifest)) {
    dependencies[id] = normalizeDependencies(getModuleDependencies(id, manifest, cache))
  }

  return {
    dependencies,
    entrypoints,
  }
}

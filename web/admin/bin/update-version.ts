import fs from 'fs'
import path from 'path'

interface VersionInfo {
  version: string
  buildDate: string
}
const versionFilePath = path.resolve(__dirname, '..', 'version.json')
const packageJson = require('../package.json')

const currentDate = new Date().toISOString()

const versionFile: VersionInfo = {
  version: packageJson.version,
  buildDate: currentDate,
}

fs.writeFileSync(versionFilePath, JSON.stringify(versionFile, null, 2))

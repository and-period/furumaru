import * as esbuild from 'esbuild';

await esbuild.build({
  // General options
  bundle: true,
  platform: 'node',
  tsconfig: 'tsconfig.json',
  // Input
  entryPoints: ['src/app.ts'],
  // Output contents
  // Output location
  outfile: 'app.js',
  // Path resolution
  external: [],
  // Transofrmation
  target: 'node20',
  // Optimization
  minify: true,
  // Source maps
  sourcemap: false,
  // Build metadata
  // Logging
});

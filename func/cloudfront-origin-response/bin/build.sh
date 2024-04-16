#/bin/sh

if [ ! -f "package.json" ]; then
  echo "Error: package.json not found."
  exit 1
fi

### clean up
rm -rf ./dist ./app.zip
mkdir -p ./dist

### build
npm run build

mv ./app.js ./dist/index.js
cp -r ./node_modules ./dist/

### remove dev dependencies
dependencies=$(jq -r '.devDependencies | keys | .[]' package.json)

for dependency in $dependencies; do
  rm -rf ./dist/node_modules/${dependency}
done

### remove unnecessary files
rm -rf ./dist/node_modules/**/test
rm -rf ./dist/node_modules/**/.eslintrc.yml

### compress
cd ./dist
zip -ry ./app.zip .

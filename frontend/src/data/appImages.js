const appImageModules = import.meta.glob('../assets/images/app/*', {
  eager: true,
  import: 'default',
})

export function appImage(filename) {
  const entry = Object.entries(appImageModules).find(([path]) => path.endsWith('/' + filename))
  return entry ? entry[1] : ''
}

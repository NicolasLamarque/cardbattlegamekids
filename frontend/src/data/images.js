const characterImages = import.meta.glob('../assets/images/jujutsu-kaisen/characters/*.jpg', {
  eager: true,
  import: 'default',
})

export function localCharacterImage(filename) {
  const entry = Object.entries(characterImages).find(([path]) => path.endsWith('/' + filename))
  return entry ? entry[1] : ''
}

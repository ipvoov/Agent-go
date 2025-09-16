function randomPart(len) {
  const chars = 'abcdefghijklmnopqrstuvwxyz0123456789'
  let out = ''
  for (let i = 0; i < len; i++) {
    out += chars[Math.floor(Math.random() * chars.length)]
  }
  return out
}

export function generateSessionId() {
  const ts = Date.now().toString(36)
  return `sess_${ts}_${randomPart(6)}`
}



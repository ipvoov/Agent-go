function cjkScore(s) {
  if (!s) return 0
  const matches = s.match(/[\u4e00-\u9fff\u3400-\u4dbf\uf900-\ufaff]/g)
  return matches ? matches.length : 0
}

export function decodeLatin1Utf8(text) {
  if (!text) return text
  let best = text
  let bestScore = cjkScore(text)

  // Strategy 1: Interpret current JS string's 0..255 codepoints as raw bytes -> UTF-8
  try {
    const bytes = new Uint8Array(Array.from(text, ch => ch.charCodeAt(0) & 0xff))
    const decoded = new TextDecoder('utf-8', { fatal: false }).decode(bytes)
    const score = cjkScore(decoded)
    if (score > bestScore) { best = decoded; bestScore = score }
  } catch {}

  // Strategy 2: legacy escape/decodeURIComponent trick
  try {
    // escape converts to %xx for 0..255; decodeURIComponent then treats as UTF-8
    /* eslint-disable deprecate/escape */
    const decoded2 = decodeURIComponent(escape(text))
    const score2 = cjkScore(decoded2)
    if (score2 > bestScore) { best = decoded2; bestScore = score2 }
  } catch {}

  return best
}


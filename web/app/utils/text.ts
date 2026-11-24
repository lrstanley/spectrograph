// Length in UTF-8 bytes, which is what the backend's length limits count.
export function byteLength(value: string) {
  return new TextEncoder().encode(value).length
}

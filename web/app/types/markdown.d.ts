// "*.md" imports are rendered to HTML at build time by build/markdown.ts.
declare module "*.md" {
  const doc: {
    html: string
    headings: { id: string; text: string; depth: number }[]
  }
  export default doc
}

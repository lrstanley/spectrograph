import { Marked } from "marked"
import type { Plugin } from "vite"

export type MarkdownHeading = {
  id: string
  text: string
  depth: number
}

function slugify(value: string) {
  return value
    .toLowerCase()
    .replace(/<[^>]+>/g, "")
    .replace(/&[a-z0-9#]+;/g, "")
    .replace(/[^a-z0-9\s-]/g, "")
    .trim()
    .replace(/\s+/g, "-")
}

// Renders "*.md" imports to HTML at build time, exporting the HTML and the
// list of headings (for a table of contents). Content is trusted (it lives in
// this repository).
export function markdown(): Plugin {
  return {
    name: "spectrograph:markdown",
    enforce: "pre",
    async transform(source, id) {
      const path = id.split("?")[0] ?? id
      if (!path.endsWith(".md")) {
        return
      }

      const headings: MarkdownHeading[] = []
      const seen = new Map<string, number>()
      const marked = new Marked({
        gfm: true,
        renderer: {
          heading({ tokens, depth }) {
            const text = this.parser.parseInline(tokens)
            let slug = slugify(text)
            const count = seen.get(slug) ?? 0
            seen.set(slug, count + 1)
            if (count > 0) {
              slug = `${slug}-${count}`
            }
            headings.push({ id: slug, text: text.replace(/<[^>]+>/g, ""), depth })
            return `<h${depth} id="${slug}"><a href="#${slug}">${text}</a></h${depth}>\n`
          },
        },
      })

      const html = await marked.parse(source)
      return { code: `export default ${JSON.stringify({ html, headings })}`, map: null }
    },
  }
}

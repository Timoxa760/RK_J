/**
 * After `npx shadcn-vue add ...`, run: npm run shadcn:patch
 *
 * Replaces `extends PrimitiveProps` with local PrimitiveForwardProps
 * so Vue SFC macro does not resolve reka-ui types (see lib/reka-props.ts).
 */
import { readFileSync, writeFileSync, globSync } from 'node:fs'
import { join } from 'node:path'

const root = join(import.meta.dirname, '..', 'components', 'ui')
const files = globSync(join(root, '**/*.vue'))

const extendsPrimitive =
  /extends\s+(?:\/\* @vue-ignore \*\/\s*)?PrimitiveProps\b/g

const definePrimitive =
  /defineProps<\s*PrimitiveProps\b/g

const importPrimitiveType =
  /import type \{([^}]*)\bPrimitiveProps\b([^}]*)\} from "reka-ui"\n/

let patched = 0
for (const file of files) {
  let src = readFileSync(file, 'utf8')
  let next = src

  if (extendsPrimitive.test(src) || definePrimitive.test(src)) {
    if (!next.includes("PrimitiveForwardProps")) {
      if (next.includes("from 'vue'") || next.includes('from "vue"')) {
        next = next.replace(
          /(import type \{ HTMLAttributes \} from "vue"\n)/,
          '$1import type { PrimitiveForwardProps } from \'~/lib/reka-props\'\n',
        )
      } else {
        next = next.replace(
          /<script setup lang="ts">\n/,
          '<script setup lang="ts">\nimport type { PrimitiveForwardProps } from \'~/lib/reka-props\'\n',
        )
      }
    }

    next = next
      .replace(extendsPrimitive, 'extends PrimitiveForwardProps')
      .replace(definePrimitive, 'defineProps<PrimitiveForwardProps')
      .replace(importPrimitiveType, (line, before, after) => {
        const types = `${before}${after}`
          .split(',')
          .map((t) => t.trim())
          .filter((t) => t && t !== 'PrimitiveProps')
        if (!types.length) return ''
        return `import type { ${types.join(', ')} } from "reka-ui"\n`
      })
  }

  if (next !== src) {
    writeFileSync(file, next)
    patched++
    console.log('patched', file)
  }
}

console.log(patched ? `Done: ${patched} file(s)` : 'No files needed patching')

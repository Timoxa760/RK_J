import type { Plugin } from 'vite'

const NEEDLE = 'instance.__vrv_devtools = info;'
const REPLACEMENT = 'if (instance) instance.__vrv_devtools = info;'

function patchVueRouterSource(code: string): string | null {
  if (!code.includes('__vrv_devtools')) return null
  if (!code.includes(NEEDLE)) return null
  return code.replaceAll(NEEDLE, REPLACEMENT)
}

/**
 * vue-router 4.6.x: ref без instance → `instance.__vrv_devtools = info` → 500 в dev.
 */
export function patchVueRouterDevtools(): Plugin {
  return {
    name: 'patch-vue-router-devtools',
    enforce: 'pre',
    transform(code, id) {
      if (!id.includes('vue-router')) return
      return patchVueRouterSource(code) ?? undefined
    },
    config() {
      return {
        optimizeDeps: {
          esbuildOptions: {
            plugins: [
              {
                name: 'patch-vue-router-devtools-deps',
                setup(build) {
                  build.onLoad({ filter: /vue-router\/dist\/vue-router\.mjs$/ }, async (args) => {
                    const { readFileSync } = await import('node:fs')
                    const contents = readFileSync(args.path, 'utf8')
                    const patched = patchVueRouterSource(contents)
                    if (!patched) return null
                    return { contents: patched, loader: 'js' }
                  })
                }
              }
            ]
          }
        }
      }
    }
  }
}

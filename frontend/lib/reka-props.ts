import type { Component } from 'vue'

/** Runtime-safe subset of reka-ui PrimitiveProps (avoids SFC macro resolving reka-ui types). */
export interface PrimitiveForwardProps {
  as?: string | Component
  asChild?: boolean
}

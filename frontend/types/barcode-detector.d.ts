interface BarcodeDetectorOptions {
  formats?: string[]
}

interface DetectedBarcode {
  rawValue?: string
}

declare class BarcodeDetector {
  constructor(options?: BarcodeDetectorOptions)
  detect(source: ImageBitmap): Promise<DetectedBarcode[]>
}

interface Window {
  BarcodeDetector?: typeof BarcodeDetector
}

export {}

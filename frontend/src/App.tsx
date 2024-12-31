import { useState, useRef, useEffect } from 'react'
import { ProcessImage } from '../wailsjs/go/main/App'
import { main } from '../wailsjs/go/models'

function App() {
  const [image, setImage] = useState<string>('')
  const [ocrResult, setOcrResult] = useState<main.OCRResult | null>(null)
  const [scale, setScale] = useState<number>(1)
  const imageRef = useRef<HTMLImageElement>(null)

  const handleImageUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (e) => {
        const result = e.target?.result as string
        setImage(result)
        processOCR(result)
      }
      reader.readAsDataURL(file)
    }
  }

  const processOCR = async (imageData: string) => {
    try {
      const result = await ProcessImage(imageData)
      setOcrResult(result)
    } catch (error) {
      console.error('OCR processing failed:', error)
    }
  }

  useEffect(() => {
    if (imageRef.current && ocrResult) {
      const img = new Image()
      img.src = image
      img.onload = () => {
        const displayWidth = imageRef.current?.clientWidth || 0
        const originalWidth = img.width
        setScale(displayWidth / originalWidth)
      }
    }
  }, [image, ocrResult])

  return (
    <div className="container">
      <input type="file" accept="image/*" onChange={handleImageUpload} />
      
      <div className="image-container" style={{ position: 'relative' }}>
        {image && <img ref={imageRef} src={image} alt="Uploaded image" />}
        
        {ocrResult?.boxes.map((box, index) => (
          <div
            key={index}
            style={{
              position: 'absolute',
              left: Math.round(box.x * scale) + 10,
              top: Math.round(box.y * scale) + 10,
              width: Math.round(box.width * scale),
              height: Math.round(box.height * scale),
              border: '1px solid red',
              cursor: 'pointer',
            }}
            onClick={() => navigator.clipboard.writeText(box.text)}
            title={box.text}
          />
        ))}
      </div>

      {ocrResult && (
        <div className="text-output">
          <h3>Extracted Text:</h3>
          <pre>{ocrResult.text}</pre>
        </div>
      )}
    </div>
  )
}

export default App
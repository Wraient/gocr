import { useState, useRef, useEffect } from 'react'
import { GetInitialImage, ProcessImageFile, GetImageData, ProcessImage } from '../wailsjs/go/main/App'
import { main } from '../wailsjs/go/models'

function App() {
  const [image, setImage] = useState<string>('')
  const [ocrResult, setOcrResult] = useState<main.OCRResult | null>(null)
  const [scale, setScale] = useState<number>(1)
  const imageRef = useRef<HTMLImageElement>(null)

  useEffect(() => {
    const loadInitialImage = async () => {
      try {
        const initialPath = await GetInitialImage()
        if (initialPath) {
          const ocrResult = await ProcessImageFile(initialPath)
          setOcrResult(ocrResult)

          const imageData = await GetImageData(initialPath)
          setImage(imageData)
        }
      } catch (error) {
        console.error('Failed to load initial image:', error)
      }
    }
    loadInitialImage()
  }, [])

  const handleImageUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = async (e) => {
        const result = e.target?.result as string
        setImage(result)
        
        // Process OCR using ProcessImage for base64 data
        const ocrResult = await ProcessImage(result)
        setOcrResult(ocrResult)
      }
      reader.readAsDataURL(file)
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
      
      <div className="image-container">
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
# go-vlfeat
vlfeat 库的 go 语言版本实现

目前版本支持的 vlfeat v0.9.21

sample code:

```
mat, _ := gocv.IMRead("1.jpg", gocv.IMReadGrayScale)
floatImg := gocv.NewMat()
mat.ConvertTo(&floatImg, gocv.MatTypeCV32F)
floaImgData, _ := floatImg.DataPtrFloat32()
sift := vlfeat.NewSift(mat.Cols(), mat.Rows(), 4, 2, 0)
if sift.ProcessFirstOctave(floaImgData) != vlfeat.VlErrorEOF {
    for {
        sift.Detect()
        detectedKeypoints := sift.GetKeypoints()
        nkeys := sift.GetNkeypoints()
        for i := 0; i < nkeys; i++ {
            angleCount, angles := sift.CalcKeypointOrientations(detectedKeypoints[i])
            for j := 0; j < angleCount; j++ {
                angleDesc := sift.CalcKeypointDescriptor(128, detectedKeypoints[i], angles[j])
                // 处理 angleDesc
            }
        }
        if sift.ProcessNextOctave() == vlfeat.VlErrorEOF {
            break
        }
    }
}
defer sift.Delete()
```
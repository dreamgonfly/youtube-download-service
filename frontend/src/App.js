import React, { Fragment, useState, useEffect } from 'react'
import api from './api'
import url from './url'
import Header from './components/Header/Header'
import Footer from './components/Footer/Footer'
import Input from './components/Input/Input'
import Output from './components/Output/Output'
import './index.css'

function App() {
  const [data, setData] = useState({ Thumbnail: "", Name: "", Formats: [] })
  const [videoId, setVideoId] = useState("")
  const [waiting, setWaiting] = useState(false)


  useEffect(() => {
    const videoIdFromQuery = url.extractVideoIdFromURL(window.location.href)
    if (videoIdFromQuery !== null) {
      setVideoId(videoIdFromQuery)
      api.preview(videoIdFromQuery).then((res) => {
        setData(res.data)
        api.updateThumbnail(videoIdFromQuery, res.data.Thumbnail, res.data.Name).then((res) => {
          setData((prev) => { return { ...prev, Thumbnail: res.data.Thumbnail } })
        })
      })
    }
  }, [])

  const inputSetVideoIdHandler = (videoId) => {
    setWaiting(true)
    setVideoId(videoId)
    api.preview(videoId).then((res) => {
      setData(res.data)
      setWaiting(false)
      api.updateThumbnail(videoId, res.data.Thumbnail, res.data.Name).then((res) => {
        setData((prev) => { return { ...prev, Thumbnail: res.data.Thumbnail } })
      })
    })
  }

  return (
    <Fragment>
      <Header />
      <main>
        <Input videoId={videoId} onSetVideoId={inputSetVideoIdHandler} />
        <Output videoId={videoId} data={data} waiting={waiting} />
      </main>
      <Footer />
    </Fragment>
  )
}

export default App

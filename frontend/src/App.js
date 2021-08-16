import React, { Fragment, useState, useEffect } from 'react'
import api from './api'
import { extractVideoIdFromURL } from './url'
import Header from './components/Header/Header'
import Footer from './components/Footer/Footer'
import Input from './components/Input/Input'
import Output from './components/Output/Output'
import './index.css'

function App() {
  const [videoId, setVideoId] = useState("")
  const [data, setData] = useState({})
  const [waiting, setWaiting] = useState(false)

  const inputHandler = (videoId) => {
    setWaiting(true)
    setVideoId(videoId)
    api.preview(videoId).then((res) => {
      setWaiting(false)
      setData(res.data)
      api.updateThumbnail(videoId, res.data.Thumbnail, res.data.Name).then((res) => {
        setData((prev) => { return { ...prev, Thumbnail: res.data.Thumbnail } })
      })
    })
  }

  useEffect(() => {
    const videoIdFromQuery = extractVideoIdFromURL(window.location.href)
    if (videoIdFromQuery) {
      inputHandler(videoIdFromQuery)
    }
  }, [])

  return (
    <Fragment>
      <Header />
      <main>
        <Input videoId={videoId} onInput={inputHandler} />
        <Output videoId={videoId} data={data} waiting={waiting} />
      </main>
      <Footer />
    </Fragment>
  )
}

export default App

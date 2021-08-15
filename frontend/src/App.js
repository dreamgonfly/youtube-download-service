import React, { Fragment, useState, useEffect, useRef } from 'react';
import api from './api';
import url from './url';
import styles from './index.module.css';
import Loader from "react-loader-spinner";

function App() {
  const [data, setData] = useState({ Thumbnail: "", Name: "", Formats: [] })
  const [videoId, setVideoId] = useState("")
  const [selectedIndex, setSelectedIndex] = useState(undefined)
  const [optionOpen, setOptionOpen] = useState(false)
  const [loading, setLoading] = useState(false)

  const inputRef = useRef()

  useEffect(() => {
    const videoIdFromQuery = url.extractVideoIdFromURL(window.location.href)
    if (videoIdFromQuery !== null) {
      inputRef.current.value = url.reconstructYoutubeURL(videoIdFromQuery)
      api.preview(videoIdFromQuery).then((res) => {
        setData(res.data)
        setSelectedIndex(res.data.Formats[0].FormatId)
        api.updateThumbnail(videoIdFromQuery, res.data.Thumbnail, res.data.Name).then((res) => {
          setData((prev) => { return { ...prev, Thumbnail: res.data.Thumbnail } })
        })
      })
    }
  }, [])

  const previewClickHandler = (event) => {
    const isValidURL = url.isValidURL(inputRef.current.value)
    let videoId = inputRef.current.value
    if (isValidURL) {
      videoId = url.extractVideoIdFromURL(inputRef.current.value)
    }
    setVideoId(videoId)
    api.preview(videoId).then((res) => {
      setData(res.data)
      setSelectedIndex(res.data.Formats[0].FormatId)
      api.updateThumbnail(videoId, res.data.Thumbnail, res.data.Name).then((res) => {
        setData((prev) => { return { ...prev, Thumbnail: res.data.Thumbnail } })
      })
    })
  }

  const viewClickHandlerConstructur = (id, format, filename) => {
    const viewClickHandler = () => {
      setLoading(true)
      api.save(id, format, filename).then((res) => {
        setLoading(false)
        // Secure solution. https://stackoverflow.com/questions/45046030/maintaining-href-open-in-new-tab-with-an-onclick-handler-in-react
        const newWindow = window.open(res.data.URL, '_blank', 'noopener,noreferrer')
        if (newWindow) newWindow.opener = null
      })
    }
    return viewClickHandler
  }

  let options = []
  if (data.Formats.length > 0) {
    options = data.Formats.filter((element) => { return element.FormatId !== selectedIndex })
  }

  const choiceClickHandler = () => {
    setOptionOpen((prev) => { return !prev })
  }

  const optionClickHandler = (event) => {
    setSelectedIndex(event.currentTarget.dataset.index)
    setOptionOpen(false)
  }

  // const data = {
  //   Thumbnail: "https://storage.googleapis.com/youtube-download-backend-beta/videos/GSVsfCCtRr0/%5B%EA%B8%B0%EC%83%9D%EC%B6%A9%5D%2030%EC%B4%88%20%EC%98%88%EA%B3%A0.jpg?X-Goog-Algorithm=GOOG4-RSA-SHA256&X-Goog-Credential=youtube-download-service%40youtube-download-service.iam.gserviceaccount.com%2F20210803%2Fauto%2Fstorage%2Fgoog4_request&X-Goog-Date=20210803T004718Z&X-Goog-Expires=899&X-Goog-Signature=9559009897db3c9d382ed73cbbfc3d349f00b798b4ab6e0980d4acabaca09f87757721933b97368d69c1f6076aa1a01f8ea7ff9ce3d3aa0df8d700bccace51f2487f168d238ab1b48b5986172bc140c3020251b22b6ea5a3fac6f588e6f5cbf192fa4a5178d9879dfe067622735350bbd14a376f5847dcdfe13ae54799222e78295372843e335f9d55172c7b98cfe71dd3cc79d7bbf53a7c908be9b4bb7b589b93337cfa0e01c08aca5dd33f025833168f1fc0ee6e48c411ed68b968d1ff2775a386a3a6b8b99962e13721785c133b2e300b1e957160d5a1c01d1f54f6719e298bd316af0fcab50c110c37ef26a878cb2ac5723e7bfcf186462c0f428e878b0b&X-Goog-SignedHeaders=host",
  //   Formats: [
  //     {
  //       Filesize: 1348634,
  //       FormatId: "18",
  //       FormatNote: "360p",
  //       Ext: "mp4"
  //     },
  //     {
  //       Filesize: 2059470,
  //       FormatId: "22",
  //       FormatNote: "720p",
  //       Ext: "mp4"
  //     }
  //   ]
  // }

  let optionsBlock = undefined
  if (optionOpen) {
    optionsBlock = <div className={styles["output-format-options"]}>
      {options.map((item, index) => { return <div key={index} data-index={item.FormatId} className={styles["option-desc"]} onClick={optionClickHandler}>{item.Ext} {item.FormatNote} {Math.round(item.Filesize / 1000 / 1000 * 10) / 10} MB</div> })}
    </div>
  }

  let LoadingIcon = undefined
  if (loading) {
    LoadingIcon = <Loader
      type="TailSpin"
      color="#00BFFF"
      height={20}
      width={20}
    />
  } else {
    LoadingIcon = "Open In New Tab"
  }

  let outputBlock = <div></div>
  if (data.Formats.length > 0 && selectedIndex) {
    const minutes = Math.floor(data.DurationSecond / 60)
    const seconds = data.DurationSecond - minutes * 60;

    outputBlock = <div className={styles["output-card"]}>
      <div className={styles["output-content"]}>
        <div className={styles["output-thumbnail"]}>
          <img className={styles["img"]} src={data.Thumbnail} alt="thumbnail" />
        </div>
        <div className={styles["output-title"]}>{data.Title}</div>
        <div className={styles["output-length"]}>{minutes}:{seconds}</div>
      </div>
      <div className={styles["output-download"]}>
        <div className={styles["output-format-choice"]}>
          <div className={styles["format-desc"]}>{data.Formats.find(element => element.FormatId === selectedIndex).Ext} {data.Formats.find(element => element.FormatId === selectedIndex).FormatNote} {Math.round(data.Formats.find(element => element.FormatId === selectedIndex).Filesize / 1000 / 1000 * 10) / 10} MB</div>
          <div className={styles["format-choice-arrow"]} onClick={choiceClickHandler}><i className={`fas fa-caret-down ${styles["arrow-down"]}`}></i></div>
        </div>
        {optionsBlock}
        <div>
          <div className={styles["output-action"]}>
            <a className={styles['output-link']} href={api.composeDownloadLink(videoId, data.Formats.find(element => element.FormatId === selectedIndex).FormatId, data.Name + "." + data.Formats.find(element => element.FormatId === selectedIndex).Ext)}><button className={styles["output-download-button"]}>Download</button></a>
          </div>
          <div className={styles["output-action"]}>
            <button className={styles["output-view"]} onClick={viewClickHandlerConstructur(videoId, data.Formats.find(element => element.FormatId === selectedIndex).FormatId, data.Name + "." + data.Formats.find(element => element.FormatId === selectedIndex).Ext)}>
              {LoadingIcon}
            </button>
          </div>
        </div>
      </div>
    </div>
  }

  return (
    <Fragment>
      {/* <h1>Youtube Download Service</h1>
      <input type="text" value={enteredInput} onChange={inputChangeHandler} />
      <button onClick={previewClickHandler}>Preview</button>
      <br />
      {data.Formats.map((item) => {
        return <Fragment key={item.FormatId}>
          <br />
          <a href={api.composeDownloadLink(videoId, item.FormatId, data.Name + "." + item.Ext)}><button>Download {item.FormatNote}</button></a>
          <button onClick={viewClickHandlerConstructur(videoId, item.FormatId, data.Name + "." + item.Ext)}>View {item.FormatNote}</button>
          <br />
        </Fragment>
      })}
      {thumbnailElement} */}
      <header>
        <div className={styles.inner}>
          <div className={styles.headerContainer}>
            <div className={styles["header-icon"]}></div>
          </div>
        </div>
      </header>
      <section className={styles["input"]}>
        <div className={styles.inner}>
          <div className={styles["title-container"]}>
            <div className={styles['title-logo']}>Dryoutube.com</div>
            <div className={styles["title-desc"]}>Online Youtube Downloader</div>
          </div>
          <div className={styles["input-container"]}>
            <div className={styles["input-logo"]}><img className={styles["logo-img"]} src="logo.png" alt="" /></div>
            <input ref={inputRef} className={styles["input-bar"]} type="text" placeholder="Put your video link here"></input>
            <button className={styles["input-button"]} onClick={previewClickHandler}>Go</button>
          </div>
        </div>
      </section>
      <section className={styles["output"]}>
        <div className={styles.inner}>
          {outputBlock}
        </div>
      </section>
      <footer>
        <div className={styles.inner}>
          <div className={styles["footer-message"]}>Online Youtube Downloader</div>
          <div className={styles["footer-contact"]}>dreamgonfly@gmail.com</div>
          <div className={styles["footer-copyright"]}>Copyright 2021 © All rights reserved.</div>
        </div>
      </footer>
    </Fragment>
  );
}

export default App;

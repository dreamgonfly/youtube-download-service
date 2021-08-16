import React, { useRef } from 'react'
import styles from './Input.module.css'
import { reconstructYoutubeURL, getVideoId } from '../../url'

const Input = (props) => {
    const inputRef = useRef()

    if (props.videoId !== "") {
        inputRef.current.value = reconstructYoutubeURL(props.videoId)
    }

    const previewClickHandler = () => {
        const videoId = getVideoId(inputRef.current.value)
        props.onInput(videoId)
    }

    return (
        <section className={styles["input"]}>
            <div className="inner">
                <div className={styles["title-container"]}>
                    <div className={styles['title-logo']}>Dryoutube.com</div>
                    <div className={styles["title-desc"]}>Online Youtube Downloader</div>
                </div>
                <div className={styles["input-container"]}>
                    <div className={styles["input-logo"]}>
                        <img className={styles["logo-img"]} src="logo.png" alt="" />
                    </div>
                    <input ref={inputRef} className={styles["input-bar"]} type="text" placeholder="Put your video link here"></input>
                    <button className={styles["input-button"]} onClick={previewClickHandler}>Go</button>
                </div>
            </div>
        </section>
    )
}

export default Input